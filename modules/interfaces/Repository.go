package interfaces

import "web-crawler/modules/types"

type Repository interface {
	GetMany() ([]Document, error)
	GetManyByKeyword(search string, pagination *types.PaginationOptions) ([]Document, error)
	Save(document Document) (Document, error)
}
