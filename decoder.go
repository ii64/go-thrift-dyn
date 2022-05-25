package thrift_dyn

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"io"
	"sync"
)

type Decoder struct {
	buf   bytes.Buffer
	prot  thrift.TProtocol
	trans *thrift.StreamTransport
	mu    sync.Mutex
}

func NewDecoder(pf thrift.TProtocolFactory) *Decoder {
	return (&Decoder{}).Init(pf)
}

func (dec *Decoder) Init(pf thrift.TProtocolFactory) *Decoder {
	dec.trans = thrift.NewStreamTransportR(&dec.buf)
	dec.prot = pf.GetProtocol(dec.trans)
	dec.buf.Reset()
	return dec
}

func (dec *Decoder) decodeInternal(reader io.Reader, valueDst any) (err error) {
	dec.trans.Reader = reader
	switch value := valueDst.(type) {
	case thrift.TStruct:
		err = value.Read(context.Background(), dec.prot)
	default:
		err = fmt.Errorf("unsupported type %T", value)
	}
	if err != nil {
		return
	}
	return
}

func (dec *Decoder) ReadFrom(reader io.Reader, valueDst any) (err error) {
	dec.mu.Lock()
	defer dec.mu.Unlock()
	if err = dec.decodeInternal(reader, valueDst); err != nil {
		return
	}
	return
}

func (dec *Decoder) Decode(src []byte, valueDst any) (err error) {
	dec.mu.Lock()
	defer dec.mu.Unlock()
	if err = dec.decodeInternal(bytes.NewReader(src), valueDst); err != nil {
		return
	}
	return
}
