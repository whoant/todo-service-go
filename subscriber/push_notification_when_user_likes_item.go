package subscriber

import (
	"context"
	"log"

	goservice "todo-service/go-sdk"
	"todo-service/pubsub"
)

type HasUserId interface {
	GetUserId() int
}

func PushNotificationWhenUserLikesItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Push notification when user likes item",
		Hdl: func(ctx context.Context, message *pubsub.Message) error {
			//data := message.Data().(HasUserId)
			data := message.Data().(map[string]interface{})
			userId := data["user_id"].(float64)
			log.Println("Push notification to user_id:", int(userId))

			return nil
		},
	}
}
