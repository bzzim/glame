package setting

import (
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/pkg/utils"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const queriesFile = "data/customQueries.json"

func (r *Controller) Queries(ctx *gin.Context) {
	queries, err := r.loadQueries()
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	responses.NewSuccessResponse(ctx, queries.Queries)
}

func (r *Controller) AddQuery(ctx *gin.Context) {
	var request models.Query

	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	queries, err := r.loadQueries()
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	for _, v := range queries.Queries {
		if v.Prefix == request.Prefix {
			responses.NewErrorResponse(ctx, http.StatusBadRequest, "Prefix must be unique")
			return
		}
	}

	queries.Queries = append(queries.Queries, request)
	err = r.saveQueries(queries)
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	responses.NewSuccessResponse(ctx, request)
}

func (r *Controller) SaveQuery(ctx *gin.Context) {
	prefix := strings.ToLower(ctx.Param("prefix"))
	if len(prefix) == 0 {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	var request models.Query

	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	queries, err := r.loadQueries()
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	for i, v := range queries.Queries {
		if v.Prefix == prefix {
			queries.Queries[i] = request
			break
		}
	}

	err = r.saveQueries(queries)
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	responses.NewSuccessResponse(ctx, queries.Queries)
}

func (r *Controller) DeleteQuery(ctx *gin.Context) {
	prefix := strings.ToLower(ctx.Param("prefix"))
	if len(prefix) == 0 {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	queries, err := r.loadQueries()
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	var newQueries = make([]models.Query, 0)
	for _, t := range queries.Queries {
		if prefix != strings.ToLower(t.Prefix) {
			newQueries = append(newQueries, t)
		}
	}

	queries.Queries = newQueries

	err = r.saveQueries(queries)
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	responses.NewSuccessResponse(ctx, newQueries)
}

func (r *Controller) loadQueries() (*models.Queries, error) {
	data, err := utils.ReadJson[models.Queries](queriesFile)
	if err != nil {
		r.log.Warn(err.Error())
	}
	return data, err
}

func (r *Controller) saveQueries(queries *models.Queries) error {
	err := utils.SaveJson(queriesFile, queries)
	if err != nil {
		r.log.Warn(err.Error())
	}
	return err
}
