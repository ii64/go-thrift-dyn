package thrift_dyn

import (
	"bytes"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/stretchr/testify/require"
)

var defaultTestTConfiguration = &thrift.TConfiguration{}
var defaultTestTProtocols = []ProtocolType{
	ProtocolType_Compact,
	ProtocolType_Binary,
}

func withBytesTTransport(t require.TestingT, f func(b *bytes.Buffer, trans thrift.TTransport)) {
	var b bytes.Buffer
	trans := thrift.NewStreamTransportRW(&b)
	defer trans.Close()
	f(&b, trans)
}

func withProtocols(t require.TestingT, protoList []ProtocolType, tcfg *thrift.TConfiguration, f func(protofactory thrift.TProtocolFactory)) {
	if tcfg == nil {
		tcfg = defaultTestTConfiguration
	}
	for _, proto := range protoList {
		protofactory := ProtocolFactory(proto, tcfg)
		f(protofactory)
	}
}
