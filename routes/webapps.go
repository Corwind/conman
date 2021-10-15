package routes

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/Corwind/conman/utils"
	"github.com/gin-gonic/gin"
)

type WebappPostSchema struct {
	Name       string               `form:"name" binding:"required"`
	ChartName  string               `form:"chartname" binding:"required"`
	FileHeader multipart.FileHeader `form:"values" binding:"required"`
}

func (env *Env) V1PostWebapps(c *gin.Context) {
	namespace := c.Param("id")
	if c.GetBool("admin_context") {
		namespace = "public"
	}
	var form WebappPostSchema
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	opened_file, err := form.FileHeader.Open()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	data, _ := ioutil.ReadAll(opened_file)

	webappParams := *utils.NewWebappParams(form.Name, form.ChartName, string(data))
	ret, err := utils.SaveWebappParams(env.DB, namespace, webappParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	params, err := utils.DecodeWebappParams(ret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, &params)
}

func (env *Env) V1GetWebapps(c *gin.Context) {
	namespace := c.Param("id")
	ret, err := utils.FetchWebappParamsList(env.DB, namespace)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, &ret)
}
