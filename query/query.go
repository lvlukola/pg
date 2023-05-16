package query

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	SELECT = "SELECT"
	INSERT = "INSERT"
)

type field struct {
	sel   string
	alias string
}

type DbClient interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type q[T any] struct {
	ctx          context.Context
	client       DbClient
	op           string
	tables       map[string]string
	joins        []join
	fields       []field
	namedArgs    pgx.NamedArgs
	where        string
	orderBy      string
	limit        string
	offset       string
	namedArgsNum int8
	values       [][]string
	returning    string
}

func New[T any](ctx context.Context, client DbClient) *q[T] {
	return &q[T]{
		ctx,
		client,
		"",
		map[string]string{},
		[]join{},
		[]field{},
		pgx.NamedArgs{},
		"",
		"",
		"",
		"",
		0,
		[][]string{},
		"",
	}
}

func (q *q[T]) Select(sel, alias string) *q[T] {
	q.op = SELECT
	selForAdd := field{sel, alias}
	q.fields = []field{selForAdd}
	return q
}

func (q *q[T]) AddSelect(sel, alias string) *q[T] {
	q.op = SELECT
	selForAdd := field{sel, alias}
	q.fields = append(q.fields, selForAdd)
	return q
}

func (q *q[T]) From(table, alias string) *q[T] {
	q.tables = map[string]string{table: alias}
	return q
}

func (q *q[T]) nextNamedArg() string {
	q.namedArgsNum += 1
	return fmt.Sprintf("a%d", q.namedArgsNum)
}

func (q *q[T]) OrderBy(field string) *q[T] {
	q.orderBy = fmt.Sprintf("ORDER BY %s", field)
	return q
}

func (q *q[T]) AndOrderBy(field string) *q[T] {
	q.orderBy += fmt.Sprintf(",  %s", field)
	return q
}

func (q *q[T]) Offset(offset uint16) *q[T] {
	if offset > 0 {
		q.offset = "OFFSET @offset"
		q.namedArgs["offset"] = offset
	}

	return q
}

func (q *q[T]) Limit(limit int) *q[T] {
	if limit > 0 {
		q.limit = "LIMIT @limit"
		q.namedArgs["limit"] = limit
	}

	return q
}
