package gin_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"todo-service/common"
)

func GetUser() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		u := ctx.MustGet(common.CurrentUser)

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(u))
	}
}
