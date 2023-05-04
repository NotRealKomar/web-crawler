package interfaces

import "web-crawler/modules/types"

type Service interface {
	SearchByKeyword(search string, pagination *types.PaginationOptions) (*types.PaginatedResponse[Model], error)
}
