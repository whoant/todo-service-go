package storage

import (
	"context"

	"go.opencensus.io/trace"
	"gorm.io/gorm"
	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

func (s *mysqlStorage) GetItem(ctx context.Context, cond map[string]interface{}) (*todomodel.TodoItem, error) {
	_, span := trace.StartSpan(ctx, "item.storage.find")
	defer span.End()

	var data todomodel.TodoItem

	if err := s.db.Where(cond).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}
		return nil, common.ErrorDB(err)
	}

	return &data, nil
}
