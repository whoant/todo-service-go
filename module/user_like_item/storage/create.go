package storage

import (
	"context"

	"todo-service/common"
	"todo-service/module/user_like_item/model"
)

func (s *sqlStore) Create(ctx context.Context, data *model.Like) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrorDB(err)
	}

	return nil
}
