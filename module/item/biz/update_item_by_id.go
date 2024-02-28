package todobiz

import (
	"context"

	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

type UpdateItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*todomodel.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *todomodel.TodoItemUpdate) error
}

type updateItemBiz struct {
	storage   UpdateItemStorage
	requester common.Requester
}

func NewUpdateItemBiz(storage UpdateItemStorage, requester common.Requester) *updateItemBiz {
	return &updateItemBiz{
		storage:   storage,
		requester: requester,
	}
}

func (biz *updateItemBiz) UpdateItemById(ctx context.Context, id int, dataUpdate *todomodel.TodoItemUpdate) error {
	data, err := biz.storage.GetItem(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		if err == common.ErrRecordNotFound {
			return common.ErrCannotGetEntity(todomodel.EntityName, err)
		}
		return common.ErrCannotUpdateEntity(todomodel.EntityName, err)
	}

	isOwner := biz.requester.GetUserId() == data.UserId

	if !isOwner && !common.IsAdmin(biz.requester) {
		return common.ErrNoPermission(err)
	}

	if data.Status != nil && *data.Status == todomodel.ItemStatusDeleted {
		return common.ErrEntityDeleted(todomodel.EntityName, todomodel.ErrItemIsDeleted)
	}

	if err = biz.storage.UpdateItem(ctx, map[string]interface{}{
		"id": id,
	}, dataUpdate); err != nil {
		return common.ErrCannotUpdateEntity(todomodel.EntityName, err)
	}

	return nil
}
