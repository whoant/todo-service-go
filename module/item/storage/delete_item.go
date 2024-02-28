package storage

import (
	"context"
	"time"

	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

func (s *mysqlStorage) DeleteItem(ctx context.Context, cond map[string]interface{}) error {

	if err := s.db.Table(todomodel.TodoItem{}.TableName()).
		Where(cond).
		Updates(map[string]interface{}{
			"status":     todomodel.ItemStatusDeleted,
			"updated_at": time.Now().UTC(),
		}).Error; err != nil {
		return common.ErrorDB(err)
	}

	return nil
}
