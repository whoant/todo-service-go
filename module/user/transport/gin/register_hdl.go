package gin_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-service/common"
	goservice "todo-service/go-sdk"
	"todo-service/module/user/biz"
	"todo-service/module/user/model"
	"todo-service/module/user/storage"
)

func Register(serviceCtx goservice.ServiceContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var data model.UserCreation

		if err := ctx.ShouldBindJSON(&data); err != nil {
			panic(common.ErrorInvalidRequest(err))
		}

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		store := storage.NewMySQLStorage(db)
		md5 := common.NewMd5Hash()
		business := biz.NewRegisterBiz(store, md5)
		if err := business.Register(ctx.Request.Context(), &data); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
