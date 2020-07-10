package models

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"reflect"
	"text/template"
	"time"

	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"github.com/stevennick/edge-client-agent/db"
	"gorm.io/gorm"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type DeploymentOptionModel struct{}

type OptionTemplate struct {
	Key   string `json:"name=key"`
	Value string `json:"name=value"`
}

type DeploymentOption struct {
	ID                           uint              `gorm:"primarykey" json:"name=id"`
	CreatedAt                    time.Time         `json:"name=created_at"`
	UpdatedAt                    time.Time         `json:"name=updated_at"`
	DeletedAt                    gorm.DeletedAt    `gorm:"index" json:"name=deleted_at"`
	Namespace                    string            `json:"name=namespace"`
	Options                      []OptionTemplate  `json:"name=options" gorm:"-"`
	OptionsSerialized            string            `json:"-"`
	DeploymentTemplate           appsv1.Deployment `json:"name=deployment_template" gorm:"-"`
	DeploymentTemplateSerialized string            `json:"-"` // appsv1.Deployment
	ServiceTemplate              corev1.Service    `json:"name=service_template" gorm:"-"`
	ServiceTemplateSerialized    string            `json:"-"` // corev1.Service
	// Image              string            `json:"name=image"`
	// Name               string            `json:"name=name"`
	// Ports              []int             `json:"name=ports"`
}

// BeforeSave Handle option struct serialized to json
func (u *DeploymentOption) BeforeSave(gorm *gorm.DB) (err error) {
	serialized, err := json.Marshal(u.Options)
	if err != nil {
		return err
	}
	u.OptionsSerialized = string(serialized)
	return
}

// AfterFind Handle serialized to json optionserialized to struct
func (u *DeploymentOption) AfterFind(gorm *gorm.DB) (err error) {
	var val []OptionTemplate
	if err := json.Unmarshal([]byte(u.OptionsSerialized), &val); err != nil {
		return err
	}
	u.Options = val
	return
}

var deploymentModel *DeploymentModel

func (m DeploymentOptionModel) GetExecutionModels() {
	if deploymentModel == nil {
		deploymentModel = new(DeploymentModel)
	}
}

func (m DeploymentOptionModel) GetORM() {
	if orm == nil {
		orm = db.GetORM()
	}
}

func (m DeploymentOptionModel) ListDeploymentOptions() *[]DeploymentOption {
	var options []DeploymentOption
	m.GetORM()
	orm.Find(&options)
	return &options
}

func (m DeploymentOptionModel) CreateDeploymentOption(deploy *DeploymentOption) *DeploymentOption {
	m.GetORM()
	orm.Create(&deploy)
	return deploy
}

func (m DeploymentOptionModel) GetDeploymentOptionByID(id uint) *DeploymentOption {
	m.GetORM()
	var deploy = new(DeploymentOption)
	deploy.ID = id
	orm.First(&deploy)
	return deploy
}

func (m DeploymentOptionModel) UpdateDeploymentOptionByID(id uint, deploy *DeploymentOption) *DeploymentOption {
	m.GetORM()
	deploy.ID = id
	orm.Updates(&deploy)
	return deploy
}

func (m DeploymentOptionModel) DeleteDeploymentOptionByID(id uint) bool {
	m.GetORM()
	deploy := m.GetDeploymentOptionByID(id)
	orm.Delete(&deploy)
	return true
}

func (m DeploymentOptionModel) Migrate() {
	m.GetORM()
	orm.AutoMigrate(&DeploymentOption{})
	if !orm.Migrator().HasTable(&DeploymentOption{}) {
		orm.Migrator().CreateTable(&DeploymentOption{})
	}
}

func (m DeploymentOptionModel) ExecuteDeploymentByID(id uint) (*appsv1.Deployment, error) {
	// Init related models
	m.GetExecutionModels()
	m.GetORM()

	// fetch DeploymentOption
	deployOption := m.GetDeploymentOptionByID(id)

	// Apply Template
	// t := template.Must(template.New("Deployment").Parse(deployOption.DeploymentTemplate))

	buf := m.ApplyTemplate("Deployment", deployOption.DeploymentTemplateSerialized, deployOption.Options)

	// Encoding to struct
	deployment := appsv1.Deployment{}
	json.Unmarshal(buf.Bytes(), &deployment)

	// Apply deployment to cluster
	deployresult, err := deploymentModel.CreateDeployment(deployOption.Namespace, &deployment)
	return deployresult, err
}

// ApplyTemplate apply options to the target article
func (m DeploymentOptionModel) ApplyTemplate(name string, rawtemplate string, replacement []OptionTemplate) bytes.Buffer {
	dynstructBuilder := dynamicstruct.NewStruct()
	// valarray := make(map[int]interface{})

	for _, r := range replacement {
		dynstructBuilder.AddField(r.Key, r.Value, "")
		// valarray[i] = r.Value
	}
	dynstruct := dynstructBuilder.Build().New()
	elem := reflect.ValueOf(&dynstruct).Elem()
	for _, r := range replacement {
		value := reflect.ValueOf(r.Value)
		elem.FieldByName(r.Key).Set(value)
	}

	t := template.Must(template.New(name).Parse(rawtemplate))
	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := t.Execute(w, dynstruct)
	if err != nil {
		log.Fatal("Execute template error:", err)
	}
	return buf
}
