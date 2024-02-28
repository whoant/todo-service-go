package storage

import (
	"context"

	"gorm.io/gorm"
	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

func (s *mysqlStorage) UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *todomodel.TodoItemUpdate) error {
	if err := s.db.Where(cond).Updates(&dataUpdate).Error; err != nil {
		return common.ErrorDB(err)
	}

	return nil
}

func (s *mysqlStorage) IncreaseLikeCount(ctx context.Context, id int) error {
	db := s.db
	if err := db.Table(todomodel.TodoItem{}.TableName()).Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count + ?", 1)).Error; err != nil {
		return common.ErrorDB(err)
	}

	return nil
}

func (s *mysqlStorage) DecreaseLikeCount(ctx context.Context, id int) error {
	db := s.db
	if err := db.Table(todomodel.TodoItem{}.TableName()).Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count - ?", 1)).Error; err != nil {
		return common.ErrorDB(err)
	}

	return nil
}
