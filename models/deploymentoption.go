package models

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"reflect"
	"text/template"
	"time"

	"github.com/itrix-edge/edge-client-agent/db"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"gorm.io/gorm"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type DeploymentOptionModel struct{}

type OptionTemplate struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DeploymentOption struct {
	ID                           uint              `gorm:"primarykey" json:"id"`
	CreatedAt                    time.Time         `json:"created_at"`
	UpdatedAt                    time.Time         `json:"updated_at"`
	DeletedAt                    gorm.DeletedAt    `gorm:"index" json:"deleted_at"`
	Namespace                    string            `json:"namespace"`
	Options                      []OptionTemplate  `json:"options" gorm:"-"`
	OptionsSerialized            string            `json:"-"`
	DeploymentTemplate           appsv1.Deployment `json:"deployment_template" gorm:"-"`
	DeploymentTemplateSerialized string            `json:"-"` // appsv1.Deployment
	ServiceTemplate              corev1.Service    `json:"service_template" gorm:"-"`
	ServiceTemplateSerialized    string            `json:"-"` // corev1.Service
	Hooks                        []Hook            `json:"hooks"`
	// Image              string            `json:"image"`
	// Name               string            `json:"name"`
	// Ports              []int             `json:"ports"`
}

// BeforeSave Handle option struct serialized to json
func (u *DeploymentOption) BeforeSave(gorm *gorm.DB) (err error) {
	u.OptionsSerialized = u.SerializeValue(u.Options)
	u.DeploymentTemplateSerialized = u.SerializeValue(u.DeploymentTemplate)
	u.ServiceTemplateSerialized = u.SerializeValue(u.ServiceTemplate)
	return
}

// SerializeValue make struct serialized to json
func (u *DeploymentOption) SerializeValue(v interface{}) string {
	byteArray, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(byteArray)
}

// UnSerializeValue make serialized json to struct
func (u *DeploymentOption) UnSerializeValue(values string, v interface{}) {
	if err := json.Unmarshal([]byte(values), v); err != nil {
		log.Fatal(err)
	}
}

// AfterFind Handle serialized to json optionserialized to struct
func (u *DeploymentOption) AfterFind(gorm *gorm.DB) (err error) {
	u.UnSerializeValue(u.OptionsSerialized, &u.Options)
	u.UnSerializeValue(u.DeploymentTemplateSerialized, &u.DeploymentTemplate)
	u.UnSerializeValue(u.ServiceTemplateSerialized, &u.ServiceTemplate)
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
