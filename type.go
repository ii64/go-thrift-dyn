package thrift_dyn

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type TDataSpec struct {
	Type     thrift.TType
	Required bool
	Protocol thrift.TProtocol
}

type TData[T any] struct {
	TDataSpec
	Value *T
}

func WriteDataGenericx(ctx context.Context, t TDataSpec, value any) (err error) {
	switch value := value.(type) {
	case bool:
		return t.Protocol.WriteBool(ctx, value)
	case byte:
		return t.Protocol.WriteByte(ctx, int8(value))
	case int8:
		return t.Protocol.WriteByte(ctx, value)
	case int16:
		return t.Protocol.WriteI16(ctx, value)
	case int32:
		return t.Protocol.WriteI32(ctx, value)
	case int64:
		return t.Protocol.WriteI64(ctx, value)
	case float64:
		return t.Protocol.WriteDouble(ctx, value)
	case []byte:
		return t.Protocol.WriteBinary(ctx, value)
	case string:
		return t.Protocol.WriteString(ctx, value)
	case thrift.TStruct: // *RPCStruct | TypeContainerImplementer
		return value.Write(ctx, t.Protocol)
	}
	return fmt.Errorf("expected %s, got %T", t.Type, value)
}

func WriteDataGeneric(ctx context.Context, t TDataSpec, value any) (err error) {
	switch t.Type {
	case thrift.BOOL:
		if value, ok := value.(bool); ok {
			return t.Protocol.WriteBool(ctx, value)
		} else if t.Required {
			return t.Protocol.WriteBool(ctx, false)
		}
	case thrift.BYTE:
		switch value := value.(type) {
		case int8:
			return t.Protocol.WriteByte(ctx, value)
		case byte:
			return t.Protocol.WriteByte(ctx, int8(value))
		}
		if t.Required {
			return t.Protocol.WriteByte(ctx, 0)
		}
	case thrift.I16:
		if value, ok := value.(int16); ok {
			return t.Protocol.WriteI16(ctx, value)
		}
		if t.Required {
			return t.Protocol.WriteI16(ctx, 0)
		}
	case thrift.I32:
		if value, ok := value.(int32); ok {
			return t.Protocol.WriteI32(ctx, value)
		}
		if t.Required {
			return t.Protocol.WriteI32(ctx, 0)
		}
	case thrift.I64:
		switch value := value.(type) {
		case int64:
			return t.Protocol.WriteI64(ctx, value)
		case int:
			return t.Protocol.WriteI64(ctx, int64(value))
		}
		if t.Required {
			return t.Protocol.WriteI64(ctx, 0)
		}
	case thrift.DOUBLE:
		if value, ok := value.(float64); ok {
			return t.Protocol.WriteDouble(ctx, value)
		}
		if t.Required {
			return t.Protocol.WriteDouble(ctx, 0)
		}
	case thrift.STRING:
		switch value := value.(type) {
		case string:
			return t.Protocol.WriteBinary(ctx, String2bs(value))
		case []byte:
			return t.Protocol.WriteBinary(ctx, value)
		}
		if t.Required {
			return t.Protocol.WriteBinary(ctx, []byte{})
		}
	case thrift.STRUCT:
		if value, ok := value.(thrift.TStruct); ok {
			return value.Write(ctx, t.Protocol)
		}
		if t.Required {
			return (&RPCStruct{}).Write(ctx, t.Protocol)
		}
		return nil
	case thrift.MAP, thrift.SET, thrift.LIST:
		if value, ok := value.(TypeContainerImplementer); ok {
			return value.Write(ctx, t.Protocol)
		}
		return nil
	}
	return fmt.Errorf("expected %s, got %T", t.Type, value)
}

func WriteData[T any](ctx context.Context, t TData[T]) (err error) {
	if t.Value == nil {
		return
	}
	return WriteDataGeneric(ctx, t.TDataSpec, *t.Value)
}

