package searchtasks

import (
	"strconv"

	"github.com/pkg/errors"
)

const(
	maxRows = 25
)

type Query struct {
	Tags   []string
	Query  string
	Limit  uint
	Offset uint
}

func NewQuery(tags []string, query, limit, offset string) (Query, error) {
	var offsetNum uint64
	if offset != "" {
		var err error
		offsetNum, err = strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return Query{}, errors.Wrap(ErrValidationFailed, err.Error())
		}
	}

	var limitNum uint64 = offsetNum + maxRows
	if limit != "" {
		var err error
		limitNum, err = strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return Query{}, errors.Wrap(ErrValidationFailed, err.Error())
		}
	}

	if limitNum < offsetNum{
	    return Query{}, ErrValidationFailed
	}

	if limitNum - offsetNum > maxRows {
	    limitNum = offsetNum + maxRows
	}

	return Query{
		Tags:   tags,
		Query:  query,
		Limit:  uint(limitNum),
		Offset: uint(offsetNum),
	}, nil
}
