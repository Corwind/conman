package routes

import (
	"net/http"

	"github.com/Corwind/conman/utils"

	k8s_types "k8s.io/apimachinery/pkg/types"

	"github.com/gin-gonic/gin"
)

type PodInfo struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	ID        k8s_types.UID `json:"id"`
}

func (env *Env) V1GetPods(c *gin.Context) {
	namespace := c.Param("id")
	ret, err := utils.PodList(env.KubernetesClientSet, namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	values := make([]PodInfo, 0, len(ret.Items))

	for _, pod := range ret.Items {
		values = append(values, PodInfo{
			Name:      pod.ObjectMeta.Name,
			Namespace: pod.ObjectMeta.Namespace,
			ID:        pod.ObjectMeta.UID,
		})
	}

	c.JSON(http.StatusOK, &values)
}

func (env *Env) V1GetPod(c *gin.Context) {
	namespace := c.Param("id")
	name := c.Param("name")
	ret, err := utils.PodGet(env.KubernetesClientSet, namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	podInfo := PodInfo{
		Name:      ret.ObjectMeta.Name,
		Namespace: ret.ObjectMeta.Namespace,
		ID:        ret.ObjectMeta.UID,
	}

	c.JSON(http.StatusOK, &podInfo)
}
