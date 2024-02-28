package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-service/common"
	goservice "todo-service/go-sdk"
	"todo-service/module/user_like_item/biz"
	"todo-service/module/user_like_item/model"
	"todo-service/module/user_like_item/storage"
	"todo-service/pubsub"
)

func LikeItem(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrorInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		store := storage.NewSQLStore(db)
		//itemStore := itemStorage.NewMySQLStorage(db)
		pb := serviceCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)
		business := biz.NewUserLikeItemBiz(store, pb)

		if err := business.LikeItem(c.Request.Context(), &model.Like{
			UserId: requester.GetUserId(),
			ItemId: int(id.GetLocalID()),
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
