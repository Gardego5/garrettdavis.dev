package access

import (
	"context"

	"github.com/Gardego5/garrettdavis.dev/resource/internal"
)

func Get[T any](c context.Context) *T { return c.Value(internal.GenericKey[T]{}).(*T) }
