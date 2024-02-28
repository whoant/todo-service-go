package ginitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-service/common"
	goservice "todo-service/go-sdk"
	todobiz "todo-service/module/item/biz"
	todomodel "todo-service/module/item/model"
	todostorage "todo-service/module/item/storage"
)

func UpdateItemHandler(serviceCtx goservice.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorInvalidRequest(err))
			return
		}

		var data todomodel.TodoItemUpdate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorInvalidRequest(err))
			return
		}
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		storage := todostorage.NewMySQLStorage(db)
		biz := todobiz.NewUpdateItemBiz(storage, requester)
		if err := biz.UpdateItemById(c.Request.Context(), int(uid.GetLocalID()), &data); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
