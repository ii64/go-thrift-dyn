
# go-thrift-dyn

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/ii64/go-thrift-dyn.svg)](https://pkg.go.dev/github.com/ii64/go-thrift-dyn)

Inspect and rebuild Apache Thrift wire data



### Example

```go
// testRpcMethod_args
var args RPCStruct
// testRpcMethod_results
var results RPCStruct
client := thrift.NewTStandardClient(iprot, oprot)

var meta thrift.ResponseMeta
var err error
// oneway
meta, err = client.Call(ctx, "testRpcMethod", &args, nil)
// duplex
meta, err = client.Call(ctx, "testRpcMethod", &args, &results)
_ = results.Fields
```

### Benchmark

Benchmark write of simple message:

```thrift
struct Model {
  1: string abc; 4: i64 sd; 9: double f64
  10: optional list<i64> listI64
}
```

```
TBinary:
BenchmarkModelRebuildOriginal-12     3283041	       341.8 ns/op
BenchmarkModelRebuildOriginal-12     3502664	       345.2 ns/op
BenchmarkModelRebuildOriginal-12     3450882	       358.4 ns/op
BenchmarkModelRebuildDyn-12          2934448	       418.1 ns/op
BenchmarkModelRebuildDyn-12          2927372	       418.8 ns/op
BenchmarkModelRebuildDyn-12    	     2921905	       425.2 ns/op

TCompact:
BenchmarkModelRebuildOriginal-12     2822864	       436.8 ns/op
BenchmarkModelRebuildOriginal-12     2777244	       413.8 ns/op
BenchmarkModelRebuildOriginal-12     2831146	       439.3 ns/op
BenchmarkModelRebuildDyn-12    	     2587606	       469.5 ns/op
BenchmarkModelRebuildDyn-12    	     2506990	       478.6 ns/op
BenchmarkModelRebuildDyn-12    	     2543522	       468.8 ns/op
```