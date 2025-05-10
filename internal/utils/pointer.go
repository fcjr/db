package utils

// ToPtr returns a pointer to the value provided.
// This is useful for setting optional parameters.
func ToPtr[T any](v T) *T {
	return &v
}
