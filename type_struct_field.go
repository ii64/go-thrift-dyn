package go_thriftproxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

var (
	ErrSkipField = errors.New("skip field")
)

type TFieldID = int16

type TField struct {
	ID       TFieldID     // i16
	Type     thrift.TType // i8
	Required bool         // i8
	Name     string

	Value any

	// caches
	writeCache func() error // func writer ptr
}

type TFieldImplementerx interface {
	GetID() TFieldID
	GetValue() any

	WriteDataThunk(ctx context.Context, p thrift.TProtocol) (doFunc func() error)
	Write(ctx context.Context, p thrift.TProtocol) (err error)
	Read(ctx context.Context, p thrift.TProtocol) (err error)
}

// NewTField create new TField.
func NewTField(id TFieldID, type_ thrift.TType, name string, required bool) *TField {
	return &TField{ID: id, Name: name, Type: type_, Required: required}
}

// WriteDataThunk write to wire data. Now, pass func ptr that'll be executed later.
func (f *TField) WriteDataThunk() (doFunc func(ctx context.Context, p thrift.TProtocol) error) {
	if _, ok := f.Value.(*TypeContainerKV); ok {
		doFunc = func(ctx context.Context, p thrift.TProtocol) error { return f.Value.(*TypeContainerKV).Write(ctx, p) }
		return
	}
	switch f.Type {
	case thrift.BOOL:
		_, ok := f.Value.(bool)
		if !ok {
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteBool(ctx, false) }
			goto CheckRequireness
		}
		doFunc = func(ctx context.Context, p thrift.TProtocol) error {
			return p.WriteBool(ctx, f.Value.(bool))
		}
	case thrift.BYTE:
		switch f.Value.(type) {
		case int8:
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteByte(ctx, f.Value.(int8)) }
		case byte:
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteByte(ctx, int8(f.Value.(int8))) }
		default:
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteByte(ctx, 0) }
			goto CheckRequireness
		}
	case thrift.I16:
		_, ok := f.Value.(int16)
		if !ok {
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteI16(ctx, 0) }
			goto CheckRequireness
		}
		doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteI16(ctx, f.Value.(int16)) }
	case thrift.I32:
		_, ok := f.Value.(int32)
		if !ok {
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteI32(ctx, 0) }
			goto CheckRequireness
		}
		doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteI32(ctx, f.Value.(int32)) }
	case thrift.I64:
		switch f.Value.(type) {
		case int64:
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteI64(ctx, f.Value.(int64)) }
		case int:
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteI64(ctx, int64(f.Value.(int))) }
		default:
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteI64(ctx, 0) }
			goto CheckRequireness
		}
	case thrift.DOUBLE:
		_, ok := f.Value.(float64)
		if !ok {
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteDouble(ctx, 0) }
			goto CheckRequireness
		}
		doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteDouble(ctx, f.Value.(float64)) }
	case thrift.STRING:
		_, ok := f.Value.(string)
		if !ok {
			doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteString(ctx, "") }
			goto CheckRequireness
		}
		doFunc = func(ctx context.Context, p thrift.TProtocol) error { return p.WriteString(ctx, f.Value.(string)) }

	case thrift.STRUCT:
		_, ok := f.Value.(thrift.TStruct)
		if !ok {
			// don't have zero-value
			goto CheckRequireness
		}
		doFunc = func(ctx context.Context, p thrift.TProtocol) error {
			return f.Value.(thrift.TStruct).Write(ctx, p)
		}
	case thrift.MAP, thrift.SET, thrift.LIST:
		_, ok := f.Value.(TypeContainerImplementer)
		if !ok {
			// don't have zero-value
			goto CheckRequireness
		}
		doFunc = func(ctx context.Context, p thrift.TProtocol) error {
			return f.Value.(TypeContainerImplementer).Write(ctx, p)
		}
	}
	return
CheckRequireness:
	if f.Required && f.Value == nil {
		if doFunc == nil {
			doFunc = func(ctx context.Context, p thrift.TProtocol) error {
				return nil
			}
		}
		return
	}
	return func(ctx context.Context, p thrift.TProtocol) error {
		return fmt.Errorf("expected %s, got %T", f.Type, f.Value)
	}
}

