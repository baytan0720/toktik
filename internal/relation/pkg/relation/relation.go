package relation

import (
	"gorm.io/gorm"

	"toktik/pkg/db/model"
)

type Relation = model.Relation

type RelationService struct {
	dbInstance func() *gorm.DB
}

func NewRelationService(db func() *gorm.DB) *RelationService {
	return &RelationService{
		dbInstance: db,
	}
}
