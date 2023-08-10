package v1

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"github.com/gin-gonic/gin"

	v1e "github.com/antonmisa/1cctl/internal/controller/http/v1/error"
	"github.com/antonmisa/1cctl/internal/controller/http/v1/middleware/clustercredentials"
	"github.com/antonmisa/1cctl/internal/entity"
	"github.com/antonmisa/1cctl/internal/usecase"
	"github.com/antonmisa/1cctl/internal/usecase/common"
	"github.com/antonmisa/1cctl/pkg/logger"
)

type ctrlRoutes struct {
	c usecase.Ctrl
	l logger.Interface
	t trace.Tracer
}

func newCtrlRoutes(handler *gin.RouterGroup, t usecase.Ctrl, l logger.Interface, tr trace.Tracer) {
	r := &ctrlRoutes{t, l, tr}

	h := handler.Group("/cluster")
	{
		h.Use(clustercredentials.UseClusterCredentials(l))

		h.GET("/list", r.clusters)
		h.GET("/:cluster/infobase/list", r.infobases)
		h.GET("/:cluster/infobase/:infobase/session/list", r.sessionsByInfobase)
		h.GET("/:cluster/infobase/:infobase/connection/list", r.connectionsByInfobase)
		h.GET("/:cluster/session/list", r.sessions)
		h.GET("/:cluster/connection/list", r.connections)
	}
}

type clusterResponse struct {
	Clusters []entity.Cluster `json:"clusters"`
}

// @Summary     Show clusters
// @Description Show all clusters with data
// @ID          clusters
// @Tags  	    cluster list
// @Produce     json
// @Param		cache	    query	 bool			false	"Firstly try to find from Cache"
// @Param       entrypoint  query    string         true 	"Entrypoint for cluster"
// @Success     200 {object} clusterResponse
// @Failure     400 {object} error.response
// @Failure     404 {object} error.response
// @Failure     500 {object} error.response
// @Router      /cluster/list [get]
func (r *ctrlRoutes) clusters(c *gin.Context) {
	ctx, span := r.t.Start(c.Request.Context(), "clusters")
	defer span.End()

	entrypoint := c.GetString(common.Entrypoint)

	args := map[string]any{
		common.UseCache: c.MustGet(common.UseCache),
	}

	span.AddEvent("get list of clusters")

	clusters, err := r.c.Clusters(ctx, entrypoint, args)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - clusters - r.c.Clusters")
		v1e.ErrorResponse(c, http.StatusInternalServerError, "internal problems")

		return
	}

	span.AddEvent("generate json response")

	c.JSON(http.StatusOK, clusterResponse{clusters})
}

type infobaseRequest struct {
	Cluster string `uri:"cluster"       binding:"required"  example:"UUID"`
}

type infobaseResponse struct {
	Infobases []entity.Infobase `json:"infobases"`
}

// @Summary     Show all infobases in cluster
// @Description Show all infobases with identifiers for current cluster
// @ID          infobases
// @Tags  	    infobase list
// @Produce     json
// @Param		cluster	    path	 string			true	"UUID of cluster"
// @Param		cache	    query	 bool			false	"Firstly try to find from Cache"
// @Param       entrypoint  query    string         true 	"Entrypoint for cluster"
// @Success     200 {object} infobaseResponse
// @Failure     400 {object} error.response
// @Failure     404 {object} error.response
// @Failure     500 {object} error.response
// @Router      /cluster/:cluster/infobase/list [get]
func (r *ctrlRoutes) infobases(c *gin.Context) {
	ctx, span := r.t.Start(c.Request.Context(), "infobases")
	defer span.End()

	var request infobaseRequest

	span.AddEvent("binding incoming params")

	if err := c.ShouldBindUri(&request); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - infobases")
		v1e.ErrorResponse(c, http.StatusBadRequest, "invalid request parameters")

		return
	}

	entrypoint := c.GetString(common.Entrypoint)

	args := map[string]any{
		common.UseCache:    c.MustGet(common.UseCache),
		common.ClusterCred: c.MustGet(common.ClusterCred),
	}

	clusterCred, _ := c.MustGet(common.ClusterCred).(entity.Credentials)

	span.AddEvent("get list of infobases")

	infobases, err := r.c.Infobases(ctx, entrypoint, entity.Cluster{ID: request.Cluster}, clusterCred, args)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - infobases")
		v1e.ErrorResponse(c, http.StatusInternalServerError, "internal problems")

		return
	}

	span.AddEvent("generate json response")

	c.JSON(http.StatusOK, infobaseResponse{infobases})
}

