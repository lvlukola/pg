package query

import "fmt"

type LogSql struct {
	Query  string `json:"query,omitempty"`
	Params any    `json:"params,omitempty"`
	Result any    `json:"result,omitempty"`
	Error  any    `json:"error,omitempty"`
}

func NewLogSql(query string, params ...any) LogSql {
	var paramsArray []any

	for _, p := range params {
		paramsArray = append(paramsArray, p)
	}

	return LogSql{query, paramsArray, nil, nil}
}

func (q *q[T]) GetLogSql() LogSql {
	return LogSql{q.GetQuery(), q.GetArgs(), nil, nil}
}

func (l LogSql) GetMsg() string {
	return fmt.Sprintf("%+v", l)
}

func (l LogSql) SetResult(result any) LogSql {
	l.Result = fmt.Sprintf("%+v", result)
	return l
}

func (l LogSql) SetError(error any) LogSql {
	l.Error = error
	return l
}
