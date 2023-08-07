package clustercredentials

import (
	"net/http"

	"github.com/gin-gonic/gin"

	e "github.com/antonmisa/1cctl/internal/controller/http/v1/error"
	"github.com/antonmisa/1cctl/internal/entity"
	"github.com/antonmisa/1cctl/internal/usecase/common"
	"github.com/antonmisa/1cctl/pkg/logger"
)

type clusterCred struct {
	Login    string `header:"login"`
	Password string `header:"password"`
}

func UseClusterCredentials(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {

		clusterCred := clusterCred{}

		if err := c.ShouldBindHeader(&clusterCred); err != nil {
			l.Error(err, "http - v1 - UseClusterCredentials")
			e.ErrorResponse(c, http.StatusBadRequest, "bad request")

			return
		}

		c.Set(common.ClusterCred, entity.Credentials{Name: clusterCred.Login, Pwd: clusterCred.Password})

		c.Next()
	}
}