type requestWoInfobase struct {
	Cluster string `uri:"cluster"       binding:"required"  example:"UUID"`
}

type sessionResponse struct {
	Sessions []entity.Session `json:"sessions"`
}

// @Summary     Show all sessions in cluster
// @Description Show all sessions with identifiers for current cluster
// @ID          sessions
// @Tags  	    session list
// @Produce     json
// @Param		cluster	    path	 string			true	"UUID of cluster"
// @Param		cache	    query	 bool			false	"Firstly try to find from Cache"
// @Param       entrypoint  query    string         true 	"Entrypoint for cluster"
// @Success     200 {object} sessionResponse
// @Failure     400 {object} error.response
// @Failure     404 {object} error.response
// @Failure     500 {object} error.response
// @Router      /cluster/:cluster/session/list [get]
func (r *ctrlRoutes) sessions(c *gin.Context) {
	ctx, span := r.t.Start(c.Request.Context(), "sessions")
	defer span.End()

	var request requestWoInfobase

	span.AddEvent("binding incoming params")

	if err := c.ShouldBindUri(&request); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - sessions")
		v1e.ErrorResponse(c, http.StatusBadRequest, "invalid request parameters")

		return
	}

	entrypoint := c.GetString(common.Entrypoint)

	args := map[string]any{
		common.UseCache:    c.MustGet(common.UseCache),
		common.ClusterCred: c.MustGet(common.ClusterCred),
	}

	clusterCred, _ := c.MustGet(common.ClusterCred).(entity.Credentials)

	span.AddEvent("get list of sessions")

	sessions, err := r.c.Sessions(ctx, entrypoint, entity.Cluster{ID: request.Cluster}, clusterCred, entity.Infobase{}, args)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - sessions")
		v1e.ErrorResponse(c, http.StatusInternalServerError, "internal problems")

		return
	}

	span.AddEvent("generate json response")

	c.JSON(http.StatusOK, sessionResponse{sessions})
}

type requestWInfobase struct {
	Cluster  string `uri:"cluster"       binding:"required"  example:"UUID"`
	Infobase string `uri:"infobase"      binding:"required"  example:"UUID"`
}

// @Summary     Show all sessions in infobase
// @Description Show all sessions with identifiers for current infobase in cluster
// @ID          sessionsByInfobase
// @Tags  	    session list infobase
// @Produce     json
// @Param		cluster	    path	 string			true	"UUID of cluster"
// @Param		infobase    path	 string			true	"UUID of infobase"
// @Param		cache	    query	 bool			false	"Firstly try to find from Cache"
// @Param       entrypoint  query    string         true 	"Entrypoint for cluster"
// @Success     200 {object} sessionResponse
// @Failure     400 {object} error.response
// @Failure     404 {object} error.response
// @Failure     500 {object} error.response
// @Router      /cluster/:cluster/infobase/:infobase/session/list [get]
func (r *ctrlRoutes) sessionsByInfobase(c *gin.Context) {
	ctx, span := r.t.Start(c.Request.Context(), "sessions")
	defer span.End()

	var request requestWInfobase

	span.AddEvent("binding incoming params")

	if err := c.ShouldBindUri(&request); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - sessionsByInfobase")
		v1e.ErrorResponse(c, http.StatusBadRequest, "invalid request parameters")

		return
	}

	entrypoint := c.GetString(common.Entrypoint)

	args := map[string]any{
		common.UseCache:    c.MustGet(common.UseCache),
		common.ClusterCred: c.MustGet(common.ClusterCred),
	}

	clusterCred, _ := c.MustGet(common.ClusterCred).(entity.Credentials)

	span.AddEvent("get list of sessions")

	sessions, err := r.c.Sessions(ctx, entrypoint, entity.Cluster{ID: request.Cluster}, clusterCred, entity.Infobase{ID: request.Infobase}, args)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - sessionsByInfobase")
		v1e.ErrorResponse(c, http.StatusInternalServerError, "internal problems")

		return
	}

	span.AddEvent("generate json response")

	c.JSON(http.StatusOK, sessionResponse{sessions})
}

