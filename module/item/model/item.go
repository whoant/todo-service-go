package model

import (
	"errors"
	"strings"
	"time"

	"todo-service/common"
)

const (
	EntityName = "item"
)

var (
	ErrTitleCannotBeEmpty = errors.New("title cannot be black")
	ErrItemIsDeleted      = errors.New("item is deleted")
)

type TodoItem struct {
	common.SQLModel
	Title       string             `json:"title" gorm:"column:title;"`
	Description string             `json:"description" gorm:"column:description;"`
	Status      *ItemStatus        `json:"status" gorm:"column:status;"`
	Image       *common.Image      `json:"image" gorm:"column:image;"`
	UserId      int                `json:"-" gorm:"column:user_id;"`
	LikedCount  int                `json:"liked_count" gorm:"-;"`
	Owner       *common.SimpleUser `json:"owner" gorm:"foreignKey:UserId;"`
}

func (TodoItem) TableName() string {
	return "todo_items"
}

func (item *TodoItem) Mask() {
	item.SQLModel.Mask(common.DbTypeItem)

	if v := item.Owner; v != nil {
		v.Mask()
	}
}

type TodoItemCreation struct {
	Id          int           `json:"_" gorm:"column:id;"`
	UserId      int           `json:"-" gorm:"column:user_id;"`
	Title       string        `json:"title" gorm:"column:title;"`
	Description string        `json:"description" gorm:"column:description;"`
	Image       *common.Image `json:"image" gorm:"column:image;"`
	CreatedAt   *time.Time    `json:"-" gorm:"column:created_at;autoCreateTime"`
}

func (t *TodoItemCreation) Validate() error {
	t.Title = strings.TrimSpace(t.Title)

	if t.Title == "" {
		return ErrTitleCannotBeEmpty
	}

	return nil
}

func (TodoItemCreation) TableName() string {
	return TodoItem{}.TableName()
}

type TodoItemUpdate struct {
	Title       *string       `json:"title" gorm:"column:title;"`
	Description *string       `json:"description" gorm:"column:description;"`
	Status      *ItemStatus   `json:"status" gorm:"column:status;"`
	Image       *common.Image `json:"image" gorm:"column:image;"`
	UpdatedAt   *time.Time    `json:"-" gorm:"column:updated_at;autoUpdateTime"`
}

func (TodoItemUpdate) TableName() string {
	return TodoItem{}.TableName()
}
