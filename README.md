# Alert System Microservice

This is the codebase for an alert system microservice that runs with alongside a Bitcoin SV Node and will produce automated RPC calls when validly signed alerts are received.

## Getting Started
### Copy settings file
```
$ cp example_settings_local.conf settings_local.conf
```

### Run the server (from source)
```
$ go run cmd/main.go
```

### Running with docker:
```
$ touch <db_file>
$ docker run -u root -e P2P_ALERT_SYSTEM_PROTOCOL_ID=/bitcoin/alert-system/0.0.1 -e P2P_BOOTSTRAP_PEER=/ip4/68.183.57.231/tcp/9906/p2p/12D3KooWQs6ptKvoKNHurCzqRaVp3uFs9731NQwS3AmVcNc2TGpb -e P2P_PORT=9906 -e P2P_IP=0.0.0.0 -v <database_file>:/alert_system_datastore.db docker.io/galtbv/alert-system:0.0.1
```

### Running with Podman
```
$ touch <db_file>
$ podman run -u root -e P2P_ALERT_SYSTEM_PROTOCOL_ID=/bitcoin/alert-system/0.0.1 -e P2P_BOOTSTRAP_PEER=/ip4/68.183.57.231/tcp/9906/p2p/12D3KooWQs6ptKvoKNHurCzqRaVp3uFs9731NQwS3AmVcNc2TGpb -e P2P_PORT=9906 -e P2P_IP=0.0.0.0 -v /git/alert-system/foo.db:/alert_system_datastore.db:Z docker.io/galtbv/alert-system:0.0.1
```
