
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

Benchmark write of simple message (sized 32 bytes) 

```
TBinary:
BenchmarkModelRebuildOriginal-12     5368544	       224.8 ns/op
BenchmarkModelRebuildOriginal-12     5474109	       230.8 ns/op
BenchmarkModelRebuildOriginal-12     5214946	       217.6 ns/op
BenchmarkModelRebuildDyn-12    	     4490389	       251.9 ns/op
BenchmarkModelRebuildDyn-12    	     4143368	       253.0 ns/op
BenchmarkModelRebuildDyn-12    	     4454946	       260.0 ns/op

TCompact:
BenchmarkModelRebuildOriginal-12     4522177	       253.9 ns/op
BenchmarkModelRebuildOriginal-12     4344098	       261.9 ns/op
BenchmarkModelRebuildOriginal-12     3954589	       259.0 ns/op
BenchmarkModelRebuildDyn-12    	     4282292	       275.9 ns/op
BenchmarkModelRebuildDyn-12    	     4444050	       279.5 ns/op
BenchmarkModelRebuildDyn-12    	     4120533	       280.7 ns/op
```