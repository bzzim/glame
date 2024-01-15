package setting

import (
	"net/http"

	"github.com/bzzim/glame/pkg/utils"
	"github.com/bzzim/glame/requests"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
)

const CssFile = "data/flame.css"

func (r *Controller) Css(ctx *gin.Context) {
	css, err := utils.ReadFile(CssFile)
	if err != nil {
		r.log.Warn(err.Error())
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	responses.NewSuccessResponse(ctx, css)
}

func (r *Controller) SaveCss(ctx *gin.Context) {
	var request requests.Style

	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	if err := utils.WriteToFile(request.Css, CssFile); err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	responses.NewSuccessResponse(ctx, "")
}
