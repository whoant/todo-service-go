package biz

import (
	"context"
	"log"

	"todo-service/common"
	"todo-service/module/user_like_item/model"
	"todo-service/pubsub"
)

type UserLikeItemStore interface {
	Create(ctx context.Context, data *model.Like) error
}

//type IncreaseItemStore interface {
//	IncreaseLikeCount(ctx context.Context, id int) error
//}

type userLikeItemBiz struct {
	store UserLikeItemStore
	//itemStore IncreaseItemStore
	ps pubsub.PubSub
}

func NewUserLikeItemBiz(store UserLikeItemStore, ps pubsub.PubSub) *userLikeItemBiz {
	return &userLikeItemBiz{
		store: store,
		ps:    ps,
	}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *model.Like) error {
	if err := biz.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}

	//go func() {
	//	defer common.Recovery()
	//
	//	if err := biz.itemStore.IncreaseLikeCount(ctx, data.ItemId); err != nil {
	//		log.Println(err)
	//	}
	//}()

	if err := biz.ps.Publish(ctx, common.TopicUserLikedItem, pubsub.NewMessage(data)); err != nil {
		log.Println(err)
	}

	return nil
}
