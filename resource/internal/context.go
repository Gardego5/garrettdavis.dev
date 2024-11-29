package internal

type Key int

const (
	_ Key = iota
	Enforcer
	Fileserver
	Logger
	RequestId
	RequestRef
	RouterMethod
	RouterPath
	Session
	Validate
	WriterRef
)

type GenericKey[T any] struct{}
