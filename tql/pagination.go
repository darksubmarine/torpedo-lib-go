package tql

import (
	"encoding/base64"
	"encoding/json"
	"github.com/darksubmarine/torpedo-lib-go/ptr"
	"strings"
)

type CursorOrderType int8

const (
	CursorNext = "next"
	CursorPrev = "prev"

	CursorAscStr  = "asc"
	CursorDescStr = "desc"

	cursorPrevPrefix = "prev::"
	cursorNextPrefix = "next::"

	_ CursorOrderType = iota
	CursorDesc
	CursorAsc
)

type ExportableCursorSort struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type ExportableCursor struct {
	Pivot     string                `json:"pivot"`
	PivotSort string                `json:"pivotSort"`
	Sort      *ExportableCursorSort `json:"sort,omitempty"`
}

func NewExportableCursorFrom(cursor string) (*ExportableCursor, error) {
	var data ExportableCursor
	var dst = make([]byte, 0)
	dst, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dst, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (ec *ExportableCursor) String() string {
	data, _ := json.Marshal(ec)
	return base64.StdEncoding.EncodeToString(data)
}

type Cursor struct {
	pivot *string
	order CursorOrderType

	// additional sorting field
	sortField string
	sortVal   interface{}
	sortType  string
}

func PaginationCursorOrder(sort string) CursorOrderType {
	switch sort {
	case CursorAscStr:
		return CursorAsc
	case CursorDescStr:
		return CursorDesc
	default:
		return CursorAsc
	}
}

func NewPaginationCursorFrom(pivot string, sort CursorOrderType, sortField string, sortVal interface{}, sortType string) *Cursor {
	c := &Cursor{}

	switch sort {
	case CursorAsc:
		c.order = CursorAsc
	case CursorDesc:
		c.order = CursorDesc
	default:
		c.order = CursorAsc
	}

	c.pivot = ptr.String(pivot)
	c.sortType = sortType
	c.sortField = sortField
	c.sortVal = sortVal
	return c
}

func NewPaginationCursorFromPivot(pivot string, sort CursorOrderType) *Cursor {
	return NewPaginationCursorFrom(pivot, sort, "", nil, "") // TODO see sort CursorAsc
}

func NewPaginationCursorFromBeginning(sort string) *Cursor {
	var s CursorOrderType
	switch sort {
	case CursorAscStr:
		s = CursorAsc
	case CursorDescStr:
		s = CursorDesc
	default:
		s = CursorAsc
	}

	return NewPaginationCursorFrom("", s, "", nil, "")
}

func NewPaginationCursor(rawCursor string) (*Cursor, error) {
	c := &Cursor{}
	if rawCursor != "" {
		data, err := c.parseCursor(rawCursor)
		if err != nil {
			return nil, err
		}

		c.pivot = ptr.String(data.Pivot)
		c.order = c.extractOrder(data.PivotSort)

		if data.Sort != nil {
			c.sortField = data.Sort.Field
			c.sortType = data.Sort.Type
			c.sortVal = data.Sort.Value
		}
	}

	return c, nil
}

func (c *Cursor) parseCursor(rawCursor string) (*ExportableCursor, error) {
	return NewExportableCursorFrom(rawCursor)
}

func (c *Cursor) extractPivot(rawCursor string) string {
	return strings.TrimLeft(strings.TrimLeft(rawCursor, "prev::"), "next::")
}

func (c *Cursor) extractOrder(pivotSort string) CursorOrderType {
	switch pivotSort {
	case CursorAscStr:
		return CursorAsc
	case CursorDescStr:
		return CursorDesc
	default:
		return CursorAsc
	}
}

func (c *Cursor) cursor(direction string) *string {
	if c.pivot == nil {
		return nil
	}

	exp := ExportableCursor{
		Pivot:     *c.pivot,
		PivotSort: c.ToOrder(),
		Sort: &ExportableCursorSort{
			Field: c.sortField,
			Value: c.sortVal,
			Type:  c.sortType,
		}}

	// TODO base64 encode
	//return ptr.String(fmt.Sprintf("%s::%s", direction, *c.pivot))
	return ptr.String(exp.String())
}

func (c *Cursor) Next() *string {
	return c.cursor(CursorNext)
}

func (c *Cursor) Prev() *string {
	return c.cursor(CursorPrev)
}

func (c *Cursor) Order() CursorOrderType {
	return c.order
}

func (c *Cursor) ToOrder() string {
	switch c.order {
	case CursorAsc:
		return CursorAscStr
	case CursorDesc:
		return CursorDescStr
	default:
		return CursorAscStr
	}
}

func (c *Cursor) Pivot() string {
	if c.pivot == nil {
		return PaginationCursorEmpty
	}

	return *c.pivot
}

func (c *Cursor) SortField() string {
	return c.sortField
}

func (c *Cursor) SortType() string {
	return c.sortType
}

func (c *Cursor) SortVal() interface{} {
	return c.sortVal
}

func CursorFromQuery(q *Query) *Cursor {
	var cursor *Cursor
	if q.PaginationCursor() != "" {
		cursor, _ = NewPaginationCursor(q.PaginationCursor())
	} else {
		if q.HasSort() {
			paginationSort := q.CursorPaginationSort()
			pivotSort := CursorDesc
			if q.CursorType() == PaginationCursortOrderASC {
				pivotSort = CursorAsc
			}

			cursor = NewPaginationCursorFrom(
				"",
				pivotSort,
				paginationSort.Field,
				nil,
				paginationSort.Kind)
		} else {
			cursor = NewPaginationCursorFromBeginning(q.PaginationCursorOrder())
		}
	}
	return cursor
}
