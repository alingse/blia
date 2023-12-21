package blia

import (
	"context"
	"net/http"

	"github.com/alingse/structquery"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Decode(r *http.Request, outPtr Validator) error {
	switch r.Method {
	case http.MethodGet:
		return DecodeQuery(r, outPtr)
	case http.MethodPost, http.MethodPut:
		_, err := DecodeBody(r.Body, outPtr)
		return err
	}
	return nil
}

type GetDB func() *gorm.DB

var StructQueryer = structquery.NewQueryer()

func SimpleQuery[T any, U any, F GetDB](r *http.Request, f F) ([]T, int64, error) {
	var query U
	var queryPtr = &query
	if queryPtrV, ok := (any(queryPtr)).(Validator); ok {
		if err := Decode(r, queryPtrV); err != nil {
			return nil, 0, err
		}
	} else {
		return nil, 0, ErrSimpleQueryInvalid
	}

	conds, err := StructQueryer.And(query)
	if err != nil {
		return nil, 0, err
	}
	var item T

	ctx := r.Context()
	db := f().WithContext(ctx).
		Model(&item).
		Where(conds)

	var result []T
	var totals int64
	if off, ok := (any(query)).(OffsetLimiter); ok {
		db.Count(&totals)
		offlimit := off.ToOffsetLimit()
		db = db.Offset(offlimit.Offset).Limit(offlimit.Limit)
	}
	err = db.Find(&result).Error
	if err != nil {
		return nil, 0, err
	}
	return result, totals, nil
}

type OrderBy interface {
	OrderBy() []clause.OrderByColumn
}

type Clauses interface {
	Clauses() clause.Expression
}

type ModelQuery interface {
	Clauses
	OrderBy
	OffsetLimiter
}

type Entity interface {
	TableName() string
}

type NoOrder struct{}

func (NoOrder) OrderBy() []clause.OrderByColumn {
	return nil
}

func SimpleFetch[T Entity](ctx context.Context, f GetDB, query ModelQuery) ([]T, int64, error) {
	var db = f().WithContext(ctx)
	var items []T
	var total int64
	var err error
	var t T

	db = db.Model(&t)
	if expr := query.Clauses(); expr != nil {
		db = db.Where(expr)
	}
	db.Count(&total)

	if orders := query.OrderBy(); orders != nil {
		for _, order := range orders {
			db = db.Order(order)
		}
	}

	off := query.ToOffsetLimit()
	err = db.Find(&items).
		Order("id desc").
		Offset(off.Offset).
		Limit(off.Limit).
		Error
	if err != nil {
		logrus.Error(query, err)
		return nil, 0, err
	}
	return items, total, nil
}
