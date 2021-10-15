package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Corwind/conman/clients"
	"github.com/Corwind/conman/routes"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

func WebappRetrieve(router *gin.RouterGroup, env *routes.Env) {
	router.GET("/:id/webapps", env.V1GetWebapps)
}

func WebappsRegistration(router *gin.RouterGroup, env *routes.Env) {
	router.POST("/:id/webapps", env.V1PostWebapps)
}

func ReposRegistration(router *gin.RouterGroup, env *routes.Env) {
	router.POST("/:id/repos", env.V1PostRepos)
}

func ReposRetrieve(router *gin.RouterGroup, env *routes.Env) {
	router.GET("/:id/repos", env.V1GetRepos)
}

func InstancesGroup(router *gin.RouterGroup, env *routes.Env) {
	router.GET("/:id/releases", env.V1GetReleases)
	router.POST("/:id/releases", env.V1PostRelease)
	router.DELETE("/:id/releases", env.V1DeleteRelease)
}

func PodsGroup(router *gin.RouterGroup, env *routes.Env) {
	router.GET("/:id/pods", env.V1GetPods)
	router.GET("/:id/pods/:name", env.V1GetPod)
}

func SetupCors(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3002"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}

func main() {
	kubeconfig, err := os.ReadFile("./kubeconfig")
	fdb.MustAPIVersion(630)
	if err != nil {
		panic(err.Error())
	}

	env := routes.Env{
		KubernetesClientSet: clients.K8sClient(kubeconfig),
		DB:                  fdb.MustOpenDefault(),
		Kubeconfig:          kubeconfig,
	}

	router := gin.Default()

	SetupCors(router)

	v1 := router.Group("/api/v1")
	v1.Use(routes.CheckAuthentication(false))
	ReposRetrieve(v1, &env)
	WebappRetrieve(v1, &env)
	InstancesGroup(v1, &env)
	PodsGroup(v1, &env)
	v1.Use(routes.CheckAuthentication(true))
	WebappsRegistration(v1, &env)
	ReposRegistration(v1, &env)

	router.Run(":3001")
}
