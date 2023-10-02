package tql

import (
	"fmt"
	"strings"
)

func ToSortSQL(items []SortItem) string {
	if len(items) == 0 {
		return ""
	}

	sf := strings.Builder{}
	sf.WriteString(" ORDER BY")
	for _, f := range items {
		if f.Kind == "asc" {
			sf.WriteString(fmt.Sprintf(" %s ASC,", f.Field))
		} else {
			sf.WriteString(fmt.Sprintf(" %s DESC,", f.Field))
		}
	}

	return sf.String()[:sf.Len()-1]
}

func CursorToSQL(q *Query, filter string) (sqlWhere string, sqlLimit int64, sqlOrder string) {

	sqlLimit = q.PaginationRAWItems()
	var op = ">"
	var sort = "ASC"

	var cursor *Cursor = CursorFromQuery(q)
	if cursor.Order() == CursorDesc {
		op = "<"
		sort = "DESC"
	}

	if q.HasSort() {
		sortItems := append([]SortItem{}, *q.CursorPaginationSort(), SortItem{Field: "id", Kind: sort})
		sqlOrder = ToSortSQL(sortItems)
		sqlWhere = toSqlCursorFilter(filter, cursor)
	} else {
		sqlOrder = fmt.Sprintf(" ORDER BY id %s", sort)
		sqlWhere = fmt.Sprintf(" WHERE (%s) AND (id %s '%s')", filter, op, cursor.Pivot())
	}

	return sqlWhere, sqlLimit, sqlOrder
}

func toSqlCursorFilter(filter string, cursor *Cursor) (q string) {

	order := cursor.Order()
	var op = ">"

	if order == CursorDesc {
		op = "<"
	}

	var sortPart string
	var sortEqual string
	var sortOp = ">"
	if cursor.SortType() != "asc" {
		sortOp = "<"
	}

	if cursor.SortVal() != nil {
		switch v := cursor.SortVal().(type) {
		case string:
			sortPart = fmt.Sprintf("%s %s '%s'", cursor.SortField(), sortOp, v)
			sortEqual = fmt.Sprintf("%s = '%s'", cursor.SortField(), v)
		default:
			sortPart = fmt.Sprintf("%s %s %v", cursor.SortField(), sortOp, v)
			sortEqual = fmt.Sprintf("%s = %s", cursor.SortField(), v)
		}
		q = fmt.Sprintf(" WHERE ((%s) AND (%s)) OR ((%s) AND (%s) AND (id %s '%s'))", sortPart, filter, filter, sortEqual, op, cursor.Pivot())
	} else {
		if cursor.Pivot() != "" {
			q = fmt.Sprintf(" WHERE ((%s) AND (id %s '%s'))", filter, op, cursor.Pivot())
		} else {
			q = fmt.Sprintf(" WHERE %s", filter)
		}

	}

	return q
}

/*
SQL:
-----

SELECT *
FROM `users`
WHERE
	((name > "6") AND (created >= 1666875856369))
    OR
    ((created >=1666875856369) AND name="6" AND id > "01GGCRPQ06W18NGGYT9R11V02E" )

 ORDER BY name DESC, id DESC
 LIMIT 3;


http://go-database-sql.org/retrieving.html
In MySQL, the parameter placeholder is ?, and in PostgreSQL it is $N, where N is a number. SQLite accepts either of these. In Oracle placeholders begin with a colon and are named, like :param1. We’ll use ? because we’re using MySQL as our example.

http://go-database-sql.org/modifying.html
_, err := db.Exec("DELETE FROM users")  // OK
_, err := db.Query("DELETE FROM users") // BAD
They do not do the same thing, and you should never use Query() like this. The Query() will return a sql.Rows, which reserves a database connection until the sql.Rows is closed.

*/
