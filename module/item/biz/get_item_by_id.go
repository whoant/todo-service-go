package todobiz

import (
	"context"

	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

type GetItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*todomodel.TodoItem, error)
}

type getItemBiz struct {
	storage   GetItemStorage
	requester common.Requester
}

func NewGetItemBiz(storage GetItemStorage, requester common.Requester) *getItemBiz {
	return &getItemBiz{
		storage:   storage,
		requester: requester,
	}
}

func (biz *getItemBiz) GetItemById(ctx context.Context, id int) (*todomodel.TodoItem, error) {
	data, err := biz.storage.GetItem(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, common.ErrCannotGetEntity(todomodel.EntityName, err)
	}

	isOwner := biz.requester.GetUserId() == data.UserId

	if !isOwner && !common.IsAdmin(biz.requester) {
		return nil, common.ErrNoPermission(err)
	}

	return data, err
}
