package query

import (
	"fmt"
	"strings"
)

type join struct {
	t     string
	alias string
	table string
	on    string
}

func (q *q[T]) LeftJoin(alias, table, on string) *q[T] {
	return q.addJoin("LEFT JOIN", alias, table, on)
}

func (q *q[T]) RightJoin(alias, table, on string) *q[T] {
	return q.addJoin("RIGHT JOIN", alias, table, on)
}

func (q *q[T]) InnerJoin(alias, table, on string) *q[T] {
	return q.addJoin("INNER JOIN", alias, table, on)
}

func (q *q[T]) addJoin(t, alias, table, on string) *q[T] {
	q.joins = append(q.joins, join{
		t,
		alias,
		table,
		on,
	})

	return q
}

func (q *q[T]) getJoins() string {
	var arr []string

	for _, v := range q.joins {
		element := fmt.Sprintf("%s %s %s ON %s", v.t, v.table, v.alias, v.on)
		arr = append(arr, strings.Trim(element, " "))
	}

	return strings.Join(arr, " ")
}
