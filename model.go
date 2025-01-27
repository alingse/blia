package blia

import (
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `json:"id"         gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at" gorm:"<-:create"`
	UpdatedAt time.Time      `json:"updated_at" gorm:""`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (m Model) GetModelID() uint {
	return m.ID
}

func (m Model) GetModel() Model {
	return m
}

type WithModel interface {
	GetModelID() uint
}

func GetModelIDs[T WithModel](ms []T) []uint {
	ids := make([]uint, 0, len(ms))
	for _, m := range ms {
		ids = append(ids, m.GetModelID())
	}
	return ids
}

func GetModelIDSet[T WithModel](ms []T) mapset.Set[uint] {
	ids := mapset.NewSet[uint]()
	for _, m := range ms {
		ids.Add(m.GetModelID())
	}
	return ids
}
