package upload

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"todo-service/common"
	goservice "todo-service/go-sdk"
)

func UploadHandler(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			panic(common.ErrorInvalidRequest(err))
		}

		dst := fmt.Sprintf("static/%d.%s", time.Now().UTC().UnixNano(), fileHeader.Filename)
		if err := ctx.SaveUploadedFile(fileHeader, dst); err != nil {

		}

		img := common.Image{
			Id:        0,
			Url:       dst,
			Width:     100,
			Height:    100,
			CloudName: "local",
			Extension: "",
		}

		img.FullFill("http://localhost:3000")

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(img))
	}
}
