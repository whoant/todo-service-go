package rpc

import (
	"context"

	"todo-service/demogrpc/demo"
)

type ItemStorage interface {
	GetItemLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type rpc struct {
	store ItemStorage
}

func NewRPCService(store ItemStorage) demo.ItemLikeServiceServer {
	return &rpc{
		store: store,
	}
}

func (c *rpc) GetItemLikes(ctx context.Context, req *demo.GetItemLikesReq) (*demo.ItemLikeResp, error) {
	ids := make([]int, len(req.Ids))
	for i := range ids {
		ids[i] = int(req.Ids[i])
	}

	result, err := c.store.GetItemLikes(ctx, ids)
	if err != nil {
		return nil, err
	}

	rs := make(map[int32]int32)

	for k, v := range result {
		rs[int32(k)] = int32(v)
	}
	return &demo.ItemLikeResp{
		Result: rs,
	}, nil
}
