#!/usr/bin/env python3
# Copyright (c) 2024 Bitcoin Association
# Distributed under the Open BSV software license, see the accompanying file LICENSE.

from datetime import datetime
import decimal
import http.client as asm
import json
import re
import subprocess
import sys
import time

# Representing a process we can open.
class Process():
    def __init__(self, command, path=None):
        self.command = command
        self.path = path
        self.process = None
    
    def open(self, blocking = True, stderr = None):
        report(f"Running {self.command}...")
        if blocking:
            # Open the process and wait for it to finish
            self.process = subprocess.Popen(self.command, 
                                            stdin=subprocess.PIPE, 
                                            stdout=subprocess.PIPE, 
                                            stderr=subprocess.PIPE, 
                                            universal_newlines=True,
                                            cwd=self.path)
            cli_stdout, cli_stderr = self.process.communicate()
            return_code = self.process.poll()
            if return_code:
                raise subprocess.CalledProcessError(return_code, self.command, stderr=cli_stderr)
            return cli_stdout
        else:
            # Open the process and return immediately
            self.process = subprocess.Popen(self.command, universal_newlines=True, cwd=self.path, stderr=stderr)

# Representing a call over the SSH with the key-based authentication.
class SSHCall():

    def __init__(self, host, user, key):
        self.ssh_command = ["ssh", "-i", key, f"{user}@{host}"]
        self.process = None
    
    def run(self, command, blocking = True, stderr=None):
        ssh_call = self.ssh_command + [command]
        self.process = Process(ssh_call)
        return self.process.open(blocking=blocking, stderr=stderr)

# Representing the Alert System Microservice.
class ASM():

    def __init__(self, port = None, ssh_args = {}, timeout = 60):
        self.timeout = timeout
        self.command = ["alert-system"]
        # default port
        if port is None:
            port = 3000
        self.ssh = None
        host = "localhost"
        if len(ssh_args) == 3:
            self.ssh = SSHCall(ssh_args.get("host"), ssh_args.get("user"), ssh_args.get("pk_path"))
            host = ssh_args.get("host")
        self.service = f"{host}:{port}"
        self.process = None
    
    def wait_for_synced(self, process = None):
        wait_until = time.time() + self.timeout
        while time.time() < wait_until:
            if process is not None:
                # Check if process (either ssh or alert-system) already terminated while waiting for health status
                return_code = process.process.poll()
                if return_code is not None:
                    command = self.command
                    if self.ssh:
                        command = self.ssh.ssh_command + command
                    stderr = None
                    if process.process.stderr is not None:
                        stderr = process.process.stderr.read()
                    raise subprocess.CalledProcessError(return_code, command, stderr=stderr)
            if (self.is_synced()):
                return
            report(f"Not synced")
            report(f"Retrying...")
            time.sleep(1.0)
        assert wait_until >= time.time(), "Alert System Microservice not synced, timeout exceeded"

    def is_synced(self):
        try:
            connection = asm.HTTPConnection(self.service)
            connection.request('GET', "/health")
            health_response = connection.getresponse()
            assert health_response.status == 200, f"Alert System Microservice response is {health_response.status}"
            health_response = health_response.read()
            health = json.loads(health_response)
            return health["synced"]
        except Exception as e:
            report_exception("Exception while requesting Alert System Microservice health", e)
            return None


    def run(self):
        process = None
        if self.ssh:
            # We want to get the stderr of the SSH process to be able to report any issues.
            self.ssh.run(' '.join(self.command), blocking=False, stderr=subprocess.PIPE)
            process = self.ssh.process
        else:
            process = Process(self.command)
            # We want to get the stderr to be able to report any alert-system issues
            process.open(blocking=False, stderr=subprocess.PIPE)
        report("Waiting for the ASM to be synced...")
        self.wait_for_synced(process)
        if self.ssh:
            # We can terminate the SSH process
            self.ssh.process.process.terminate()

    def start(self):
        report("Checking if the ASM has already started...")
        is_synced = self.is_synced()
        if is_synced is not None:
            if is_synced:
                print("Alert System Microservice has already started")
                return
            else:
                print("Alert System Microservice has already started, waiting to be synced...")
                self.wait_for_synced()
                return
        print("Starting the Alert System Microservice...")
        self.run()

