package query

import (
	"film-management/pkg/query/sort"
)

type FilterSortLimit struct {
	Sort   sort.Sortable
	Filter Filter
	Limit  int
	Offset int
}

type Filter map[string]interface{}

type FilterSortLimitBuilder interface {
	SetSort(sort.Sortable) FilterSortLimitBuilder
	SetFilter(Filter) FilterSortLimitBuilder
	SetLimit(int) FilterSortLimitBuilder
	SetOffset(int) FilterSortLimitBuilder
	Build() FilterSortLimit
}

type filterSortLimitBuilder struct {
	sort   sort.Sortable
	filter Filter
	limit  int
	offset int
}

func NewFilterSortLimitBuilder() FilterSortLimitBuilder {
	return &filterSortLimitBuilder{}
}

func (b *filterSortLimitBuilder) SetSort(sort sort.Sortable) FilterSortLimitBuilder {
	b.sort = sort
	return b
}

func (b *filterSortLimitBuilder) SetFilter(filter Filter) FilterSortLimitBuilder {
	b.filter = filter
	return b
}

func (b *filterSortLimitBuilder) SetLimit(limit int) FilterSortLimitBuilder {
	b.limit = limit
	return b
}

func (b *filterSortLimitBuilder) SetOffset(offset int) FilterSortLimitBuilder {
	b.offset = offset
	return b
}

func (b *filterSortLimitBuilder) Build() FilterSortLimit {
	return FilterSortLimit{
		Sort:   b.sort,
		Filter: b.filter,
		Limit:  b.limit,
		Offset: b.offset,
	}
}
