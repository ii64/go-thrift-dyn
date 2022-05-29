package test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/davecgh/go-spew/spew"
	th "github.com/ii64/go-thrift-dyn"
	"github.com/ii64/go-thrift-dyn/internal/test/base"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"testing"
	"time"
)

const dmp = true
const pprofX = false

type readWriteCloser struct {
	*io.PipeReader
	*io.PipeWriter
	b      *bytes.Buffer
	readed *int
	mu     sync.Mutex
}

func (rw readWriteCloser) Write(b []byte) (n int, err error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	rw.b.Write(b)
	return rw.PipeWriter.Write(b)
}

func (rw readWriteCloser) Read(b []byte) (n int, err error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	if *rw.readed > rw.b.Len() {
		err = io.EOF
		return
	}
	n, err = rw.PipeReader.Read(b)
	if err != nil {
		return
	}
	*rw.readed = *rw.readed + n
	return
}

func newReadWriteCloser() readWriteCloser {
	reader, writer := io.Pipe()
	return readWriteCloser{
		PipeReader: reader,
		PipeWriter: writer,
		b:          &bytes.Buffer{},
		readed:     new(int),
	}
}
func (p readWriteCloser) Close() error {
	if err := p.PipeReader.Close(); err != nil {
		return err
	}
	if err := p.PipeWriter.Close(); err != nil {
		return err
	}
	return nil
}

//go:generate thriftgo -o . --gen go:gen_deep_equal=true,use_type_alias=true,typed_enum_string=false,json_enum_as_text=true,reorder_fields=true,scan_value_for_enum=true,frugal_tag=true,gen_db_tag=true model.thrift

func BenchmarkModelRebuildOriginal(bn *testing.B) {
	ctx := context.Background()
	var err error
	var b bytes.Buffer
	b.Grow(2 << 11)
	trans := thrift.NewStreamTransportRW(&b)
	protofactory := th.ProtocolFactory(th.ProtocolType_Compact, &thrift.TConfiguration{})
	prot := protofactory.GetProtocol(trans)

	listData := []int64{1, 2, 3}
	mapI64 := map[int64]int64{1: 2}
	mapI32 := map[int32]int32{3: 4}

	var m base.Model
	m.Abc = "hello"
	m.Sd = 0xcafe
	m.ListI64 = listData
	m.MapI64 = mapI64
	m.MapI32 = mapI32

	m.Write(ctx, prot)
	prot.Flush(ctx)
	fmt.Println(b.Bytes())

	m.MapI32 = map[int32]int32{}
	bn.ResetTimer()
	for i := 0; i < bn.N; i++ {
		b.Reset()
		m.Abc = "hello"
		m.Sd = 0xcafe
		m.ListI64 = listData
		m.MapI64 = mapI64
		// m.MapI32 = mapI32 // fair bench with dyn one.

		err = m.Write(ctx, prot)
		if err != nil {
			bn.Fail()
		}
		err = prot.Flush(ctx)
		if err != nil {
			bn.Fail()
		}
	}
	bn.StopTimer()
}

func BenchmarkModelRebuildDyn(bn *testing.B) {
	if pprofX {
		{
			f, err := os.Create(strconv.FormatInt(time.Now().Unix(), 10) + ".cpu")
			if err != nil {
				panic(err)
			}
			if err = pprof.StartCPUProfile(f); err != nil {
				panic(err)
			}
		}
		defer pprof.StopCPUProfile()
		defer func() {
			f, err := os.Create(strconv.FormatInt(time.Now().Unix(), 10) + ".heap")
			if err != nil {
				panic(err)
			}
			defer f.Close()
			if err = pprof.WriteHeapProfile(f); err != nil {
				panic(err)
			}
		}()
	}
	ctx := context.Background()
	var err error
	var b bytes.Buffer
	b.Grow(2 << 11)
	trans := thrift.NewStreamTransportRW(&b)
	protofactory := th.ProtocolFactory(th.ProtocolType_Compact, &thrift.TConfiguration{})
	prot := protofactory.GetProtocol(trans)

	listData := []int64{1, 2, 3}
	mapI64 := map[int64]int64{1: 2}
	mapI32 := map[int32]int32{3: 4}

	var m th.RPCStruct
	f1 := th.NewTField(1, thrift.STRING, "abc", true)
	f2 := th.NewTField(4, thrift.I64, "sd", true)
	f3 := th.NewTField(9, thrift.DOUBLE, "f64", true)
	f4 := th.NewTField(10, thrift.LIST, "listI64", true)
	f4Value := th.NewTypeContainerList[int64](th.TypeContainerDesc{Value: thrift.I64}, true)
	f4.SetValue(f4Value)
	f5 := th.NewTField(11, thrift.MAP, "mapI64", true)
	f5Value := th.NewTypeContainerMapUnordered[int64, int64](th.TypeContainerDesc{Key: thrift.I64, Value: thrift.I64}, true)
	f5.SetValue(f5Value)
	f6 := th.NewTField(12, thrift.MAP, "mapI32", false)
	f6Value := th.NewTypeContainerMap[int32, int32](th.TypeContainerDesc{Key: thrift.I32, Value: thrift.I32}, false)
	f6.SetValue(f6Value) // comment this to make it act as optional
	m.AddField(f1, f2, f3, f4, f5, f6)

	f1.SetValue("hello")
	f2.SetValue(0xcafe)
	f4Value.Value = listData
	f5Value.Value = mapI64
	f6Value.FromMap(mapI32)

	// check
	m.Write(ctx, prot)
	prot.Flush(ctx)
	fmt.Println(b.Bytes())
	bn.ResetTimer()
	for i := 0; i < bn.N; i++ {
		b.Reset()
		f1.SetValue("hello")
		f2.SetValue(0xcafe)
		f4Value.Value = listData
		f5Value.Value = mapI64
		// f6Value.FromMap(mapI32) // fair bench with dyn one.

		err = m.Write(ctx, prot)
		if err != nil {
			fmt.Println(err)
			bn.Fail()
		}
		err = prot.Flush(ctx)
		if err != nil {
			fmt.Println(err)
			bn.Fail()
		}
		// fmt.Println(b.Bytes())
	}
	bn.StopTimer()

}

