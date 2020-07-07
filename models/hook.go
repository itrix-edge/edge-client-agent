package models

import (
	"github.com/stevennick/edge-client-agent/db"
	"gorm.io/gorm"
)

type Hook struct {
	gorm.Model
	ID           int64  `db:"id, primarykey, autoincrement" json:"id"`
	Name         string `db:"name" json:"name"`
	Key          string `db:"string" json:"string"`
	DeploymentID int64  `db:"deployment_id" json:"deplyoment_id"`
	UpdatedAt    int64  `db:"updated_at" json:"-"`
	CreatedAt    int64  `db:"created_at" json:"-"`
}

// HookModel
type HookModel struct{}

var orm *gorm.DB

func init() {
	orm := db.GetORM()
	if orm != nil {
		orm.AutoMigrate(&Hook{})
	}
}

func (m HookModel) CreateHook(hook *Hook) *Hook {
	orm.Create(&hook)
	return hook
}

func (m HookModel) ReadHook(id int64) *Hook {
	var hook = new(Hook)
	orm.First(&hook, id)
	return hook
}

func (m HookModel) UpdateHook(hook *Hook) *Hook {
	orm.Model(&hook).Updates(&hook)
	return hook
}

func (m HookModel) DeleteHook(id int64) bool {
	hook := m.ReadHook(id)
	orm.Delete(&hook)
	return true
}
