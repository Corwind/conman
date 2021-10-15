package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Corwind/conman/config"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	helmclient "github.com/mittwald/go-helm-client"
	"k8s.io/client-go/kubernetes"
)

var internalServerError = fmt.Errorf("Internal Server Error")

type Env struct {
	DB                  fdb.Database
	KubernetesClientSet *kubernetes.Clientset
	HelmClient          helmclient.Client
	Kubeconfig          []byte
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func CheckAuthentication(admin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := &http.Client{}

		endpoint := config.GOAUTH_ENDPOINTS["users"] + "/" + c.Param("id")
		req, err := http.NewRequest("GET", endpoint, nil)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, internalServerError)
			return
		}

		req.Header.Add("Authorization", c.Request.Header.Get("Authorization"))
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, internalServerError)
			return
		}

		if resp.StatusCode != 200 {
			c.AbortWithError(http.StatusInternalServerError, internalServerError)
		}

		defer resp.Body.Close()
		var user User
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			var tmp []interface{}
			json.NewDecoder(resp.Body).Decode(&tmp)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if admin && user.Id == config.ADMIN_ID {
			c.Set("admin_context", true)
		} else {
			c.Set("admin_context", false)
		}
		c.Set("user_id", user.Id)
	}
}