# Representing the BSV bitcoin-cli.
class BSVCLI():

    # Provide additional bitcoin-cli parameters if needed for RPC commands
    def __init__(self, args = [], ssh = None):
        self.command = ["bitcoin-cli"] + args
        self.ssh = ssh

    # Runs bitcoin-cli command locally and returns the result
    def run_command_locally(self, command):
        process = Process(command)
        return process.open()
    
    # Runs bitcoin-cli command remotely (SSH) and returns the result
    def run_command_ssh(self, command):
        return self.ssh.run(command)

    # Sends RPC and returns the result as a JSON object
    def rpc(self, rpc, rpc_args = []):
        command = self.command + [rpc] + rpc_args
        json_output = None
        if self.ssh:
            json_output = self.run_command_ssh(' '.join(command))
        else:
            json_output = self.run_command_locally(command)
        try:
            return json.loads(json_output, parse_float=decimal.Decimal)
        except json.JSONDecodeError as e:
            report(f"{command} returned:\n{json_output}")
            report_exception("Not a JSON string", e)
        return None

# Representing the BSV node.
class BSVNode():

    # Provide additional bitcoind parameters if needed to start the node properly
    def __init__(self, args = [], ssh_args = {}, timeout = 60):
        self.timeout = timeout
        self.command = ["bitcoind"] + args
        self.ssh = None
        if len(ssh_args) == 3:
            self.ssh = SSHCall(ssh_args.get("host"), ssh_args.get("user"), ssh_args.get("pk_path"))
        self.cli = BSVCLI(args=args, ssh=self.ssh)
        self.process = None

    def run_node_ssh(self):
        # We want to get the stderr to be able to report SSH issues
        self.ssh.run(' '.join(self.command), blocking=False, stderr=subprocess.PIPE)
        self.wait_for_node_ready(self.ssh.process)
        # We can terminate the SSH process
        self.ssh.process.process.terminate()
    
    def run_node_locally(self):
        process = Process(self.command)
        # We want to get the stderr to be able to report bitcoind issues
        process.open(blocking=False, stderr=subprocess.PIPE)
        self.wait_for_node_ready(process)
    
    def wait_for_node_ready(self, process):
        report("Waiting for RPC connection...")
        self.wait_for_rpc_connection(process)
        report("Waiting for node initialization...")
        self.wait_for_initialization()

    def wait_for_rpc_connection(self, process):
        running = False
        for _ in range(self.timeout):
            # Check if process (either ssh or bitcoind) already terminated while waiting for the RPC connection
            return_code = process.process.poll()
            if return_code is not None:
                command = self.command
                if self.ssh:
                    command = self.ssh.ssh_command + command
                stderr = None
                if process.process.stderr is not None:
                    stderr = process.process.stderr.read()
                raise subprocess.CalledProcessError(return_code, command, stderr=stderr)
            try:
                self.cli.rpc("getblockcount")
                # RPC connection is up
                running = True
                break
            except Exception as e:
                report_exception("Failed with", e)
                report("Retrying...")
                time.sleep(1.0)
        if not running:
            raise AssertionError("RPC connection timeout exceeded")
    
    def wait_for_initialization(self):
        wait_until = time.time() + self.timeout
        while time.time() < wait_until:
            if (self.cli.rpc("getinfo")["initcomplete"]):
                return
            time.sleep(1.0)
        assert wait_until >= time.time(), "Initialization timeout exceeded"

    def run_node(self):
        if self.ssh:
            self.run_node_ssh()
        else:
            self.run_node_locally()

    def start(self):
        try:
            report("Checking if the BSV node has already started...")
            self.cli.rpc("getblockcount")
            print("The BSV node has already started, waiting for initialization...")
            self.wait_for_initialization()
            return
        except Exception as e:
            # At this point we assume the node is not running yet
            report("BSV node has not started yet")
            report_exception("Exception was", e)
        print("Starting the BSV node...")
        self.run_node()

verbose = False
def report(message):
    if verbose:
        print(f"[{datetime.now().strftime('%d-%m-%Y %H:%M:%S.%f')}]: {message}")

def report_exception(message, exception, exit = False):
    error = None
    if hasattr(exception, 'stderr'):
        error = exception.stderr
    if exit:
        print(f"{message}: {exception}")
        if error is not None:
            print(f"Error: {error}")
        sys.exit(1)
    else:
        report(f"{message}: {exception}")
        if error is not None:
            report(f"Error: {error}")

