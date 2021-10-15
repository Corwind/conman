package routes

import (
	"net/http"

	"github.com/Corwind/conman/clients"
	"github.com/Corwind/conman/utils"
	"github.com/gin-gonic/gin"
)

type ReleasePostSchema struct {
	WebappName string   `form:"webappname" json:"webappname" binding:"required"`
	Hostname   string   `form:"hostname" json:"hostname" binding:"required"`
	ValuesYaml []string `form:"valuesyaml" json:"valuesyaml"`
}

type ReleaseDeleteSchema struct {
	Name string `json:"name"`
}

type ReleaseReturnSchema struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (env *Env) V1GetReleases(c *gin.Context) {
	namespace := c.Param("id")
	releases, err := utils.FetchReleases(env.DB, namespace)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &releases)
}

func (env *Env) V1PostRelease(c *gin.Context) {
	namespace := c.Param("id")
	var form ReleasePostSchema
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	client := clients.GetHelmClient(env.Kubeconfig, namespace)
	utils.InitRepos(client, env.DB, namespace)

	params, err := utils.FetchWebappParamsByName(env.DB, namespace, form.WebappName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	_, err = utils.GetNamespace(*env.KubernetesClientSet, namespace)
	if err != nil {
		utils.CreateNamespace(*env.KubernetesClientSet, namespace)
	}

	ret, err := utils.InstallWebapp(
		client,
		env.DB,
		params.(utils.WebappParams),
		namespace,
		form.Hostname,
		form.ValuesYaml,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, &ret)
}

func (env *Env) V1DeleteRelease(c *gin.Context) {
	namespace := c.Param("id")
	var form ReleaseDeleteSchema
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	client := clients.GetHelmClient(env.Kubeconfig, namespace)
	utils.InitRepos(client, env.DB, namespace)

	client.UninstallReleaseByName(form.Name)

	utils.DeleteRelease(env.DB, namespace, form.Name)
}
