package storage

import (
	"context"

	"todo-service/common"
	"todo-service/module/user_like_item/model"
)

func (s *sqlStore) Delete(ctx context.Context, userId, itemId int) error {
	var data model.Like
	if err := s.db.Table(data.TableName()).
		Where("user_id = ? and item_id = ?", userId, itemId).
		Delete(nil).Error; err != nil {
		return common.ErrorDB(err)
	}

	return nil
}
