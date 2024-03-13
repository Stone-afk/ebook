package es

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
)

type Searcher[T any] struct {
	client  *elastic.Client
	idxName []string
	query   elastic.Query
}

func NewSearcher[T any](client *elastic.Client, idxName ...string) *Searcher[T] {
	return &Searcher[T]{
		client:  client,
		idxName: idxName,
	}
}

func (s *Searcher[T]) Query(q elastic.Query) *Searcher[T] {
	s.query = q
	return s
}

//
//func (s *Searcher[T]) Do1(ctx context.Context) (T, error) {
//
//}

func (s *Searcher[T]) Do(ctx context.Context) ([]T, error) {
	resp, err := s.client.Search(s.idxName...).Do(ctx)
	res := make([]T, 0, resp.Hits.TotalHits.Value)
	for _, hit := range resp.Hits.Hits {
		var t T
		err = json.Unmarshal(hit.Source, &t)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

//func (s *Searcher[T]) Resp() *elastic.SearchResult {
//
//}
