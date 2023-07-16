# potee-tasks-checker

## generate proto
```
protoc --go-grpc_out=. --go_out=. proto/data.proto 
```

## Check
```
grpcurl -import-path proto -proto data.proto -plaintext \
    -d '{"address":"localhost", "commands":[{"command":"python3 examples/checker.py", "input":"test"}]}' \
    localhost:50051 Checker/Run
```