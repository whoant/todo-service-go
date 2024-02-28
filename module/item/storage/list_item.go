package storage

import (
	"context"

	"todo-service/common"
	todomodel "todo-service/module/item/model"
)

func (s *mysqlStorage) ListItem(ctx context.Context,
	filter *todomodel.Filter,
	paging *common.Paging,
	moreKeys ...string) ([]todomodel.TodoItem, error) {

	db := s.db.Where("status <> ?", "Deleted")

	if f := filter; f != nil {
		if v := f.Status; v != "" {
			db = db.Where("status = ?", v)
		}
	}

	if err := db.Select("id").Table(todomodel.TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	var result []todomodel.TodoItem

	if v := paging.FakeCursor; v != "" {
		uid, err := common.FromBase58(v)
		if err != nil {
			return nil, common.ErrorDB(err)
		}
		db = db.Where("id < ?", uid.GetLocalID())
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.Select("*").Order("id desc").
		Limit(paging.Limit).
		Find(&result).Error; err != nil {
		return nil, common.ErrorDB(err)
	}

	if len(result) > 0 {
		result[len(result)-1].Mask()
		paging.NextCursor = result[len(result)-1].FakeId.String()
	}

	return result, nil
}