type connectionResponse struct {
	Connections []entity.Connection `json:"connections"`
}

// @Summary     Show all connections in cluster
// @Description Show all connections with identifiers for current cluster
// @ID          connections
// @Tags  	    connection list
// @Produce     json
// @Param		cluster	    path	 string			true	"UUID of cluster"
// @Param		cache	    query	 bool			false	"Firstly try to find from Cache"
// @Param       entrypoint  query    string         true 	"Entrypoint for cluster"
// @Success     200 {object} connectionResponse
// @Failure     400 {object} error.response
// @Failure     404 {object} error.response
// @Failure     500 {object} error.response
// @Router      /cluster/:cluster/connection/list [get]
func (r *ctrlRoutes) connections(c *gin.Context) {
	ctx, span := r.t.Start(c.Request.Context(), "sessions")
	defer span.End()

	var request requestWoInfobase

	span.AddEvent("binding incoming params")

	if err := c.ShouldBindUri(&request); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - connections")
		v1e.ErrorResponse(c, http.StatusBadRequest, "invalid request parameters")

		return
	}

	entrypoint := c.GetString(common.Entrypoint)

	args := map[string]any{
		common.UseCache:    c.MustGet(common.UseCache),
		common.ClusterCred: c.MustGet(common.ClusterCred),
	}

	clusterCred, _ := c.MustGet(common.ClusterCred).(entity.Credentials)

	span.AddEvent("get list of connections")

	connections, err := r.c.Connections(ctx, entrypoint, entity.Cluster{ID: request.Cluster}, clusterCred, entity.Infobase{}, args)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - connections")
		v1e.ErrorResponse(c, http.StatusInternalServerError, "internal problems")

		return
	}

	span.AddEvent("generate json response")

	c.JSON(http.StatusOK, connectionResponse{connections})
}

// @Summary     Show all connections in infobase
// @Description Show all connections with identifiers for current infobase in cluster
// @ID          connectionsByInfobase
// @Tags  	    connection list infobase
// @Produce     json
// @Param		cluster	    path	 string			true	"UUID of cluster"
// @Param		infobase    path	 string			true	"UUID of infobase"
// @Param		cache	    query	 bool			false	"Firstly try to find from Cache"
// @Param       entrypoint  query    string         true 	"Entrypoint for cluster"
// @Success     200 {object} connectionResponse
// @Failure     400 {object} error.response
// @Failure     404 {object} error.response
// @Failure     500 {object} error.response
// @Router      /cluster/:cluster/infobase/:infobase/connection/list [get]
func (r *ctrlRoutes) connectionsByInfobase(c *gin.Context) {
	ctx, span := r.t.Start(c.Request.Context(), "sessions")
	defer span.End()

	var request requestWInfobase

	span.AddEvent("binding incoming params")

	if err := c.ShouldBindUri(&request); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - sessionsByInfobase")
		v1e.ErrorResponse(c, http.StatusBadRequest, "invalid request parameters")

		return
	}

	entrypoint := c.GetString(common.Entrypoint)

	args := map[string]any{
		common.UseCache:    c.MustGet(common.UseCache),
		common.ClusterCred: c.MustGet(common.ClusterCred),
	}

	clusterCred, _ := c.MustGet(common.ClusterCred).(entity.Credentials)

	span.AddEvent("get list of connections")

	connections, err := r.c.Connections(ctx, entrypoint, entity.Cluster{ID: request.Cluster}, clusterCred, entity.Infobase{ID: request.Infobase}, args)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		r.l.Error(err, "http - v1 - sessionsByInfobase")
		v1e.ErrorResponse(c, http.StatusInternalServerError, "internal problems")

		return
	}

	c.JSON(http.StatusOK, connectionResponse{connections})
}
