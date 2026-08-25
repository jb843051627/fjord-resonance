package api

import "github.com/jb843051627/fjord-resonance/internal/store"

type PageResponse[T any] struct {
	Items   []T  `json:"items"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

func NewPageResponse[T any](items []T, page store.Page) PageResponse[T] {
	return PageResponse[T]{Items: items, Limit: page.Limit, Offset: page.Offset, HasMore: len(items) >= page.Limit}
}
