package types

import "web-crawler/modules/interfaces"

type searchResponseTotal struct {
	Value int `json:"value"`
}

type searchResponseDocument[T interfaces.Document] struct {
	Source T      `json:"_source"`
	Index  string `json:"_index"`
	Id     string `json:"_id"`
}

type searchResponseHits[T interfaces.Document] struct {
	Total searchResponseTotal         `json:"total"`
	Hits  []searchResponseDocument[T] `json:"hits"`
}

type SearchResponse[T interfaces.Document] struct {
	Took     int                   `json:"took"`
	Hits     searchResponseHits[T] `json:"hits"`
	TimedOut bool                  `json:"timed_out"`
}
