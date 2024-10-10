package middleware

type key int

const (
	markdownKey key = iota
	fileserverKey
	loggerKey
	requestIdKey
	writerrefKey
	dbKey
)
