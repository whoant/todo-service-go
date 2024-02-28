package ginitem

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-service/common"
	goservice "todo-service/go-sdk"
	todobiz "todo-service/module/item/biz"
	todomodel "todo-service/module/item/model"
	todostorage "todo-service/module/item/storage"
)

func CreateNewItemHandler(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataItem todomodel.TodoItemCreation

		if err := c.ShouldBind(&dataItem); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorInvalidRequest(err))
			return
		}

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		dataItem.UserId = requester.GetUserId()
		dataItem.Title = strings.TrimSpace(dataItem.Title)

		storage := todostorage.NewMySQLStorage(db)
		biz := todobiz.NewCreateToDoItemBiz(storage)

		if err := biz.CreateNewItem(c.Request.Context(), &dataItem); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": dataItem.Id})
	}
}