func TestModelRebuild(t *testing.T) {
	ctx := context.Background()
	is := require.New(t)
	var err error
	var b bytes.Buffer
	trans := thrift.NewStreamTransportRW(&b)
	protofactory := th.ProtocolFactory(th.ProtocolType_Binary, &thrift.TConfiguration{})
	prot := protofactory.GetProtocol(trans)

	var m base.Model
	m.Abc = "hello"
	m.Sd = 0xcafe
	m.ListI64 = []int64{1, 2, 3}
	m.MapI64 = map[int64]int64{
		8: 9,
	}
	m.MapI32 = map[int32]int32{
		2: 8,
	}

	err = m.Write(ctx, prot)
	is.NoError(err)

	err = prot.Flush(ctx)
	is.NoError(err)

	var expected = make([]byte, b.Len())
	copy(expected, b.Bytes())
	b.Reset()

	// rebuild model

	var m2 th.RPCStruct
	f1 := th.NewTField(1, thrift.STRING, "abc", true)
	f2 := th.NewTField(4, thrift.I64, "sd", true)
	f3 := th.NewTField(9, thrift.DOUBLE, "f64", true)
	f4 := th.NewTField(10, thrift.LIST, "listI64", true)
	f4Value := th.NewTypeContainerList[int64](th.TypeContainerDesc{Value: thrift.I64}, true)
	f4.SetValue(f4Value)
	f5 := th.NewTField(11, thrift.MAP, "mapI64", true)
	f5Value := th.NewTypeContainerMapUnordered[int64, int64](th.TypeContainerDesc{Key: thrift.I64, Value: thrift.I64}, true)
	f5.SetValue(f5Value)
	f6 := th.NewTField(12, thrift.MAP, "mapI32", false)
	f6Value := th.NewTypeContainerMap[int32, int32](th.TypeContainerDesc{Key: thrift.I32, Value: thrift.I32}, false)
	f6.SetValue(f6Value) // comment this to make it act as optional
	m2.AddField(f1, f2, f3, f4, f5, f6)
	//
	f1.SetValue(m.Abc)
	f2.SetValue(m.Sd)
	f4Value.Value = m.ListI64
	f5Value.Value = m.MapI64
	f6Value.FromMap(m.MapI32)

	err = m2.Write(ctx, prot)
	is.NoError(err)

	err = prot.Flush(ctx)
	is.NoError(err)

	var actual = make([]byte, b.Len())
	copy(actual, b.Bytes())
	b.Reset()

	fmt.Printf("%+#v\n%+#v\n - f1: %+#v\n - f2: %+#v\n - f3: %+#v\n - f4: %+#v\n",
		m, m2, f1, f2, f3, f4)

	is.Equal(expected, actual)
}

func TestModelMapOnly(t *testing.T) {
	var err error
	var b bytes.Buffer
	trans := thrift.NewStreamTransportRW(&b)
	protofactory := th.ProtocolFactory(th.ProtocolType_Compact, &thrift.TConfiguration{})
	prot := protofactory.GetProtocol(trans)

	var m = base.MapOnly{
		ListById: map[int64][]int32{
			1: {1, 2, 3, 4, 5},
		},
		StringById: map[int64]string{
			123: "hello",
			923: "world",
		},
		ModelById: map[int64]*base.Model{
			886:  base.NewModel(),
			1314: base.NewModel(),
		},
	}
	err = m.Write(context.Background(), prot)
	require.NoError(t, err)
	err = prot.Flush(context.Background())
	require.NoError(t, err)

	prevB := make([]byte, b.Len())
	copy(prevB, b.Bytes())
	fmt.Println(prevB)

	var m2 th.RPCStruct
	err = m2.Read(context.Background(), prot)
	require.NoError(t, err)

	if dmp {
		spew.Dump(m2)
	}

	b.Reset()
	err = m2.Write(context.Background(), prot)
	require.NoError(t, err)
	err = prot.Flush(context.Background())
	postB := make([]byte, b.Len())
	copy(postB, b.Bytes())
	fmt.Println(postB)
}

