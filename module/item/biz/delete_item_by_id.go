package todobiz

import (
	"context"

	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

type DeleteItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*todomodel.TodoItem, error)
	DeleteItem(ctx context.Context, cond map[string]interface{}) error
}

type deleteItemBiz struct {
	storage   DeleteItemStorage
	requester common.Requester
}

func NewDeleteItemBiz(storage DeleteItemStorage, requester common.Requester) *deleteItemBiz {
	return &deleteItemBiz{
		storage:   storage,
		requester: requester,
	}
}

func (biz *deleteItemBiz) DeleteItemById(ctx context.Context, id int) error {
	data, err := biz.storage.GetItem(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		if err == common.ErrRecordNotFound {
			return common.ErrCannotGetEntity(todomodel.EntityName, err)
		}
		return common.ErrCannotDeletedEntity(todomodel.EntityName, err)
	}

	isOwner := biz.requester.GetUserId() == data.UserId
	if !isOwner && !common.IsAdmin(biz.requester) {
		return common.ErrNoPermission(err)
	}

	if data != nil && *data.Status == todomodel.ItemStatusDeleted {
		return common.ErrEntityDeleted(todomodel.EntityName, todomodel.ErrItemIsDeleted)
	}

	if err = biz.storage.DeleteItem(ctx, map[string]interface{}{
		"id": id,
	}); err != nil {
		return common.ErrCannotDeletedEntity(todomodel.EntityName, err)
	}

	return nil
}
