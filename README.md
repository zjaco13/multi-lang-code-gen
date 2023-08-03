# Multi Lang Code Gen
Protobuf boilerplate generation framework for EKS Blueprints

## How to run

Build into binary then run:
```bash
go build
./multi-lang-code-gen
```

Run with go:
```bash
go run main.go
```

### Optional Flags
```bash
-obj        specify text file containing objects to generate.  Defaults to obj.txt

-addon      specify proto file to write new addons to, must already exist.  Defaults to addons.proto

-clp        specify proto file to write new cluster providers to, must already exist.  Defaults to cluster_providers.proto

-rp         specify proto file to write new resource providers to, must already exist.  Defaults to resource_providers.proto

-team       specify proto file to write new teams to, must already exist.  Defaults to teams.proto

-rpc        specify proto file to write new rpc calls to, must already exist.  Defaults to cluster.proto
```

Running this program will generate a new rpc file at the passed in rpcfile.temp or cluster.proto.temp as default.  You must then overwrite the original file with an external command.

Example overwrite command:
```bash
mv cluster.proto.temp cluster.proto
```

