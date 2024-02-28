package ginitem

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-service/common"
	"todo-service/demogrpc/demo"
	goservice "todo-service/go-sdk"
	todoBiz "todo-service/module/item/biz"
	todoModel "todo-service/module/item/model"
	itemRepo "todo-service/module/item/repo"
	todoStorage "todo-service/module/item/storage"
	"todo-service/module/item/storage/rpc"
)

func ListItemHandler(serviceCtx goservice.ServiceContext, client demo.ItemLikeServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorInvalidRequest(err))
			return
		}

		paging.Process()

		var filter todoModel.Filter
		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrorInvalidRequest(err))
			return
		}
		//apiItemCaller := serviceCtx.MustGet(common.PluginItemAPI).(interface {
		//	GetServiceUrl() string
		//})
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		storage := todoStorage.NewMySQLStorage(db)
		likeStore := rpc.NewClient(client)
		//likeStore := likeStorage.NewSQLStore(db)
		//likeStore := restapi.New(apiItemCaller.GetServiceUrl(), serviceCtx.Logger("restapi.itemlikes"))
		repo := itemRepo.NewListItemRepo(storage, likeStore, requester)
		biz := todoBiz.NewListItemsBiz(repo, requester)
		data, err := biz.ListItem(c.Request.Context(), &filter, &paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for i := range data {
			data[i].Mask()
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(data, paging, filter))
	}
}
