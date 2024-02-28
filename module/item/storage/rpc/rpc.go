package rpc

import (
	"context"

	"todo-service/demogrpc/demo"
)

type rpcClient struct {
	client demo.ItemLikeServiceClient
}

func NewClient(client demo.ItemLikeServiceClient) *rpcClient {
	return &rpcClient{
		client: client,
	}
}

func (c *rpcClient) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	reqIds := make([]int32, len(ids))

	for i := range ids {
		reqIds[i] = int32(ids[i])
	}

	resp, err := c.client.GetItemLikes(ctx, &demo.GetItemLikesReq{Ids: reqIds})
	if err != nil {
		return nil, err
	}

	res := make(map[int]int)
	for k, v := range resp.Result {
		res[int(k)] = int(v)
	}

	return res, nil
}
