package todobiz

import (
	"context"

	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

type ListItemRepo interface {
	ListItem(ctx context.Context,
		filter *todomodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]todomodel.TodoItem, error)
}

type listItemBiz struct {
	repo      ListItemRepo
	requester common.Requester
}

func NewListItemsBiz(repo ListItemRepo, requester common.Requester) *listItemBiz {
	return &listItemBiz{
		repo:      repo,
		requester: requester,
	}
}

func (biz *listItemBiz) ListItem(ctx context.Context,
	filter *todomodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]todomodel.TodoItem, error) {
	ctxStore := context.WithValue(ctx, common.CurrentUser, biz.requester)

	data, err := biz.repo.ListItem(ctxStore, filter, paging, "Owner")
	if err != nil {
		return nil, common.ErrCannotListEntity(todomodel.EntityName, err)
	}

	return data, nil
}
