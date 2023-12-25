package setting

import (
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/pkg/utils"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const themesFile = "data/themes.json"

func (r *Controller) Themes(ctx *gin.Context) {
	themes, err := r.loadThemes()
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	responses.NewSuccessResponse(ctx, themes.Themes)
}

func (r *Controller) AddTheme(ctx *gin.Context) {
	var request models.Theme

	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	themes, err := r.loadThemes()
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	for _, t := range themes.Themes {
		if strings.ToLower(t.Name) == strings.ToLower(request.Name) {
			responses.NewErrorResponse(ctx, http.StatusBadRequest, "Name must be unique")
			return
		}
	}

	themes.Themes = append(themes.Themes, request)

	err = r.saveThemes(themes)
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	responses.NewSuccessResponse(ctx, request)
}

func (r *Controller) SaveTheme(ctx *gin.Context) {
	name := strings.ToLower(ctx.Param("name"))
	if len(name) == 0 {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	var request models.Theme
	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	themes, err := r.loadThemes()
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	var customThemes = make([]models.Theme, 0) // слайс, который нужно вернуть
	for i, t := range themes.Themes {
		if !t.IsCustom {
			continue
		}
		if strings.ToLower(t.Name) == name {
			themes.Themes[i] = request
		}
		customThemes = append(customThemes, themes.Themes[i])
	}

	err = r.saveThemes(themes)
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	responses.NewSuccessResponse(ctx, customThemes)
}

func (r *Controller) DeleteTheme(ctx *gin.Context) {
	name := strings.ToLower(ctx.Param("name"))
	if len(name) == 0 {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	themes, err := r.loadThemes()
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	var newThemes = make([]models.Theme, 0)    // слайс, который нужно в файл записать
	var customThemes = make([]models.Theme, 0) // слайс, который нужно вернуть
	for _, t := range themes.Themes {
		if !t.IsCustom {
			newThemes = append(newThemes, t)
		} else {
			if name != strings.ToLower(t.Name) {
				newThemes = append(newThemes, t)
				customThemes = append(customThemes, t)
			}
		}
	}

	themes.Themes = newThemes

	err = r.saveThemes(themes)
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	responses.NewSuccessResponse(ctx, customThemes)
}

func (r *Controller) loadThemes() (*models.Themes, error) {
	data, err := utils.ReadJson[models.Themes](themesFile)
	if err != nil {
		r.log.Warn(err.Error())
	}
	return data, err
}

func (r *Controller) saveThemes(themes *models.Themes) error {
	err := utils.SaveJson(themesFile, themes)
	if err != nil {
		r.log.Warn(err.Error())
	}
	return err
}
