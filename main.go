/*
Copyright 2020, Yi-Fu Ciou and the ITRIX-EDGE authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Project template original from https://github.com/Massad/gin-boilerplate

*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/gin-contrib/gzip"
	"github.com/itrix-edge/edge-client-agent/controllers"
	"github.com/itrix-edge/edge-client-agent/db"
	"github.com/itrix-edge/edge-client-agent/models"
	"github.com/joho/godotenv"

	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"

	"github.com/gin-gonic/gin"
	uuid "github.com/twinj/uuid"
)

const LocalEnv = ".env"
const SystemEnv = "/config/.env"

//CORSMiddleware ...
//CORS (Cross-Origin Resource Sharing)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

//RequestIDMiddleware ...
//Generate a unique ID and attach it to each request for future reference or use
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}

var auth = new(controllers.AuthController)

//TokenAuthMiddleware ...
//JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.TokenValid(c)
		c.Next()
	}
}

// FileExists checks env and change them between local and system level.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {

	//Start the default gin server
	r := gin.Default()
	var env string
	if FileExists(SystemEnv) {
		env = SystemEnv
	}
	if FileExists(LocalEnv) {
		env = LocalEnv
	}
	//Load the .env file
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file, please create one in the root directory")
	}

	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	dbInit, err := strconv.ParseBool(os.Getenv("DB"))
	if err != nil {
		log.Fatal("Error parse .env file with key 'DB'")
	}

	if dbInit {
		//Start PostgreSQL database
		//Example: db.GetDB() - More info in the models folder
		dbDSN := fmt.Sprintf("user=%s password=%s DB.name=%s host=%s port=%s %s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_OPTIONS"))
		debugFlag := os.Getenv("APP_DEBUG")
		debug, err := strconv.ParseBool(debugFlag)
		if err != nil {
			log.Fatal("Error parse .env file DEBUG flag parse fail.")
		}
		db.Init(dbDSN, debug)
		// db.InitORM()

		//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
		//Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()
		// db.InitRedis("1")
	}

	// Kubernetes client
	kubeconfig := os.Getenv("KUBE_CONFIG")
	if len(kubeconfig) == 0 {
		models.InitKuberClient(nil)
	} else {
		models.InitKuberClient(&kubeconfig)
	}

	v1 := r.Group("/v1")
	{
		/*** START USER ***/
		user := new(controllers.UserController)

		v1.POST("/user/login", user.Login)
		v1.POST("/user/register", user.Register)
		v1.GET("/user/logout", user.Logout)

		/*** START AUTH ***/
		auth := new(controllers.AuthController)

		//Rerfresh the token when needed to generate new access_token and refresh_token for the user
		v1.POST("/token/refresh", auth.Refresh)

		/*** Namespace ***/
		namespace := new(controllers.NamespaceController)

		v1.GET("/namespaces", namespace.GetNamespaces)

		/*** Deployment ***/

		deployment := new(controllers.DeploymentController)
		v1.GET("/deployment", deployment.GetDeployments)
		v1.POST("/deployment", deployment.CreateNewDeployment)

		/*** DeploymentOption ***/
		deploymentOption := new(controllers.DeploymentOptionController)
		v3 := v1.Group("/deploymentTemplate")
		v1.GET("/migrate/deploymentTemplate", deploymentOption.MigrateDeploymentOption)
		v3.GET("", deploymentOption.ListDeploymentOptions)
		v3.POST("", deploymentOption.CreateDeploymentOption)
		v3.GET("/:id", deploymentOption.ReadDeploymentOptionByID)
		v3.PUT("/:id", deploymentOption.UpdateDeploymentOptionByID)
		v3.DELETE("/:id", deploymentOption.DeleteDeploymentOptionByID)
		/*** Services ***/
		/*** Presistent Volume ***/
		/*** Presistent Volume Clain ***/
		/*** ConfigMap ***/
		/*** Hook ***/
		hooks := new(controllers.HookController)
		v2 := v1.Group("/hook")
		v1.GET("/migrate/hook", hooks.MigrateHook)
		v2.GET("", hooks.ListHooks)
		v2.POST("", hooks.CreateHook)
		v2.GET("/:id", hooks.ReadHookByID)
		v2.PUT("/:id", hooks.UpdateHookByID)
		v2.DELETE("/:id", hooks.DeleteHookByID)
		v2.POST("/:id", hooks.ExecuteHookByID)
		v4 := v1.Group("/key")
		v4.GET("/:key", hooks.ExecuteHookByKey)
		v4.POST("/:key", hooks.ExecuteHookByKey)

	}

	r.LoadHTMLGlob("./public/html/*")

	r.Static("/public", "./public")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginVersion": "v1.6.3",
			"goVersion":  runtime.Version(),
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	fmt.Println("SSL", os.Getenv("SSL"))
	port := os.Getenv("PORT")

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	if os.Getenv("SSL") == "TRUE" {

		SSLKeys := &struct {
			CERT string
			KEY  string
		}{}

		//Generated using sh generate-certificate.sh
		SSLKeys.CERT = "./cert/myCA.cer"
		SSLKeys.KEY = "./cert/myCA.key"

		r.RunTLS(":"+port, SSLKeys.CERT, SSLKeys.KEY)
	} else {
		r.Run(":" + port)
	}
}
