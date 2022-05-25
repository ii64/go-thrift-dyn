package thrift_dyn

import "github.com/apache/thrift/lib/go/thrift"

type ProtocolType = string

const (
	ProtocolType_Compact    ProtocolType = "tcompact"
	ProtocolType_Binary     ProtocolType = "tbinary"
	ProtocolType_SimpleJSON ProtocolType = "tsimplejson"
	ProtocolType_JSON       ProtocolType = "tjson"
)

var ProtocolType_VALUES = []ProtocolType{
	ProtocolType_Compact,
	ProtocolType_Binary,
	ProtocolType_SimpleJSON,
	ProtocolType_JSON,
}

func ProtocolFactory(ptype ProtocolType, conf *thrift.TConfiguration) thrift.TProtocolFactory {
	switch ptype {
	case ProtocolType_Compact:
		return thrift.NewTCompactProtocolFactoryConf(conf)
	case ProtocolType_Binary:
		return thrift.NewTBinaryProtocolFactoryConf(conf)
	case ProtocolType_SimpleJSON:
		return thrift.NewTSimpleJSONProtocolFactoryConf(conf)
	case ProtocolType_JSON:
		return thrift.NewTJSONProtocolFactory()
	default:
		panic("unsupported protocol type")
	}
}
