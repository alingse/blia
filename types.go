package blia

type StandList[T any] struct {
	Data   []T         `json:"data"`
	Paging *PagingMeta `json:"paging"`
}

type OffsetLimiter interface {
	ToOffsetLimit() OffsetLimit
}

func NewDataList[T any](data []T, off OffsetLimiter) *StandList[T] {
	ol := off.ToOffsetLimit()
	return &StandList[T]{
		Data:   data,
		Paging: &PagingMeta{IsEnd: len(data) < ol.Limit},
	}
}

func NewAllDataList[T any](data []T) *StandList[T] {
	paging := &PagingMeta{IsEnd: true, Totals: len(data)}
	return &StandList[T]{Data: data, Paging: paging}
}

type OffsetLimit struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (o *OffsetLimit) Validate() error {
	if o.Offset < 0 || o.Limit < 0 || o.Limit > 200 {
		return ErrInvalidOffsetLimit
	}
	if o.Limit == 0 {
		o.Limit = 20
	}
	return nil
}

func (o OffsetLimit) ToOffsetLimit() OffsetLimit {
	return o
}

type PageAndPageSize struct {
	OffsetLimit OffsetLimit `json:"-"`
	Page        int         `json:"page"`
	PageSize    int         `json:"page_size"`
}

func (o *PageAndPageSize) Validate() error {
	if o.Page == 0 {
		o.Page = 1
	}
	if o.PageSize == 0 {
		o.PageSize = 20
	}

	if o.Page < 1 || o.PageSize < 0 || o.PageSize > 200 {
		return ErrInvalidOffsetLimit
	}
	o.OffsetLimit = o.ToOffsetLimit()
	return nil
}

func (o PageAndPageSize) ToOffsetLimit() OffsetLimit {
	return OffsetLimit{
		Offset: (o.Page - 1) * o.PageSize,
		Limit:  o.PageSize,
	}
}

type PagingMeta struct {
	Type   string `json:"type"`
	IsEnd  bool   `json:"is_end"`
	Totals int    `json:"totals"`
}

func Empty() map[string]interface{} {
	return map[string]interface{}{}
}

type Validator interface {
	Validate() error
}

type validator struct {
	f func() error
}

func (v validator) Validate() error {
	return v.f()
}

func NewValidator(f func() error) Validator {
	return validator{f: f}
}

func JoinValidator(validators ...Validator) Validator {
	return NewValidator(func() error {
		for _, v := range validators {
			if err := v.Validate(); err != nil {
				return err
			}
		}
		return nil
	})
}
