# gRPC event streaming example

Simple [client](cmd/client/main.go) and [server](cmd/server/main.go) to demonstrate how a [gRPC](https://grpc.io/) server might stream events to multiple clients (based on this [proto spec](https://github.com/al-maisan/rgsproto/blob/main/proto/api.proto)).

# how to build

1. go mod tidy
1. make

# how to play with the code

1. run `bin/server` in a one terminal
1. run `bin/client client-1` in a separate terminal
1. run `bin/client client-2` in a separate terminal