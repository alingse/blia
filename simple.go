package blia

import (
	"net/http"

	"github.com/alingse/structquery"
	"gorm.io/gorm"
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
