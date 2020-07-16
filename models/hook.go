package models

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"math"
	"time"

	"github.com/itrix-edge/edge-client-agent/db"
	"gorm.io/gorm"
)

const HashLength = 32

type Hook struct {
	// gorm.Model
	ID                 uint           `gorm:"primarykey" json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name               string         `json:"name"`
	Key                string         `gorm:"unique_index" json:"key"`
	DeploymentOptionID uint           `json:"deplyoment_option_id"`
	// DeploymentOption DeploymentOption `json:"deployment_option"`
}

func (m Hook) randomBase16String(l int) string {
	buff := make([]byte, int(math.Round(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l]
}

// BeforeSave create key automatic
func (m *Hook) BeforeSave(gorm *gorm.DB) (err error) {
	if !(m.ID != 0) {
		m.Key = m.randomBase16String(HashLength)
	}
	return
}

// BeforeUpdate remove custom key
func (m *Hook) BeforeUpdate(gorm *gorm.DB) (err error) {
	// Omit key string changes.
	m.Key = ""
	return
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

func (m HookModel) ReadHookByKey(key string) *Hook {
	var hook = new(Hook)
	m.GetORM()
	orm.Model(&hook).Where("key = ?", key).First(&hook)
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
func (m HookModel) ExecuteHook(hook *Hook, options []OptionTemplate) bool {
	m.GetExecutionModels()
	depStatus, svcStatus, err1, err2 := deploymentOptionModel.ExecuteDeploymentByID(hook.DeploymentOptionID, options)
	if err1 != nil || err2 != nil {
		log.Fatal(err1)
		log.Fatal(err2)
		return false
	}
	log.Print("Successfuly execute deployment: " + depStatus.ObjectMeta.SelfLink)
	log.Print("Successfuly execute service: " + svcStatus.ObjectMeta.SelfLink)
	log.Print(depStatus)
	log.Print(svcStatus)
	return true
}
func (m HookModel) ExecuteHookByID(id int64, options []OptionTemplate) bool {
	hook := m.ReadHook(id)
	return m.ExecuteHook(hook, options)
}

func (m HookModel) ExecuteHookByKey(key string, options []OptionTemplate) bool {
	hook := m.ReadHookByKey(key)
	return m.ExecuteHook(hook, options)
}
