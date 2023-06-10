package mocks

import (
	"testing"
	"web-crawler/modules/documents"
	"web-crawler/modules/elastic/repositories"
	"web-crawler/modules/logger"
	generalTypes "web-crawler/modules/types"
)

type MockContentRepository struct {
	repositories.ContentRepository
	*testing.T
}

func NewMockContentRepository(t *testing.T) *MockContentRepository {
	return &MockContentRepository{
		*repositories.NewContentRepository(
			&logger.NewMockLoggerService(t).LoggerService,
		),
		t,
	}
}

func (repository *MockContentRepository) GetMany() ([]documents.ContentDocument, error) {
	// Do nothing

	return []documents.ContentDocument{}, nil
}

func (repository *MockContentRepository) GetManyByKeyword(
	search string,
	pagination *generalTypes.PaginationOptions,
) ([]documents.ContentDocument, error) {
	// Do nothing

	return []documents.ContentDocument{}, nil
}

func (repository *MockContentRepository) Save(document documents.ContentDocument) {
	// Do nothing
}
