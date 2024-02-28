package ginitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-service/common"
	goservice "todo-service/go-sdk"
	todobiz "todo-service/module/item/biz"
	todostorage "todo-service/module/item/storage"
)

func GetItemHandler(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorInvalidRequest(err))
			return
		}
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		storage := todostorage.NewMySQLStorage(db)
		biz := todobiz.NewGetItemBiz(storage, requester)
		data, err := biz.GetItemById(c.Request.Context(), int(uid.GetLocalID()))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		data.Mask()
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
