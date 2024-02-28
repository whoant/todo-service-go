package repo

import (
	"context"

	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

type ListItemStorage interface {
	ListItem(ctx context.Context,
		filter *todomodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]todomodel.TodoItem, error)
}

type ItemLikeStorage interface {
	GetItemLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type listItemRepo struct {
	store     ListItemStorage
	likeStore ItemLikeStorage
	requester common.Requester
}

func NewListItemRepo(store ListItemStorage, likeStore ItemLikeStorage, requester common.Requester) *listItemRepo {
	return &listItemRepo{
		store:     store,
		likeStore: likeStore,
		requester: requester,
	}
}

func (repo *listItemRepo) ListItem(ctx context.Context, filter *todomodel.Filter, paging *common.Paging, moreKey ...string) ([]todomodel.TodoItem, error) {
	ctxStore := context.WithValue(ctx, common.CurrentUser, repo.requester)
	data, err := repo.store.ListItem(ctxStore, filter, paging, moreKey...)
	if err != nil {
		return nil, common.ErrCannotListEntity(todomodel.EntityName, err)
	}

	if len(data) == 0 {
		return data, nil
	}

	ids := make([]int, len(data))
	for i := range ids {
		ids[i] = data[i].Id
	}

	likeUserMap, err := repo.likeStore.GetItemLikes(ctxStore, ids)
	if err != nil {
		return data, err
	}

	for i := range data {
		data[i].LikedCount = likeUserMap[data[i].Id]
	}

	return data, err
}
