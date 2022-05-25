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

const dmp = false
const pprofX = false

type readWriteCloser struct {
	*io.PipeReader
	*io.PipeWriter
	b      *bytes.Buffer
	readed *int
}

func (rw readWriteCloser) Write(b []byte) (n int, err error) {
	rw.b.Write(b)
	return rw.PipeWriter.Write(b)
}

func (rw readWriteCloser) Read(b []byte) (n int, err error) {
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
	trans := thrift.NewStreamTransportRW(&b)
	protofactory := th.ProtocolFactory(th.ProtocolType_Compact, &thrift.TConfiguration{})
	prot := protofactory.GetProtocol(trans)

	var m base.Model
	m.Write(ctx, prot)
	prot.Flush(ctx)
	fmt.Println(b.Bytes())
	bn.ResetTimer()
	for i := 0; i < bn.N; i++ {
		b.Reset()
		m.Abc = "hello"
		m.Sd = 0xcafe
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
	trans := thrift.NewStreamTransportRW(&b)
	protofactory := th.ProtocolFactory(th.ProtocolType_Compact, &thrift.TConfiguration{})
	prot := protofactory.GetProtocol(trans)

	var m th.RPCStruct
	f1 := th.NewTField(1, thrift.STRING, "abc", true)
	f2 := th.NewTField(4, thrift.I64, "sd", true)
	f3 := th.NewTField(9, thrift.DOUBLE, "f64", true)
	m.AddField(f1)
	m.AddField(f2)
	m.AddField(f3)
	// check
	m.Write(ctx, prot)
	prot.Flush(ctx)
	fmt.Println(b.Bytes())
	bn.ResetTimer()
	for i := 0; i < bn.N; i++ {
		b.Reset()
		f1.SetValue("hello")
		f2.SetValue(0xcafe)
		err = m.Write(ctx, prot)
		if err != nil {
			fmt.Println(err)
			bn.Fail()
		}
		err = prot.Flush(ctx)
		if err != nil {
			bn.Fail()
		}
	}
	bn.StopTimer()

}

func TestModelRebuild(t *testing.T) {
	ctx := context.Background()
	is := require.New(t)
	var err error
	var b bytes.Buffer
	trans := thrift.NewStreamTransportRW(&b)
	protofactory := th.ProtocolFactory(th.ProtocolType_Compact, &thrift.TConfiguration{})
	prot := protofactory.GetProtocol(trans)

	var m base.Model
	err = m.Write(ctx, prot)
	is.NoError(err)

	err = prot.Flush(ctx)
	is.NoError(err)

	var expected = make([]byte, b.Len())
	copy(expected, b.Bytes())
	b.Reset()

	// rebuild model

	var m2 th.RPCStruct
	m2.AddField(th.NewTField(1, thrift.STRING, "abc", true))
	m2.AddField(th.NewTField(4, thrift.I64, "sd", true))
	m2.AddField(th.NewTField(9, thrift.DOUBLE, "f64", true))
	err = m2.Write(ctx, prot)
	is.NoError(err)

	err = prot.Flush(ctx)
	is.NoError(err)

	var actual = make([]byte, b.Len())
	copy(actual, b.Bytes())
	b.Reset()

	is.Equal(expected, actual)
}

func TestModelMapOnly(t *testing.T) {
	var err error
	var b bytes.Buffer
	trans := thrift.NewStreamTransportRW(&b)
	protofactory := th.ProtocolFactory(th.ProtocolType_Compact, &thrift.TConfiguration{})
	prot := protofactory.GetProtocol(trans)

	var m = base.MapOnly{
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
		th.ProtocolType_Compact,
		th.ProtocolType_Binary,
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
							panic("not equal")
						}
					case *thrift.TBinaryProtocolFactory:
						start := 4 // int32( VERSION | TMessageType )
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
						0xcafe: base.NewModel(),
						// 0xff:   nil, // TODO: not nullable?
					},
					Modset: nil,
				}
				actual base.Request
			)
			// var (
			// 	expected = base.SimpleListMap{
			// 		ListI64:  []int64{1, 2, 3},
			// 		I64ByI64: map[int64]int64{3: 3},
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
