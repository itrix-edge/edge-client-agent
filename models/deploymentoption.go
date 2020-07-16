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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Hooks                        []Hook            `gorm:"foreignkey:DeploymentOptionID" json:"hooks"`
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

// AfterSave apply new hook into this deployment option
func (u *DeploymentOption) AfterSave(gorm *gorm.DB) (err error) {
	if len(u.Hooks) == 0 {
		var hook = Hook{Name: u.Namespace + "." + u.DeploymentTemplate.Name, DeploymentOptionID: u.ID}
		gorm.Save(&hook)
		u.Hooks = append(u.Hooks, hook)
	}
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
// apply related hook into this deployment option
func (u *DeploymentOption) AfterFind(gorm *gorm.DB) (err error) {
	u.UnSerializeValue(u.OptionsSerialized, &u.Options)
	u.UnSerializeValue(u.DeploymentTemplateSerialized, &u.DeploymentTemplate)
	u.UnSerializeValue(u.ServiceTemplateSerialized, &u.ServiceTemplate)
	gorm.Where("deployment_option_id = ?", u.ID).Find(&u.Hooks)
	return
}

var deploymentModel *DeploymentModel
var serviceModel *ServiceModel
var namespaceModel *NamespaceModel

func (m DeploymentOptionModel) GetExecutionModels() {
	if deploymentModel == nil {
		deploymentModel = new(DeploymentModel)
	}
	if serviceModel == nil {
		serviceModel = new(ServiceModel)
	}
	if namespaceModel == nil {
		namespaceModel = new(NamespaceModel)
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

func (m DeploymentOptionModel) ExecuteDeploymentByID(id uint, options []OptionTemplate) (*appsv1.Deployment, *corev1.Service, error, error) {
	// Init related models
	m.GetExecutionModels()
	m.GetORM()

	// fetch DeploymentOption
	deployOption := m.GetDeploymentOptionByID(id)

	// Apply Template to Deployment & Service
	var deplyomentTemplate []byte
	var serviceTemplate []byte
	if len(options) > 0 {
		// if len(options) > 0 || len(deployOption.Options) > 0 {
		var depBuf, svcBuf bytes.Buffer
		if options != nil {
			depBuf = m.ApplyTemplate("Deployment", deployOption.DeploymentTemplateSerialized, options)
			svcBuf = m.ApplyTemplate("Service", deployOption.ServiceTemplateSerialized, options)
			// } else {
			// 	// Use default options array
			// 	depBuf = m.ApplyTemplate("Deployment", deployOption.DeploymentTemplateSerialized, deployOption.Options)
			// 	svcBuf = m.ApplyTemplate("Service", deployOption.ServiceTemplateSerialized, deployOption.Options)
		}
		deplyomentTemplate = depBuf.Bytes()
		serviceTemplate = svcBuf.Bytes()
	} else {
		deplyomentTemplate = []byte(deployOption.DeploymentTemplateSerialized)
		serviceTemplate = []byte(deployOption.ServiceTemplateSerialized)
	}

	// Encoding to struct
	deployment := appsv1.Deployment{}
	json.Unmarshal(deplyomentTemplate, &deployment)
	service := corev1.Service{}
	json.Unmarshal(serviceTemplate, &service)

	// Query namespace and create namespace if required
	ns, err := namespaceModel.GetNamespace(deployOption.Namespace, v1.GetOptions{})
	if err != nil || ns == nil {
		var namespace = corev1.Namespace{}
		namespace.Name = deployOption.Namespace
		namespace.Namespace = deployOption.Namespace
		ns, err = namespaceModel.CreateNamespace(&namespace)
	}

	// Apply deployment to cluster
	deployresult, deperr := deploymentModel.CreateDeployment(deployOption.Namespace, &deployment)
	serviceresult, svcerr := serviceModel.CreateService(deployOption.Namespace, &service)
	// serviceModel.CreateService(deployOption.Namespace, &service)
	return deployresult, serviceresult, deperr, svcerr
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