func ReadDataGeneric(ctx context.Context, t TDataSpec, value *any) (err error) {
	if value == nil {
		return
	}
	if value, ok := (*value).(thrift.TStruct); ok { // cache hit
		return value.Read(ctx, t.Protocol)
	}
	switch t.Type {
	case thrift.STOP, thrift.VOID:
		// nop.
	case thrift.BOOL:
		*value, err = t.Protocol.ReadBool(ctx)
		return
	case thrift.BYTE:
		*value, err = t.Protocol.ReadByte(ctx)
		return
	case thrift.I16:
		*value, err = t.Protocol.ReadI16(ctx)
		return
	case thrift.I32:
		*value, err = t.Protocol.ReadI32(ctx)
		return
	case thrift.I64:
		*value, err = t.Protocol.ReadI64(ctx)
	case thrift.DOUBLE:
		*value, err = t.Protocol.ReadDouble(ctx)
	case thrift.STRING:
		*value, err = t.Protocol.ReadString(ctx)
		return
	case thrift.STRUCT:
		st := &RPCStruct{}
		err = st.Read(ctx, t.Protocol)
		if err != nil {
			return
		}
		*value = st
		return
	case thrift.MAP:
		var (
			desc TypeContainerDesc
			size int
		)
		if desc.Key, desc.Value, size, err = t.Protocol.ReadMapBegin(ctx); err != nil {
			return
		}
		// if size > 0 {
		var typ TypeContainerImplementer
		typ, err = NewTypeContainerMapOfTType(desc, t.Required)
		if err != nil {
			return
		}
		typ.SetSize(size)
		if err = typ.Read(ctx, t.Protocol); err != nil {
			return
		}
		*value = typ
		// }
		if err = t.Protocol.ReadMapEnd(ctx); err != nil {
			return
		}
		return
	case thrift.SET:
		var (
			elemType thrift.TType
			size     int
		)
		if elemType, size, err = t.Protocol.ReadSetBegin(ctx); err != nil {
			return
		}
		// if size > 0 {
		var typ TypeContainerImplementer
		typ, err = NewTypeContainerSetOfTType(TypeContainerDesc{Value: elemType}, t.Required)
		if err != nil {
			return
		}
		typ.SetSize(size)
		if err = typ.Read(ctx, t.Protocol); err != nil {
			return
		}
		*value = typ
		// }
		if err = t.Protocol.ReadSetEnd(ctx); err != nil {
			return
		}
		return
	case thrift.LIST:
		var (
			elemType thrift.TType
			size     int
		)
		if elemType, size, err = t.Protocol.ReadListBegin(ctx); err != nil {
			return
		}
		// if size > 0 {
		var typ TypeContainerImplementer
		typ, err = NewTypeContainerListOfTType(TypeContainerDesc{Value: elemType}, t.Required)
		if err != nil {
			return
		}
		typ.SetSize(size)
		if err = typ.Read(ctx, t.Protocol); err != nil {
			return
		}
		*value = typ
		// }
		if err = t.Protocol.ReadListEnd(ctx); err != nil {
			return
		}
		return

	default:
		return fmt.Errorf("expected %s, got %T", t.Type, *value)
	}
	return
}

func ReadData[T any](ctx context.Context, t TData[T]) (err error) {
	if t.Value == nil {
		return
	}
	if value, ok := ((any)(*t.Value)).(thrift.TStruct); ok { // cache hit
		return value.Read(ctx, t.Protocol)
	}
	switch value := (any)(t.Value).(type) {
	case *bool:
		*value, err = t.Protocol.ReadBool(ctx)
		return
	case *byte:
		var tmp int8
		tmp, err = t.Protocol.ReadByte(ctx)
		*value = byte(tmp)
		return
	case *int8:
		*value, err = t.Protocol.ReadByte(ctx)
		return
	case *int16:
		*value, err = t.Protocol.ReadI16(ctx)
		return
	case *int32:
		*value, err = t.Protocol.ReadI32(ctx)
		return
	case *int64:
		*value, err = t.Protocol.ReadI64(ctx)
		return
	case *float64:
		*value, err = t.Protocol.ReadDouble(ctx)
		return
	case *string:
		*value, err = t.Protocol.ReadString(ctx)
		return
	case *[]byte:
		*value, err = t.Protocol.ReadBinary(ctx)
		return
	case *thrift.TStruct:
		st := &RPCStruct{}
		err = st.Read(ctx, t.Protocol)
		if err != nil {
			return
		}
		*value = st
		return
	case *TypeContainerImplementer:
		var (
			desc TypeContainerDesc
			size int
			typ  TypeContainerImplementer
		)
		switch t.Type {
		case thrift.MAP:
			if desc.Key, desc.Value, size, err = t.Protocol.ReadMapBegin(ctx); err != nil {
				return
			}
			typ, err = NewTypeContainerMapOfTType(desc, t.Required)
			if err != nil {
				return
			}
			typ.SetSize(size)
			if err = typ.Read(ctx, t.Protocol); err != nil {
				return
			}
			if err = t.Protocol.ReadMapEnd(ctx); err != nil {
				return
			}
		case thrift.SET:
			if desc.Value, size, err = t.Protocol.ReadSetBegin(ctx); err != nil {
				return
			}
			typ, err = NewTypeContainerSetOfTType(desc, t.Required)
			if err != nil {
				return
			}
			typ.SetSize(size)
			if err = typ.Read(ctx, t.Protocol); err != nil {
				return
			}
			if err = t.Protocol.ReadSetEnd(ctx); err != nil {
				return
			}
		case thrift.LIST:
			if desc.Value, size, err = t.Protocol.ReadListBegin(ctx); err != nil {
				return
			}
			typ, err = NewTypeContainerListOfTType(desc, t.Required)
			if err != nil {
				return
			}
			typ.SetSize(size)
			if err = typ.Read(ctx, t.Protocol); err != nil {
				return
			}
			if err = t.Protocol.ReadListEnd(ctx); err != nil {
				return
			}
		}
		*value = typ
	default:
		return fmt.Errorf("expected %s, got %T (generic)", t.Type, value)
	}
	return
}
