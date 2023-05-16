package query

type List[T any] struct {
	RecordsTotal    int64 `json:"records_total"`
	RecordsFiltered int64 `json:"records_filtered"`
	Data            []T   `json:"data"`
}

func (q *q[T]) GetList() (List[T], error) {
	var resp List[T]
	var err error

	if resp.RecordsTotal, err = q.CountTotal(); err != nil {
		return List[T]{}, err
	}

	if resp.RecordsFiltered, err = q.CountFiltered(); err != nil {
		return List[T]{}, err
	}

	if resp.Data, err = q.All(); err != nil {
		return List[T]{}, err
	}

	return resp, nil
}
