package tql

import "go.mongodb.org/mongo-driver/bson"

func CursorToDocument(q *Query, filter bson.D) (cursorFilter bson.D, limit int64, sortD bson.D) {

	var cursor *Cursor = CursorFromQuery(q)

	order := cursor.Order()
	items := q.PaginationRAWItems()

	var op = "$gt"
	var sort = 1

	if order == CursorDesc {
		op = "$lt"
		sort = -1
	}

	if q.HasSort() {
		sortD = append(ToSortDocument([]SortItem{*q.CursorPaginationSort()}), bson.E{Key: "_id", Value: sort})
		cursorFilter = toCursorFilter(filter, cursor)
	} else {
		sortD = bson.D{{Key: "_id", Value: sort}}
		cursorFilter = bson.D{{"$and", []interface{}{
			bson.D{{Key: "_id", Value: bson.D{bson.E{Key: op, Value: cursor.Pivot()}}}},
			filter,
		}}}
	}

	if cursor.Pivot() == PaginationCursorEmpty {
		return filter, items, sortD
	}

	return cursorFilter, items, sortD
}

func ToSortDocument(items []SortItem) bson.D {
	sf := bson.D{}
	for _, f := range items {
		var sortType int8
		if f.Kind == "asc" {
			sortType = 1
		} else {
			sortType = -1
		}
		sf = append(sf, bson.E{Key: f.Field, Value: sortType})
	}

	return sf
}

func toCursorFilter(filter bson.D, cursor *Cursor) bson.D {

	order := cursor.Order()
	var op = "$gt"

	if order == CursorDesc {
		op = "$lt"
	}

	var f bson.D
	if cursor.SortVal() != nil {
		sortDocument := bson.D{}
		if cursor.SortType() == "asc" {
			sortDocument = bson.D{{cursor.SortField(), bson.D{{"$gt", cursor.SortVal()}}}}
		} else {
			sortDocument = bson.D{{cursor.SortField(), bson.D{{"$lt", cursor.SortVal()}}}}
		}

		f = bson.D{
			{"$or", []interface{}{
				bson.D{{"$and", []interface{}{
					sortDocument,
					filter,
				}}},

				bson.D{{"$and", []interface{}{
					filter,
					bson.D{{cursor.SortField(), bson.D{{"$eq", cursor.SortVal()}}}},
					bson.D{{Key: "_id", Value: bson.D{bson.E{Key: op, Value: cursor.Pivot()}}}},
				}}},
			}},
		}
	} else {
		if cursor.Pivot() != "" {
			f = bson.D{{"$and", []interface{}{
				filter,
				bson.D{{Key: "_id", Value: bson.D{bson.E{Key: op, Value: cursor.Pivot()}}}},
			}}}
		} else {
			f = filter
		}
	}

	return f
}

/*

IMPORTANT::: Lesson learned: NEVER sort by fields that may contain null values

https://brunoscheufler.com/blog/2022-01-01-paginating-large-ordered-datasets-with-cursor-based-pagination
https://shopify.engineering/pagination-relative-cursors
// cursor= next::01GGCRPWY5Q2BXMBQZY2QRS1SR::sortField=name::sortValue=7::sortType=desc
db.getCollection("users").find(
    {

       $or:[
       {$and:[{name: {"$lt":"7"}},{created: {"$gte":1666875856369}}]},

       {$and:[
           {created: {"$gte":1666875856369}},

           {name: {"$eq":"7"}},
           {_id: {"$gt": "01GGCRPWY5Q2BXMBQZY2QRS1SR"}}

           ]}
       ]


    }
).limit(30)
.sort({name:-1},{"_id":-1})


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
