package storage

import (
	"context"

	"gorm.io/gorm"
	"todo-service/common"
	"todo-service/module/user_like_item/model"
)

func (s *sqlStore) Find(ctx context.Context, userId, itemId int) (*model.Like, error) {
	var data model.Like

	if err := s.db.Where("user_id = ? and item_id = ?", userId, itemId).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}

		return nil, common.ErrorDB(err)
	}

	return &data, nil
}