def help(error = None):
    if error:
        print(f"{error}\n")
    print("Usage:\n" \
          "start_aks_bsv.py [-h[elp]] [-asm_port=PORT] [ASM SSH OPTIONS] [BSV SSH OPTIONS] [BSV OPTIONS] [-v[erbose]]\n\n" \
          "-h[elp]           Prints out this help message\n" \
          "-asm_port=PORT    Alert System Microservice HTTP port (3000 by default)\n" \
          "ASM SSH OPTIONS   SSH key-based authentication options to access the remote Alert System Microservice:\n" \
          "  -asm_host=HOST  IP or hostname of the remote Alert System Microservice\n" \
          "  -asm_user=USER  Username for the SSH connection\n" \
          "  -asm_pk_path=PK Private key file path\n" \
          "BSV SSH OPTIONS   SSH key-based authentication options to access the remote BSV node:\n" \
          "  -bsv_host=HOST  IP or hostname of the remote BSV node\n" \
          "  -bsv_user=USER  Username for the SSH connection\n" \
          "  -bsv_pk_path=PK Private key file path\n" \
          "BSV OPTIONS       Any additional bitcoind and bitcoin-cli parameters as -key or -key=value\n" \
          "-v[erbose]        Prints out details during the startup\n\n"
          "Example:\n" \
          "start_aks_bsv.py -datadir=/data/bsv -bsv_host=bsvhost.com -bsv_user=bsv_usr1 -bsv_pk_path=/home/bsv_usr1/.ssh/id_ed25519 -verbose")

def parse_arguments(*args):
    global verbose
    asm_args = {"ssh": {}}
    bsv_args = {"ssh": {}, "options": []}
    show_help = False
    # Argument can be either a -key or -key=value
    arg_pattern = re.compile(r'^-([a-z]+[a-z0-9_]*|[a-z]+[a-z0-9_]*=.+)$')
    for arg in args:
        if not bool(arg_pattern.match(arg))  :
            help(error=f"Error: Wrong argument {arg}")
            sys.exit(1)
        key = arg[1:]
        value = None
        if "=" in key:
            key, value = key.split('=', 1)
        # -v[erbose]
        if key == "v" or key == "verbose":
            verbose = True
        # -h[elp]
        elif key == "h" or key == "help":
            show_help = True
        # -asm_port
        elif key == "asm_port":
            if value is None:
                help(error=f"Error: {arg} is missing a value")
                sys.exit(1)
            asm_args[key[4:]] = value
        # ASM SSH OPTIONS
        elif key == "asm_host" or key == "asm_user" or key == "asm_pk_path":
            if value is None:
                help(error=f"Error: {arg} is missing a value")
                sys.exit(1)
            asm_args["ssh"][key[4:]] = value
        # BSV SSH OPTIONS
        elif key == "bsv_host" or key == "bsv_user" or key == "bsv_pk_path":
            if value is None:
                help(error=f"Error: {arg} is missing a value")
                sys.exit(1)
            bsv_args["ssh"][key[4:]] = value
        # BSV OPTIONS for everything else
        else:
           bsv_args["options"].append(arg) 
    report(f"Input parameters: {args}")
    if show_help:
        help()
        sys.exit(0)
    # With SSH, all options must be provided
    if len(asm_args.get("ssh")) > 0 and len(asm_args.get("ssh")) != 3:
        help(error="Error: Not all Alert System Microservice SSH parameters were provided.")
        sys.exit(1)
    if len(bsv_args.get("ssh")) > 0 and len(bsv_args.get("ssh")) != 3:
        help(error="Error: Not all BSV SSH parameters were provided.")
        sys.exit(1)
    return asm_args, bsv_args

def main():
    # Parse arguments
    asm_args, bsv_args = parse_arguments(*sys.argv[1:])
    # Start the Alert System Microservice
    asm = ASM(port=asm_args.get("port"), ssh_args=asm_args.get("ssh"))
    try:
        asm.start()
        print("Alert System Microservice is up and running")
    except Exception as e:
        report_exception("Failed to start the Alert System Microservice", e, exit=True)
    # Start the node
    bsv_node = BSVNode(args=bsv_args.get("options"), ssh_args=bsv_args.get("ssh"))
    try:
        bsv_node.start()
        print("BSV node is up and running")
    except Exception as e:
        report_exception("Failed to start the BSV node", e, exit=True)


if __name__ == "__main__":
    main()