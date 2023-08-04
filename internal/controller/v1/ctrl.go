package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/antonmisa/1cctl/internal/entity"
	"github.com/antonmisa/1cctl/internal/usecase"
	"github.com/antonmisa/1cctl/pkg/logger"
)

type ctrl1CRoutes struct {
	t usecase.Ctrl
	l logger.Interface
}

func newCtrl1CRoutes(handler *gin.RouterGroup, t usecase.Ctrl, l logger.Interface) {
	r := &ctrl1CRoutes{t, l}

	h := handler.Group("/cluster")
	{
		h.GET("/list", r.clusters)
		h.GET("/:cluster/infobase/list", r.infobases)
		h.GET("/:cluster/infobase/:infobase/session/list", r.sessionsByInfobase)
		h.GET("/:cluster/session/list", r.sessions)
	}
}

type clusterResponse struct {
	Clusters []entity.Cluster `json:"clusters"`
}

// @Summary     Show clusters
// @Description Show all clusters with identifiers
// @ID          clusters
// @Tags  	    cluster list
// @Accept      nil
// @Produce     json
// @Success     200 {object} clusterResponse
// @Failure     500 {object} response
// @Router      /cluster/list [get]
func (r *ctrl1CRoutes) clusters(c *gin.Context) {
	clusters, err := r.t.Clusters(c.Request.Context())
	if err != nil {
		r.l.Error(err, "http - v1 - clusters")
		errorResponse(c, http.StatusInternalServerError, "internal problems")

		return
	}

	c.JSON(http.StatusOK, clusterResponse{clusters})
}

type infobaseRequest struct {
	Cluster string `json:"cluster"       binding:"required"  example:"UUID"`
}

type infobaseResponse struct {
	Infobases []entity.Infobase `json:"infobases"`
}

// @Summary     Show all infobases in cluster
// @Description Show all infobases with identifiers for current cluster
// @ID          infobases
// @Tags  	    infobase list
// @Accept      uri
// @Produce     json
// @Success     200 {object} infobaseResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /cluster/:cluster/infobase/list [get]
func (r *ctrl1CRoutes) infobases(c *gin.Context) {
	var request infobaseRequest
	if err := c.ShouldBindUri(&request); err != nil {
		r.l.Error(err, "http - v1 - infobases")
		errorResponse(c, http.StatusBadRequest, "invalid request parameters")

		return
	}

	infobases, err := r.t.Infobases(c.Request.Context(), entity.Cluster{ID: request.Cluster})
	if err != nil {
		r.l.Error(err, "http - v1 - infobases")
		errorResponse(c, http.StatusInternalServerError, "internal problems")

		return
	}

	c.JSON(http.StatusOK, infobaseResponse{infobases})
}

type sessionRequest struct {
	Cluster  string `json:"cluster"       binding:"required"  example:"UUID"`
	Infobase string `json:"infobase"      example:"UUID"`
}

type sessionResponse struct {
	Sessions []entity.Session `json:"sessions"`
}

// @Summary     Show all sessions in cluster
// @Description Show all sessions with identifiers for current cluster
// @ID          sessions
// @Tags  	    session list
// @Accept      uri
// @Produce     json
// @Success     200 {object} sessionResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /cluster/:cluster/session/list [get]
func (r *ctrl1CRoutes) sessions(c *gin.Context) {
	var request sessionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		r.l.Error(err, "http - v1 - sessions")
		errorResponse(c, http.StatusBadRequest, "invalid request parameters")

		return
	}

	sessions, err := r.t.Sessions(c.Request.Context(), entity.Cluster{ID: request.Cluster}, entity.Infobase{ID: request.Infobase})
	if err != nil {
		r.l.Error(err, "http - v1 - sessions")
		errorResponse(c, http.StatusInternalServerError, "internal problems")

		return
	}

	c.JSON(http.StatusOK, sessionResponse{sessions})
}

// @Summary     Show all sessions in infobase
// @Description Show all sessions with identifiers for current infobase in cluster
// @ID          sessionsByInfobase
// @Tags  	    session list infobase
// @Accept      uri
// @Produce     json
// @Success     200 {object} sessionResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /cluster/:cluster/infobase/:infobase/session/list [get]
func (r *ctrl1CRoutes) sessionsByInfobase(c *gin.Context) {
	var request sessionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		r.l.Error(err, "http - v1 - sessionsByInfobase")
		errorResponse(c, http.StatusBadRequest, "invalid request parameters")

		return
	}

	sessions, err := r.t.Sessions(c.Request.Context(), entity.Cluster{ID: request.Cluster}, entity.Infobase{ID: request.Infobase})
	if err != nil {
		r.l.Error(err, "http - v1 - sessionsByInfobase")
		errorResponse(c, http.StatusInternalServerError, "internal problems")

		return
	}

	c.JSON(http.StatusOK, sessionResponse{sessions})
}
