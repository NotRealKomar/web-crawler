package documents

import (
	"time"

	"github.com/google/uuid"
)

type ContentDocument struct {
	Id        string    `json:"id"`
	Data      string    `json:"data"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewContentDocument() *ContentDocument {
	return &ContentDocument{
		Id:        uuid.NewString(),
		CreatedAt: time.Now(),
	}
}
