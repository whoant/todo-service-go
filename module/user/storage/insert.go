package storage

import (
	"context"

	"todo-service/common"
	"todo-service/module/user/model"
)

func (s *mysqlStorage) CreateUser(ctx context.Context, data *model.UserCreation) error {
	db := s.db.Begin()

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback()
		return common.ErrorDB(err)
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return common.ErrorDB(err)
	}

	return nil
}
