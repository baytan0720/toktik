package feed

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	
	"toktik/pkg/db/model"
	"toktik/pkg/test/testutil"
)

func newMockDB(t *testing.T) *gorm.DB {
	db := testutil.NewMockDB()
	require.NoError(t, db.AutoMigrate(&model.Video{}))
	return db
}

//func TestRelationOperator_Feed(t *testing.T) {
//	db := newMockDB(t)
//	f := NewFeedService(func() *gorm.DB {
//		return db
//	})
//
//	testFeedCase := &model.Video{
//
//	}
//}
