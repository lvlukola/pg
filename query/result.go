package query

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (q *q[T]) All() ([]T, error) {
	rows, _ := q.client.Query(q.ctx, q.GetQuery(), q.GetArgs())

	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[T])

	if err != nil && err != pgx.ErrNoRows {
		return []T{}, err
	}

	return data, nil
}

func (q *q[T]) One() (T, error) {
	q.Limit(1)

	rows, _ := q.client.Query(q.ctx, q.GetQuery(), q.GetArgs())

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[T])

	if err == nil /*|| err == pgx.ErrNoRows*/ {
		return data, nil
	} else {
		return data, err
	}
}

func (q *q[T]) Exists() (bool, error) {
	query := strings.Join([]string{
		q.op,
		"COUNT(*)",
		"FROM",
		getStringFromMap(q.tables),
		q.where,
	}, " ")

	var count int64

	if err := q.client.QueryRow(q.ctx, query, q.GetArgs()).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (q *q[T]) Count() (int64, error) {
	query := strings.Join([]string{
		q.op,
		"COUNT(*)",
		"FROM",
		getStringFromMap(q.tables),
		q.where,
		q.limit,
		q.offset,
	}, " ")

	var count int64

	if err := q.client.QueryRow(q.ctx, query, q.GetArgs()).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (q *q[T]) CountFiltered() (int64, error) {
	query := strings.Join([]string{
		q.op,
		"COUNT(*)",
		"FROM",
		getStringFromMap(q.tables),
		q.where,
	}, " ")

	var count int64

	if err := q.client.QueryRow(q.ctx, query, q.GetArgs()).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (q *q[T]) CountTotal() (int64, error) {
	query := strings.Join([]string{
		q.op,
		"COUNT(*)",
		"FROM",
		getStringFromMap(q.tables),
	}, " ")

	var count int64

	if err := q.client.QueryRow(q.ctx, query, q.GetArgs()).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (q *q[T]) GetArgs() pgx.NamedArgs {
	return q.namedArgs
}

func (q *q[T]) GetQuery() string {
	switch q.op {
	case SELECT:
		return q.getSelectQuery()
	case INSERT:
		return q.getInsertQuery() //нужно тестить
	}

	return fmt.Sprintf("Unknown query type: %s", q.op)
}

func (q *q[T]) getSelectQuery() string {
	query := []string{
		q.op,
		getSelectString(q.fields),
		"FROM",
		getStringFromMap(q.tables),
		q.getJoins(),
		q.where,
		q.orderBy,
		q.limit,
		q.offset,
	}

	return strings.Join(query, " ")
}

func getSelectString(elements []field) string {
	var arr []string

	for _, v := range elements {
		element := v.sel

		if v.alias != "" {
			element = fmt.Sprintf("%s as %s", v.sel, v.alias)
		}

		arr = append(arr, strings.Trim(element, " "))
	}

	return strings.Join(arr, ", ")
}

func getStringFromMap(elements map[string]string) string {
	var arr []string

	for k, v := range elements {
		element := fmt.Sprintf("%s %s", k, v)
		arr = append(arr, strings.Trim(element, " "))
	}

	return strings.Join(arr, ", ")
}
