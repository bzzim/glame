package setting

import (
	"github.com/bzzim/glame/helper"
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/pkg/utils"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

const ConfigFile = "data/config.json"

func (r *Controller) Config(ctx *gin.Context) {
	config, err := helper.LoadConfig(ConfigFile)
	if err != nil {
		r.log.Warn(err.Error())
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	responses.NewSuccessResponse(ctx, config)
}

func (r *Controller) SaveConfig(ctx *gin.Context) {
	var request models.Config

	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	err := utils.SaveJson(ConfigFile, &request)
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		r.log.Warn(err.Error())
	}
	responses.NewSuccessResponse(ctx, request)
}
