package models

import (
	"log"
	"time"

	"github.com/stevennick/edge-client-agent/db"
	"gorm.io/gorm"
)

type Hook struct {
	// gorm.Model
	ID                 uint           `gorm:"primarykey" json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name               string         `json:"name"`
	Key                string         `json:"key"`
	DeploymentOptionID uint           `json:"deplyoment_option_id"`
	// DeploymentOption DeploymentOption `json:"deployment_option"`
}

// HookModel
type HookModel struct{}

var orm *gorm.DB
var deploymentOptionModel *DeploymentOptionModel

func (m HookModel) GetExecutionModels() {
	if deploymentOptionModel == nil {
		deploymentOptionModel = new(DeploymentOptionModel)
	}
}

func (m HookModel) GetORM() {
	if orm == nil {
		orm = db.GetORM()
	}
}

func (m HookModel) Migrate() {
	m.GetORM()
	orm.AutoMigrate(&Hook{})
	if !orm.Migrator().HasTable(&Hook{}) {
		orm.Migrator().CreateTable(&Hook{})
	}
}

func (m HookModel) ListHooks() *[]Hook {
	var hooks []Hook
	m.GetORM()
	orm.Find(&hooks)
	return &hooks
}

func (m HookModel) CreateHook(hook *Hook) *Hook {
	m.GetORM()
	orm.Create(&hook)
	return hook
}

func (m HookModel) ReadHook(id int64) *Hook {
	var hook = new(Hook)
	m.GetORM()
	orm.First(&hook, id)
	return hook
}

func (m HookModel) UpdateHook(hook *Hook) *Hook {
	m.GetORM()
	orm.Model(&hook).Updates(&hook)
	return hook
}

func (m HookModel) DeleteHook(id int64) bool {
	m.GetORM()
	hook := m.ReadHook(id)
	orm.Delete(&hook)
	return true
}

// ExecuteHook executes predefined deeployment and its services
// 1. Get Hook obj
// 2. Get assoicated deploymentOption obj
// 3. Use deploymentOption obj to execute Deployment, Service (inside deploymentModel)
func (m HookModel) ExecuteHook(id int64) bool {
	hook := m.ReadHook(id)
	m.GetExecutionModels()
	status, err := deploymentOptionModel.ExecuteDeploymentByID(hook.DeploymentOptionID)
	if err != nil {
		log.Fatal(err)
		return false
	}
	log.Print(status)
	return true
}