func TestModelListBytes(t *testing.T) {
	var err error
	var b bytes.Buffer
	trans := thrift.NewStreamTransportRW(&b)
	protofactory := th.ProtocolFactory(th.ProtocolType_Compact, &thrift.TConfiguration{})
	prot := protofactory.GetProtocol(trans)

	var m = base.ListOnly{
		Names: []string{"hello", "world"},
	}
	err = m.Write(context.Background(), prot)
	require.NoError(t, err)
	err = prot.Flush(context.Background())
	require.NoError(t, err)

	prevB := b.Bytes()
	fmt.Println(prevB)

	var m2 th.RPCStruct
	err = m2.Read(context.Background(), prot)
	require.NoError(t, err)

	if dmp {
		spew.Dump(m2)
	}

	b.Reset()
	err = m2.Write(context.Background(), prot)
	require.NoError(t, err)
	err = prot.Flush(context.Background())
	require.NoError(t, err)
	postB := b.Bytes()
	fmt.Println(postB)
}

func TestModelReadWireData(t *testing.T) {
	prots := []th.ProtocolType{
		th.ProtocolType_Binary,
		th.ProtocolType_Compact,
	}
	var wg sync.WaitGroup
	wg.Add(len(prots))
	for _, protoName := range prots {
		protofactory := th.ProtocolFactory(protoName, &thrift.TConfiguration{})
		t.Run(protoName, func(t *testing.T) {
			var err error

			s2cPipe := newReadWriteCloser()
			c2sPipe := newReadWriteCloser()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// echo back
			go func() {
				defer wg.Done()

				itrans := thrift.NewStreamTransportRW(c2sPipe)
				otrans := thrift.NewStreamTransportRW(s2cPipe)
				iprot := protofactory.GetProtocol(itrans)
				oprot := protofactory.GetProtocol(otrans)
				// srv := thrift.NewTStandardClient(nil, nil)
				var (
					// m   base.Request
					m   th.RPCStruct
					err error
				)
				for ctx.Err() == nil {
					var name string
					var seqid int32
					name, _, seqid, err = iprot.ReadMessageBegin(ctx)
					err = m.Read(ctx, iprot)
					require.NoError(t, err)
					err = iprot.ReadMessageEnd(ctx)
					require.NoError(t, err)

					expectedBytes := c2sPipe.b.Bytes()
					fmt.Println("c2s", expectedBytes)

					err = oprot.WriteMessageBegin(ctx, name, thrift.REPLY, seqid)
					require.NoError(t, err)
					err = m.Write(ctx, oprot)
					require.NoError(t, err)
					err = oprot.WriteMessageEnd(ctx)
					require.NoError(t, err)
					err = otrans.Flush(ctx)
					require.NoError(t, err)

					actualBytes := s2cPipe.b.Bytes()
					fmt.Println("s2c", actualBytes)

					// skipping TMessageType
					switch protofactory.(type) {
					case *thrift.TCompactProtocolFactory:
						start := 2 // 2nd byte
						if !bytes.Equal(expectedBytes[start:], actualBytes[start:]) {
							fmt.Println(expectedBytes)
							fmt.Println(actualBytes)
							panic("not equal")
						}
					case *thrift.TBinaryProtocolFactory:
						start := 4 // int32( VERSION | TMessageType )
						fmt.Println(expectedBytes)
						fmt.Println(actualBytes)
						if !bytes.Equal(expectedBytes[start:], actualBytes[start:]) {
							panic("not equal")
						}
					}

					if dmp {
						spew.Dump(m)
					}
				}
			}()

			itrans := thrift.NewStreamTransportRW(s2cPipe)
			defer itrans.Close()
			otrans := thrift.NewStreamTransportRW(c2sPipe)
			defer otrans.Close()

			iprot := protofactory.GetProtocol(itrans)
			oprot := protofactory.GetProtocol(otrans)

			var (
				expected = base.Request{
					Model: base.NewModel(),
					Models: []*base.Model{
						base.NewModel(),
					},
					ModelById: map[int64]*base.Model{
						1234: base.NewModel(),
						// 0xff:   nil, // TODO: not nullable?
					},
					ModelByTime: map[int64][]*base.Model{
						567: {base.NewModel()},
					},
					Modset: nil,
				}
				actual base.Request
			)
			// var (
			// 	expected = base.SimpleListMap{
			// 		ListListI32: [][]int32{
			// 			{1, 2, 3, 4},
			// 			{5, 6, 7, 8},
			// 		},
			// 		ListI64:  []int64{1, 2, 3},
			// 		I64ByI64: map[int64]int64{4: 5},
			// 	}
			// 	actual base.SimpleListMap
			// )

			client := thrift.NewTStandardClient(iprot, oprot)
			_, err = client.Call(context.Background(), "testRpcMethod", &expected, &actual)
			require.NoError(t, err)

			require.Equal(t, expected, actual)
		})
	}
	wg.Wait() // wait srv goroutine end.
}
