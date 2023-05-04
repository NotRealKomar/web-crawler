package types

type PaginatedResponse[T any] struct {
	Models []T
	Page   int
}
