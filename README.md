# synse-juniper-jti-plugin

Synse plugin for Juniper JTI metrics over UDP stream

> **Note**: This project is currently in development



## An aside on (re)compiling the .proto source files

I either do not know enough about how `protoc` works, or am just otherwise struggling
to get each source file to compile to a package correctly.

At present, running `$ ./scripts/gen_proto.sh` will generate the compiled Go source
for the proto files. The issue lies with the import paths.  The generated Go source
is put into `pkg/protocol/jti/protos/{name}/{name}.pb.go`. All files (except
`telemetry_top`) import `telemetry_top`, but the import is wrong as generated. As such,
it needs to be updated from

```go
import (
	telemetry_top "protos/telemetry_top"
)
``` 

to

```go
import (
	telemetry_top "github.com/vapor-ware/synse-juniper-jti-plugin/pkg/protocol/jti/protos/telemetry_top"
)
``` 

Note also that the `.proto` files should have a `go_import` option added, e.g. for `port.proto`:

```proto
option go_package = "protos/port";
```
