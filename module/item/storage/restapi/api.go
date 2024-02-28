package restapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"todo-service/go-sdk/logger"
)

type itemService struct {
	client     *resty.Client
	serviceURL string
	logger     logger.Logger
}

func New(serviceUrl string, logger logger.Logger) *itemService {
	return &itemService{
		client:     resty.New(),
		serviceURL: serviceUrl,
		logger:     logger,
	}
}

func (service *itemService) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	type requestBody struct {
		Ids []int `json:"ids"`
	}

	var response struct {
		Data map[int]int `json:"data"`
	}

	resp, err := service.client.R().SetHeader("Content-Type", "application/json").
		SetBody(requestBody{Ids: ids}).SetResult(&response).
		Post(fmt.Sprintf("%s/%s", service.serviceURL, "v1/rpc/get_item_likes"))
	if err != nil {
		service.logger.Errorln(err)
		return nil, err
	}

	if !resp.IsSuccess() {
		service.logger.Errorln(resp.RawResponse)
		return nil, errors.New("cannot call api get item likes")
	}

	return response.Data, nil
}
