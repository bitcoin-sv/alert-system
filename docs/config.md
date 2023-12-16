# Alert System Configuration

| Parameter                      | Default Value                         | Description                                         |
|--------------------------------|---------------------------------------|-----------------------------------------------------|
| alert_webhook_url              | ""                                    | URL for alert webhook notifications                 |
| request_logging                | true                                  | Enable or disable request logging                   |
| alert_processing_interval      | "5m"                                  | Interval for alert processing                       |
| environment                    | "local"                               | Environment setting (e.g., local, production)       |
| **web_server**                 | `<Object>`                            | Nested configuration for the web server             |
| web_server.idle_timeout        | "60s"                                 | Idle timeout for the web server                     |
| web_server.port                | "3000"                                | Port on which the web server listens                |
| web_server.read_timeout        | "15s"                                 | Read timeout for the web server                     |
| web_server.write_timeout       | "15s"                                 | Write timeout for the web server                    |
| **datastore**                  | `<Object>`                            | Configuration for the datastore                     |
| datastore.auto_migrate         | true                                  | Automatically migrate the datastore                 |
| datastore.debug                | true                                  | Enable or disable debugging for the datastore       |
| datastore.engine               | "sqlite"                              | Database engine (e.g., sqlite, postgresql)          |
| datastore.password             | ""                                    | Password for the database                           |
| datastore.table_prefix         | "alert_system"                        | Prefix for database table names                     |
| **datastore.sqlite**           | `<Object>`                            | SQLite specific configuration                       |
| datastore.sqlite.database_path | "alert_system_datastore.db"           | Path to the SQLite database file                    |
| datastore.sqlite.shared        | false                                 | Use a shared SQLite database                        |
| **sql_read**                   | `<Object>`                            | Configuration for the read SQL database connection  |
| **sql_write**                  | `<Object>`                            | Configuration for the write SQL database connection |
| sql_read/write.driver          | "postgresql"                          | Database driver (e.g., postgresql)                  |
| sql_read/write.host            | "localhost"                           | Hostname for the database server                    |
| ...                            |                                       | (Additional SQL read/write parameters)              |
| **p2p**                        | `<Object>`                            | P2P network configuration                           |
| p2p.ip                         | "0.0.0.0"                             | IP address for P2P communication                    |
| p2p.port                       | "9906"                                | Port for P2P communication                          |
| p2p.alert_system_protocol_id   | "/bitcoin-testnet/alert-system/0.0.1" | Protocol ID for the alert system on the P2P network |
| ...                            |                                       | (Additional P2P parameters)                         |
| **rpc_connections**            | `[]<Object>`                          | List of RPC connections                             |
| rpc_connections[0].user        | "testUser"                            | RPC username                                        |
| rpc_connections[0].password    | "testPw"                              | RPC password                                        |
| rpc_connections[0].host        | "http://localhost:8333"               | RPC host                                            |
