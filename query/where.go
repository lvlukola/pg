package query

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	EQUEL       = uint8(0)
	LIKE        = uint8(1)
	ILIKE       = uint8(2)
	NOT_EQUEL   = uint8(3)
	IN          = uint8(4)
	GREAT       = uint8(5)
	GREAT_OR_EQ = uint8(6)

	LITTLE       = uint8(7)
	LITTLE_OR_EQ = uint8(8)

	IS_NULL     = uint8(9)
	IS_NOT_NULL = uint8(10)
)

var FIND_TYPES = [11]string{"=", "LIKE", "ILIKE", "!=", "IN", ">", ">=", "<", "<=", "is", "is not"}

func (q *q[T]) Where(findType uint8, field, value any) *q[T] {
	q.namedArgsNum = 0
	q.where = ""

	q.AndWhere(findType, field, value)

	return q
}

func (q *q[T]) FilterWhere(findType uint8, field, value any) *q[T] {
	if isFiltered(value) {
		return q
	}

	q.Where(findType, field, value)
	return q
}

func (q *q[T]) AndWhere(findType uint8, field, value any) *q[T] {
	return q.addWhere("AND", findType, field, value)
}

func (q *q[T]) AndFilterWhere(findType uint8, field, value any) *q[T] {
	if isFiltered(value) {
		return q
	}

	q.AndWhere(findType, field, value)
	return q
}

func (q *q[T]) OrWhere(findType uint8, field, value any) *q[T] {
	return q.addWhere("OR", findType, field, value)
}

func (q *q[T]) OrFilterWhere(findType uint8, field, value any) *q[T] {
	if isFiltered(value) {
		return q
	}

	q.OrWhere(findType, field, value)
	return q
}

func (q *q[T]) addWhere(t string, findType uint8, field, value any) *q[T] {
	if findType == LIKE || findType == ILIKE {
		q.addLikeWhere(t, findType, field, value)
	} else if findType == IN {
		q.addInWhere(t, findType, field, value)
	} else if findType == IS_NULL || findType == IS_NOT_NULL {
		q.writeWhere(t, fmt.Sprintf("%s %s %s", field, FIND_TYPES[findType], "null"))
	} else {
		q.addDefaultWhere(t, findType, field, value)
	}

	return q
}

func (q *q[T]) addDefaultWhere(t string, findType uint8, field, value any) *q[T] {
	namedWhere := q.nextNamedArg()
	q.namedArgs[namedWhere] = value
	q.writeWhere(t, fmt.Sprintf("%s %s @%s", field, FIND_TYPES[findType], namedWhere))

	return q
}

func (q *q[T]) addInWhere(t string, findType uint8, field, value any) *q[T] {
	values := reflect.ValueOf(value)
	if reflect.TypeOf(value).Kind() != reflect.Slice || values.Len() == 0 {
		return q
	}
	var namedParams []string

	for i := 0; i < values.Len(); i++ {
		namedWhere := q.nextNamedArg()
		namedParams = append(namedParams, fmt.Sprintf("@%s", namedWhere))
		q.namedArgs[namedWhere] = values.Index(i).Interface()
	}

	q.writeWhere(t, fmt.Sprintf("%s %s (%s)", field, FIND_TYPES[findType], strings.Join(namedParams, ", ")))

	return q
}

func (q *q[T]) addLikeWhere(t string, findType uint8, field, value any) *q[T] {
	namedWhere := q.nextNamedArg()
	q.namedArgs[namedWhere] = fmt.Sprintf("%%%s%%", value)
	q.writeWhere(t, fmt.Sprintf("%s::text %s @%s", field, FIND_TYPES[findType], namedWhere))

	return q
}

func (q *q[T]) writeWhere(t, where string) *q[T] {
	if len(q.where) > 0 {
		q.where += fmt.Sprintf(" %s %s", t, where)
	} else {
		q.where += fmt.Sprintf("WHERE %s", where)
	}

	return q
}

func isFiltered(value any) bool {
	if value == "" || fmt.Sprintf("%v", value) == "<nil>" { //TODO придумать что-то без nil
		return true
	}

	return false
}
