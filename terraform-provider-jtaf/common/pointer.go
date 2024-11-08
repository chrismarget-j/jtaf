package common

func ToPtr[A any](a A) *A {
	return &a
}
