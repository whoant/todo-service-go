package todobiz

import (
	"context"

	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

type CreateTodoItemStorage interface {
	CreateItem(ctx context.Context, data *todomodel.TodoItemCreation) error
}

type createBiz struct {
	store CreateTodoItemStorage
}

func NewCreateToDoItemBiz(store CreateTodoItemStorage) *createBiz {
	return &createBiz{
		store: store,
	}
}

func (biz *createBiz) CreateNewItem(ctx context.Context, data *todomodel.TodoItemCreation) error {
	if data.Title == "" {
		return common.ErrorInvalidRequest(todomodel.ErrTitleCannotBeEmpty)
	}

	if err := biz.store.CreateItem(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(todomodel.EntityName, err)
	}

	return nil
}
