package types

const DEFAULT_ITEMS_PER_PAGE = 10

type PaginationOptions struct {
	Page  int
	Items int
}

func NewPaginationOptions() *PaginationOptions {
	return &PaginationOptions{
		Items: DEFAULT_ITEMS_PER_PAGE,
		Page:  0,
	}
}
