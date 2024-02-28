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
	token_provider "todo-service/plugin/token_provider"
)

func Login(serviceCtx goservice.ServiceContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var loginUserData model.UserLogin

		if err := ctx.ShouldBindJSON(&loginUserData); err != nil {
			panic(common.ErrorInvalidRequest(err))
		}

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		provider := serviceCtx.MustGet(common.PluginJWT).(token_provider.Provider)

		md5 := common.NewMd5Hash()

		store := storage.NewMySQLStorage(db)
		business := biz.NewLoginBiz(store, md5, provider, 60*60*24*30)
		account, err := business.Login(ctx.Request.Context(), &loginUserData)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
