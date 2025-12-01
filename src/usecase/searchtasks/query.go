package searchtasks


const(
	maxRows = 25
)

type Query struct {
	Tags   []string
	Query  string
	Limit  uint
	Offset uint
}

func NewQuery(tags []string, query string, limit, offset uint) (Query, error) {
	if limit == 0 || limit > maxRows {
		limit = maxRows
	}

	return Query{
		Tags:   tags,
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}, nil
}
