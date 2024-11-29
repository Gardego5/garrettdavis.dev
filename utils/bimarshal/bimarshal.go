package bimarshal

import (
	"encoding"
	"encoding/json"

	"github.com/tinylib/msgp/msgp"
)

type (
	msgpImpl interface {
		msgp.Unmarshaler
		msgp.Marshaler
	}

	Bimarshal interface {
		encoding.BinaryUnmarshaler
		encoding.BinaryMarshaler
	}
)

type msgpBimarshal[T msgpImpl] struct{ d T }

func (a *msgpBimarshal[T]) MarshalBinary() ([]byte, error) { return a.d.MarshalMsg(nil) }
func (a *msgpBimarshal[T]) UnmarshalBinary(b []byte) error { _, err := a.d.UnmarshalMsg(b); return err }

var _ Bimarshal = (*msgpBimarshal[msgpImpl])(nil)

func MessagePack[T msgpImpl](data T) Bimarshal { return &msgpBimarshal[T]{d: data} }

type jsonBimarshal[T any] struct{ d T }

func (a *jsonBimarshal[T]) MarshalBinary() ([]byte, error) { return json.Marshal(a.d) }
func (a *jsonBimarshal[T]) UnmarshalBinary(b []byte) error { return json.Unmarshal(b, &a.d) }

var _ Bimarshal = (*jsonBimarshal[any])(nil)

func JSON[T any](data T) Bimarshal { return &jsonBimarshal[T]{d: data} }
