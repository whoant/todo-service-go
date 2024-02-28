package biz

import (
	"context"
	"log"

	"todo-service/common"
	"todo-service/module/user_like_item/model"
	"todo-service/pubsub"
)

type UserUnlikeItemStore interface {
	Find(ctx context.Context, userId, itemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

//type DecreaseItemStore interface {
//	DecreaseLikeCount(ctx context.Context, id int) error
//}

type userUnlikeItemBiz struct {
	store UserUnlikeItemStore
	//itemStore DecreaseItemStore
	ps pubsub.PubSub
}

func NewUserUnlikeItemBiz(store UserUnlikeItemStore, ps pubsub.PubSub) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{
		store: store,
		ps:    ps,
	}
}

func (biz *userUnlikeItemBiz) UnlikeItem(ctx context.Context, userId, itemId int) error {
	_, err := biz.store.Find(ctx, userId, itemId)
	if err == common.ErrRecordNotFound {
		return model.ErrDidNotLikeItem(err)
	}

	if err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	if err := biz.store.Delete(ctx, userId, itemId); err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	//go func() {
	//	defer common.Recovery()
	//	if err := biz.itemStore.DecreaseLikeCount(ctx, itemId); err != nil {
	//		log.Println(err)
	//	}
	//}()

	if err := biz.ps.Publish(ctx, common.TopicUserUnlikedItem, pubsub.NewMessage(&model.Like{
		UserId: userId,
		ItemId: itemId,
	})); err != nil {
		log.Println(err)
	}

	return nil
}
