package documents

import "time"

type ContentDocument struct {
	Id        string    `json:"id"`
	Data      string    `json:"data"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"createdAt"`
}
