package util

type Result[T interface{}] struct {
	Result T
	Error  error
}

type ResultNone struct {
	Result interface{}
}

type ResultErr struct {
	Error error
}
