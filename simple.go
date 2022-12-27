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

var queryer = structquery.NewQueryer()

func SimpleQuery[T any, U Validator, F GetDB](r *http.Request, f F) ([]T, error) {
	var query U
	if err := Decode(r, query); err != nil {
		return nil, err
	}
	conds, err := queryer.And(query)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	db := f().WithContext(ctx)
	var item T
	var result []T
	err = db.Model(&item).Where(conds).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
