
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
    1: string abc; 4: i64 sd; 9: double f64 // "hello", 0xcafe, 0.0
    10: optional list<i64> listI64    // {1, 2, 3}
    11: optional map<i64, i64> mapI64 // {1: 2}
    12: optional map<i32, i32> mapI32 // empty
}
```

```
TBinary:
BenchmarkModelRebuildOriginal-12     1971684	       611.0 ns/op
BenchmarkModelRebuildOriginal-12     1982988	       598.5 ns/op
BenchmarkModelRebuildOriginal-12     1901395	       597.0 ns/op
BenchmarkModelRebuildDyn-12    	     1000000	      1040 ns/op
BenchmarkModelRebuildDyn-12    	     1000000	      1063 ns/op
BenchmarkModelRebuildDyn-12    	     1000000	      1025 ns/op

TCompact:
BenchmarkModelRebuildOriginal-12     1681024	       748.5 ns/op
BenchmarkModelRebuildOriginal-12     1709762	       690.0 ns/op
BenchmarkModelRebuildOriginal-12     1666425	       769.9 ns/op
BenchmarkModelRebuildDyn-12          1124538	      1177 ns/op
BenchmarkModelRebuildDyn-12          1000000	      1091 ns/op
BenchmarkModelRebuildDyn-12          1063579	      1171 ns/op
```