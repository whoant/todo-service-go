package storage

import (
	"context"

	"go.opencensus.io/trace"
	"gorm.io/gorm"
	"todo-service/common"
	"todo-service/module/user/model"
)

func (s *mysqlStorage) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error) {
	_, span := trace.StartSpan(ctx, "user.storage.find")
	defer span.End()

	db := s.db.Table(model.User{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	var user model.User

	if err := db.Where(conditions).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.ErrorDB(err)
	}

	return &user, nil
}
