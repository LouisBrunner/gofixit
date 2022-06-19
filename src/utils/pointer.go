package utils

func Pointerize[T any](v T) *T {
	return &v
}
