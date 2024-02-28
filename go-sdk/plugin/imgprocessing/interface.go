package imgprocessing

import (
	"mime/multipart"
	"todo-service/go-sdk/sdkcm"
)

type ImgProcessing interface {
	// call img processing service to resize img and upload to s3
	Resize(file *multipart.FileHeader, folder string, longEdge int, quality int) (*sdkcm.Image, error)

	ResizeFile(filePath string, folder string, longEdge int, quality int) (*sdkcm.Image, error)
}
