package utils

func Ptr[T any](v T) *T {
	return &v
}

func Deref[T any](v *T) T {
	var zero T
	if v == nil {
		return zero
	}
	return *v
}

func IsNil[T any](ptr *T) bool {
	return ptr == nil
}
