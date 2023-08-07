package commonqueryparams

import (
	"net/http"

	"github.com/gin-gonic/gin"

	e "github.com/antonmisa/1cctl/internal/controller/http/v1/error"
	"github.com/antonmisa/1cctl/internal/usecase/common"
	"github.com/antonmisa/1cctl/pkg/logger"
)

type commonRequestQuery struct {
	Cache      bool   `form:"cache"`
	Entrypoint string `form:"entrypoint" binding:"required"`
}

func UseCommonQueryParams(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {

		commonRequestQuery := commonRequestQuery{}

		if err := c.ShouldBind(&commonRequestQuery); err != nil {
			l.Error(err, "http - v1 - UseCommonQueryParams")
			e.ErrorResponse(c, http.StatusBadRequest, "bad request")

			return
		}

		c.Set(common.UseCache, commonRequestQuery.Cache)
		c.Set(common.Entrypoint, commonRequestQuery.Entrypoint)

		c.Next()
	}
}
