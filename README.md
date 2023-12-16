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

TODO: Running with docker:
