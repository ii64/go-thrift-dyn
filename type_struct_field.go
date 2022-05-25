package thrift_dyn

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

// WriteData write data to wire, without writing struct field.
func (f *TField) WriteData(ctx context.Context, p thrift.TProtocol) (err error) {
	if value, ok := f.Value.(*TypeContainerKV); ok {
		return value.Write(ctx, p)
	}
	switch f.Type {
	case thrift.STOP:
		// nop.
	case thrift.BOOL:
		if value, ok := f.Value.(bool); ok {
			return p.WriteBool(ctx, value)
		} else if f.Required {
			return p.WriteBool(ctx, false)
		}
	case thrift.BYTE:
		switch value := f.Value.(type) {
		case int8:
			return p.WriteByte(ctx, value)
		case byte:
			return p.WriteByte(ctx, int8(value))
		}
		if f.Required {
			return p.WriteByte(ctx, 0)
		}
	case thrift.I16:
		if value, ok := f.Value.(int16); ok {
			return p.WriteI16(ctx, value)
		}
		if f.Required {
			return p.WriteI16(ctx, 0)
		}
	case thrift.I32:
		if value, ok := f.Value.(int32); ok {
			return p.WriteI32(ctx, value)
		}
		if f.Required {
			return p.WriteI32(ctx, 0)
		}
	case thrift.I64:
		switch value := f.Value.(type) {
		case int64:
			return p.WriteI64(ctx, value)
		case int:
			return p.WriteI64(ctx, int64(value))
		}
		if f.Required {
			return p.WriteI64(ctx, 0)
		}
	case thrift.DOUBLE:
		if value, ok := f.Value.(float64); ok {
			return p.WriteDouble(ctx, value)
		}
		if f.Required {
			return p.WriteDouble(ctx, 0)
		}
	case thrift.STRING:
		if value, ok := f.Value.(string); ok {
			return p.WriteString(ctx, value)
		}
		if f.Required {
			return p.WriteString(ctx, "")
		}
	case thrift.STRUCT:
		if value, ok := f.Value.(thrift.TStruct); ok {
			return value.Write(ctx, p)
		}
		if f.Required {
			return nil // nop.
		}
	case thrift.MAP, thrift.SET, thrift.LIST:
		if value, ok := f.Value.(TypeContainerImplementer); ok {
			return value.Write(ctx, p)
		}
		if f.Required {
			return nil // nop.
		}
	}
	return fmt.Errorf("expected %s, got %T", f.Type, f.Value)
}

// Write write struct field with its value.
func (f *TField) Write(ctx context.Context, p thrift.TProtocol) (err error) {
	if f.Value != nil || f.Required {
		if err = p.WriteFieldBegin(ctx, f.Name, f.Type, f.ID); err != nil {
			return
		}
		if err = f.WriteData(ctx, p); err != nil {
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
		typ := NewTypeContainerSet(f.Type, TypeContainerDesc{
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
