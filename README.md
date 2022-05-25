
# go-thrift-dyn

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
BenchmarkModelRebuildOriginal-12    5177924	       232.8 ns/op
BenchmarkModelRebuildOriginal-12    5016391	       225.4 ns/op
BenchmarkModelRebuildOriginal-12    5406008	       221.7 ns/op
BenchmarkModelRebuildDyn-12    	    2428218	       485.1 ns/op
BenchmarkModelRebuildDyn-12    	    2289037	       500.8 ns/op
BenchmarkModelRebuildDyn-12         2549248	       545.1 ns/op

TCompact:
BenchmarkModelRebuildOriginal-12    4135563	       249.8 ns/op
BenchmarkModelRebuildOriginal-12    4033410	       250.9 ns/op
BenchmarkModelRebuildOriginal-12    4571152	       252.9 ns/op
BenchmarkModelRebuildDyn-12         2381887	       650.6 ns/op
BenchmarkModelRebuildDyn-12         2359174	       525.4 ns/op
BenchmarkModelRebuildDyn-12    	    2321109	       571.8 ns/op

```