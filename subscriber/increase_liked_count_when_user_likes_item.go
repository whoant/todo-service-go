package subscriber

import (
	"context"

	"gorm.io/gorm"
	"todo-service/common"
	goservice "todo-service/go-sdk"
	"todo-service/module/item/storage"
	"todo-service/pubsub"
)

type HasItemId interface {
	GetItemId() int
}

func IncreaseLikedCountWhenUserLikesItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Increase liked_count when user likes item",
		Hdl: func(ctx context.Context, message *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
			//data := message.Data().(HasItemId)
			data := message.Data().(map[string]interface{})
			itemId := data["item_id"].(float64)

			if err := storage.NewMySQLStorage(db).IncreaseLikeCount(ctx, int(itemId)); err != nil {
				return err
			}

			_ = message.Ack()

			return nil
		},
	}
}