// Write write struct field with its value.
func (f *TField) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	thunk := f.WriteDataThunk() // pending write action.
	if thunk != nil {
		if err = p.WriteFieldBegin(ctx, f.Name, f.Type, f.ID); err != nil {
			return
		}
		if err = thunk(ctx, p); err != nil {
			return
		}
		if err = p.WriteFieldEnd(ctx); err != nil {
			return
		}
	}
	return
}

// Read read struct field with its value.
func (f *TField) Read(ctx context.Context, p thrift.TProtocol) (err error) {
	if value, ok := f.Value.(*TypeContainerKV); ok {
		return value.Read(ctx, p)
	}
	switch f.Type {
	case thrift.BOOL:
		if f.Value, err = p.ReadBool(ctx); err != nil {
			return
		}
	case thrift.BYTE:
		if f.Value, err = p.ReadByte(ctx); err != nil {
			return
		}
	case thrift.I16:
		if f.Value, err = p.ReadI16(ctx); err != nil {
			return
		}
	case thrift.I32:
		if f.Value, err = p.ReadI32(ctx); err != nil {
			return
		}
	case thrift.I64:
		if f.Value, err = p.ReadI64(ctx); err != nil {
			return
		}
	case thrift.DOUBLE:
		if f.Value, err = p.ReadDouble(ctx); err != nil {
			return
		}
	case thrift.STRING:
		if f.Value, err = p.ReadString(ctx); err != nil {
			return
		}
	case thrift.STRUCT:
		st := &RPCStruct{}
		if err = st.Read(ctx, p); err != nil {
			return
		}
		f.Value = st

	case thrift.MAP:
		var (
			keyType, valueType thrift.TType
			size               int
		)
		if keyType, valueType, size, err = p.ReadMapBegin(ctx); err != nil {
			return
		}
		typ := NewTypeContainerMap(f.Type, TypeContainerDesc{
			Key:   keyType,
			Value: valueType,
		}, f.Required)
		typ.SetSize(size)
		if err = typ.Read(ctx, p); err != nil {
			return
		}
		if err = p.ReadMapEnd(ctx); err != nil {
			return
		}
		f.Value = typ
	case thrift.SET:
		var (
			elemType thrift.TType
			size     int
		)
		if elemType, size, err = p.ReadSetBegin(ctx); err != nil {
			return
		}
		typ := NewTypeContainerList(f.Type, TypeContainerDesc{
			Value: elemType,
		}, f.Required)
		typ.SetSize(size)
		if err = typ.Read(ctx, p); err != nil {
			return
		}
		if err = p.ReadSetEnd(ctx); err != nil {
			return
		}
		f.Value = typ
	case thrift.LIST:
		var (
			elemType thrift.TType
			size     int
		)
		if elemType, size, err = p.ReadListBegin(ctx); err != nil {
			return
		}
		typ := NewTypeContainerList(f.Type, TypeContainerDesc{
			Value: elemType,
		}, f.Required)
		typ.SetSize(size)
		if err = typ.Read(ctx, p); err != nil {
			return
		}
		if err = p.ReadListEnd(ctx); err != nil {
			return
		}
		f.Value = typ

	default:
		return ErrSkipField
	}
	return
}

// GetID get field ID.
func (f TField) GetID() TFieldID {
	return f.ID
}

// GetValue get field value.
func (f *TField) GetValue() any {
	return f.Value
}

// SetKey set key of KV container (map)
func (f *TField) SetKey(value any) *TField {
	if kv, ok := f.Value.(*TypeContainerKV); ok {
		if kv.Key != nil {
			kv.Key.Value = value
		}
	}
	return f
}

// SetValue set value of field, or value of KV container (map)
func (f *TField) SetValue(value any) *TField {
	if kv, ok := f.Value.(*TypeContainerKV); ok {
		kv.Value.Value = value
	} else {
		f.Value = value
	}
	return f
}

// SetKeyValue set key and value of KV container (map)
func (f *TField) SetKeyValue(key, value any) *TField {
	if kv, ok := f.Value.(*TypeContainerKV); ok {
		if kv.Key != nil {
			kv.Key.Value = key
		}
		if kv.Value != nil {
			kv.Value.Value = value
		}
	}
	return f
}
