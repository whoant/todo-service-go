package biz

import (
	"context"

	"todo-service/common"
	"todo-service/module/user/model"
	token_provider "todo-service/plugin/token_provider"
)

type loginStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
}

type loginBiz struct {
	storage       loginStorage
	hasher        Hasher
	tokenProvider token_provider.Provider
	expiry        int
}

func NewLoginBiz(storage loginStorage, hasher Hasher, tokenProvider token_provider.Provider, expiry int) *loginBiz {
	return &loginBiz{storage: storage, hasher: hasher, tokenProvider: tokenProvider, expiry: expiry}
}

func (biz *loginBiz) Login(ctx context.Context, data *model.UserLogin) (token_provider.Token, error) {
	user, err := biz.storage.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if err != nil {
		return nil, model.ErrEmailOrPasswordInvalid
	}

	passwordHashed := biz.hasher.Hash(data.Password + user.Salt)
	if passwordHashed != user.Password {
		return nil, model.ErrEmailOrPasswordInvalid
	}

	payload := &common.TokenPayload{
		UId:   user.Id,
		URole: user.Role.String(),
	}

	accessToken, err := biz.tokenProvider.Generate(payload, biz.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	return accessToken, nil
}
