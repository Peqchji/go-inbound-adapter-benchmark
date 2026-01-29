package pkg

type Result[T any] struct {
	Res T
	Err error
}