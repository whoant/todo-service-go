package biz

import (
	"context"

	"todo-service/common"
	"todo-service/module/user/model"
)

type RegisterStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
	CreateUser(ctx context.Context, data *model.UserCreation) error
}

type Hasher interface {
	Hash(data string) string
}

type registerBusiness struct {
	storage RegisterStorage
	hasher  Hasher
}

func NewRegisterBiz(storage RegisterStorage, hasher Hasher) *registerBusiness {
	return &registerBusiness{storage, hasher}
}

func (biz *registerBusiness) Register(ctx context.Context, data *model.UserCreation) error {
	user, _ := biz.storage.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if user != nil {
		return model.ErrEmailExisted
	}

	salt := common.GenSalt(50)
	data.Password = biz.hasher.Hash(data.Password + salt)
	data.Salt = salt
	data.Role = "user"

	if err := biz.storage.CreateUser(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	return nil
}
