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

func DeleteItemHandler(serviceCtx goservice.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorInvalidRequest(err))

			return
		}
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		storage := todostorage.NewMySQLStorage(db)
		biz := todobiz.NewDeleteItemBiz(storage, requester)
		if err := biz.DeleteItemById(c.Request.Context(), int(uid.GetLocalID())); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
