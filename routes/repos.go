package routes

import (
	"net/http"

	"github.com/Corwind/conman/utils"
	"github.com/gin-gonic/gin"
)

type RepoPostSchema struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

func (env *Env) V1PostRepos(c *gin.Context) {
	namespace := c.Param("id")
	admin := c.GetBool("admin_context")
	if admin {
		namespace = "public"
	}
	var form RepoPostSchema
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	repoParams := *utils.NewRepoParams(form.Name, form.URL)
	ret, err := utils.SaveRepoParams(env.DB, namespace, repoParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	params, err := utils.DecodeRepoParams(ret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, &params)
}

func (env *Env) V1GetRepos(c *gin.Context) {
	namespace := c.Param("id")
	ret, err := utils.FetchRepoParamsList(env.DB, namespace)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, internalServerError)
		return
	}

	c.JSON(http.StatusOK, &ret)
}
