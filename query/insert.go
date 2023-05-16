package query

import (
	"fmt"
	"strings"
)

func (q *q[T]) Insert(table, fields string) *q[T] {
	q.op = INSERT
	q.tables = map[string]string{table: ""}
	q.fields = []field{}
	return q
}

func (q *q[T]) getInsertQuery() string {
	query := []string{
		q.op,
		"INTO",
		getStringFromMap(q.tables),
		fmt.Sprintf("(%s)", getSelectString(q.fields)),
		"VALUES",
		getStringFromSlice(q.values),
	}

	if q.returning != "" {
		query = append(query, "RETURNING", q.returning)
	}

	return strings.Join(query, " ")
}

func getStringFromSlice[T any](sl [][]T) string {
	var arr []string

	for _, v := range sl {
		var row []string
		for _, vv := range v {
			row = append(row, fmt.Sprintf("%v", vv))
		}
		arr = append(arr, fmt.Sprintf("(%s)", strings.Join(row, ", ")))
	}
	return strings.Join(arr, ", ")
}

func (q *q[T]) Values(values []any) *q[T] {
	var arr []string

	for _, v := range values {
		namedWhere := q.nextNamedArg()
		arr = append(arr, fmt.Sprintf("@%v", namedWhere))
		q.namedArgs[namedWhere] = v
	}

	q.values = [][]string{arr}
	return q
}

func (q *q[T]) Returning(fields string) *q[T] {
	q.returning = fields
	return q
}
