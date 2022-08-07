package utils

type ErrorKind int

const (
	KindErrHashedPassword ErrorKind = iota
	KindErrUniqueViolation
	KindErrForeignKeyViolation
	KindErrInternalServerError
	KindErrNotFound
	KindErrBadRequest
	KindErrInsufficientBalance
	KindErrUnauthorized
)

type WrapError struct {
	Value error
	Kind  ErrorKind
}
