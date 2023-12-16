# Info
```
go run publish.go -type=1 -sequence=1
```

# Ban peer
```
go run publish.go -type=5 -sequence=2 -peer=192.168.1.2/24
```

# Unban peer
```
go run publish.go -type=6 -sequence=3 -peer=192.168.1.2/24
```

# Invalidate Block
```
go run publish.go -type=7 -block-hash=<hash> -sequence=4
```

# Set Keys
Note: the keys should be compressed public keys. If `-signing-keys` is not set,
this message will be signed with the genesis keys.
```
go run publish.go -type=8 -sequence=5 -pub-keys=<key1>,<key2>,<key3>,<key4>,<key5>
```

# Test alert with different signing keys
```
go run publish.go -type=1 -sequence=6 -signing-keys=<key1>,<key2>,<key3>
```

