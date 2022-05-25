package thrift_dyn

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"io"
	"sync"
)

type Encoder struct {
	buf   bytes.Buffer
	prot  thrift.TProtocol
	trans *thrift.StreamTransport
	mu    sync.Mutex
}

func NewEncoder(pf thrift.TProtocolFactory) *Encoder {
	return (&Encoder{}).Init(pf)
}

func (enc *Encoder) Init(pf thrift.TProtocolFactory) *Encoder {
	enc.trans = thrift.NewStreamTransportRW(&enc.buf)
	enc.prot = pf.GetProtocol(enc.trans)
	enc.buf.Reset()
	return enc
}

func (enc *Encoder) encodeInternal(value any) (err error) {
	enc.buf.Reset()
	switch value := value.(type) {
	case thrift.TStruct:
		err = value.Write(context.Background(), enc.prot)
	default:
		err = fmt.Errorf("uns"+
			"unsupported type %T", value)
	}
	if err != nil {
		return
	}
	err = enc.prot.Flush(context.Background())
	return
}

func (enc *Encoder) WriteTo(writer io.Writer, value any) (n int64, err error) {
	enc.mu.Lock()
	defer enc.mu.Unlock()
	err = enc.encodeInternal(value)
	if err != nil {
		return
	}
	n, err = enc.buf.WriteTo(writer)
	return
}

func (enc *Encoder) EncodeTo(dst []byte, value any) (n int, err error) {
	enc.mu.Lock()
	defer enc.mu.Unlock()
	err = enc.encodeInternal(value)
	if err != nil {
		return
	}
	n = copy(dst, enc.buf.Bytes())
	return
}

func (enc *Encoder) Encode(value any) (bb []byte, err error) {
	enc.mu.Lock()
	defer enc.mu.Unlock()
	err = enc.encodeInternal(value)
	if err != nil {
		return
	}
	bb = make([]byte, enc.buf.Len())
	copy(bb, enc.buf.Bytes())
	return
}
