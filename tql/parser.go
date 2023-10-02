// Package tql Torpedo Query Language
package tql

import (
	"fmt"
	"github.com/darksubmarine/torpedo-lib-go/entity"
	"github.com/darksubmarine/torpedo-lib-go/ptr"
)

const (
	PageTypeOffset = "offset"
	PageTypeCursor = "cursor"

	PaginationCursortOrderASC = "asc"
	PaginationCursorOrderDESC = "desc"

	PaginationCursorEmpty = ""

	paginationOffsetDefaultItems = 10
	paginationOffsetDefaultPage  = 1
	paginationCursorDefaultOrder = "asc"
)

type Filter struct {
	Type_  *string      `json:"type"`
	Fields []FilterItem `json:"fields"`
}

func (f *Filter) Type() string {
	if f.Type_ != nil {
		return *f.Type_
	}

	return "all"
}

type FilterItem struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type SortItem struct {
	Field string `json:"field"`
	Kind  string `json:"type"`
}

type PaginationMeta struct {
	// common
	Items *int64 `json:"items,omitempty"`

	// offset pagination
	Page *int64 `json:"page,omitempty"`

	// cursor pagination
	Cursor *string `json:"cursor,omitempty"`
	//Order  *string `json:"order,omitempty"`
}

type PaginationCursorMeta struct {
	Sort *SortItem `json:"sort,omitempty"`
	Kind *string   `json:"type,omitempty"`
}

type PaginationCursor struct {
	Meta *PaginationCursorMeta `json:"meta,omitempty"`
	Mark *string               `json:"nextToken,omitempty"`
}

type PaginationOffset struct {
	Page *int64     `json:"page,omitempty"`
	Sort []SortItem `json:"sort,omitempty"`
}

type Pagination struct {
	//Type_ string         `json:"type"`
	//Meta  PaginationMeta `json:"meta"`

	Items  *int64            `json:"items,omitempty"`
	Cursor *PaginationCursor `json:"cursor,omitempty"`
	Offset *PaginationOffset `json:"offset,omitempty"`
}

type Query struct {
	Filter     *Filter     `json:"filter"`
	Projection []string    `json:"projection,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func (q *Query) Validate(fields entity.FieldMap) error {

	// projection
	if len(q.Projection) > 0 {
		for _, f := range q.Projection {
			if !fields.HasField(f) {
				return fmt.Errorf("%w: %s", ErrInvalidFieldNameAtProjection, f)
			}
		}
	}

	// Pagination
	if q.HasPagination() && q.IsOffsetPagination() && q.IsCursorPagination() || q.HasPagination() && !q.HasValidPagination() {
		return ErrInvalidPaginationType
	}

	// Sort
	if q.HasPagination() && q.HasSort() {
		if q.IsCursorPagination() {
			if q.Pagination.Cursor.Meta.Sort != nil {
				if !fields.HasField(q.Pagination.Cursor.Meta.Sort.Field) {
					return fmt.Errorf("%w: %s", ErrInvalidFieldNameAtSort, q.Pagination.Cursor.Meta.Sort.Field)
				}
			}
		} else if q.IsOffsetPagination() {
			if len(q.Pagination.Offset.Sort) > 0 {
				for _, f := range q.Pagination.Offset.Sort {
					if !fields.HasField(f.Field) {
						return fmt.Errorf("%w: %s", ErrInvalidFieldNameAtSort, f.Field)
					}
				}
			}
		}
	}

	return nil
}

func (q *Query) HasSort() bool {
	if q.HasPagination() && q.IsCursorPagination() && q.HasCursorMeta() {
		return q.Pagination.Cursor.Meta.Sort != nil
	}

	if q.HasPagination() && q.IsOffsetPagination() {
		return len(q.Pagination.Offset.Sort) > 0
	}

	return false
}

func (q *Query) OffsetPaginationSort() []SortItem {
	return q.Pagination.Offset.Sort
}

func (q *Query) CursorPaginationSort() *SortItem {
	return q.Pagination.Cursor.Meta.Sort
}

func (q *Query) HasPagination() bool {
	if q.Pagination != nil {
		return true
	}

	return false
}

func (q *Query) HasValidPagination() bool {
	return q.IsOffsetPagination() || q.IsCursorPagination()
}

func (q *Query) IsOffsetPagination() bool {
	if q.Pagination != nil && q.Pagination.Offset != nil {
		return true
	}

	return false
}

func (q *Query) IsCursorPagination() bool {
	if q.Pagination != nil && q.Pagination.Cursor != nil {
		return true
	}

	return false
}

func (q *Query) HasCursorMeta() bool {
	if q.Pagination != nil && q.Pagination.Cursor != nil && q.Pagination.Cursor.Meta != nil {
		return true
	}

	return false
}

func (q *Query) CursorType() string {
	if q.HasCursorMeta() {
		return ptr.ToString(q.Pagination.Cursor.Meta.Kind)
	}

	return ""
}

func (q *Query) IsCursorPaginationPrev() bool {
	cur, _ := NewPaginationCursor(q.PaginationCursor())
	if cur.Order() == CursorDesc {
		return true
	}
	return false
}

func (q *Query) IsCursorPaginationNext() bool {
	return !q.IsCursorPaginationPrev()
}

func (q *Query) PaginationOffset() int64 {
	var items int64 = paginationOffsetDefaultItems // default items
	var page int64 = paginationOffsetDefaultPage   //default page

	if q.Pagination.Items != nil {
		items = *q.Pagination.Items
	}

	if q.Pagination.Offset.Page != nil && *q.Pagination.Offset.Page > 0 {
		page = *q.Pagination.Offset.Page
	}

	return items * (page - 1)
}

func (q *Query) PaginationRAWItems() int64 {
	if q.Pagination.Items != nil {
		return *q.Pagination.Items
	}

	return paginationOffsetDefaultItems
}

func (q *Query) PaginationItems() int64 {
	return q.PaginationRAWItems() + 1
}

func (q *Query) PaginationOffsetPage() int64 {
	if q.Pagination.Offset.Page != nil {
		return *q.Pagination.Offset.Page
	}

	return paginationOffsetDefaultPage
}

func (q *Query) PaginationCursor() string {
	if q.Pagination.Cursor.Mark == nil {
		return PaginationCursorEmpty
	}

	return *q.Pagination.Cursor.Mark
}

func (q *Query) IsPaginationCursorEmpty() bool {
	return q.PaginationCursor() == PaginationCursorEmpty
}

func (q *Query) PaginationCursorOrder() string {

	if q.IsCursorPagination() && q.HasCursorMeta() {
		return ptr.ToString(q.Pagination.Cursor.Meta.Kind)
	}

	return paginationCursorDefaultOrder
}
