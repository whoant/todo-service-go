package storage

import (
	"context"

	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

func (s *mysqlStorage) CreateItem(ctx context.Context, data *todomodel.TodoItemCreation) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrorDB(err)
	}

	return nil
}
