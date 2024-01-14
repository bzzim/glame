package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bzzim/glame/pkg/token"
	"github.com/bzzim/glame/requests"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	jwt      *token.JWT
	password string
}

func NewAuthController(jwt *token.JWT, password string) AuthController {
	return AuthController{jwt: jwt, password: password}
}

func (r *AuthController) Login(ctx *gin.Context) {
	var request requests.Login

	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	if request.Password != r.password {
		responses.NewErrorResponse(ctx, http.StatusUnauthorized, responses.UnauthorizedMessage)
		return
	}

	d := request.Duration[len(request.Duration)-1:]
	i := request.Duration[:len(request.Duration)-1]

	duration, _ := strconv.Atoi(i)
	expireAt := time.Now()
	switch {
	case d == "h":
		expireAt = expireAt.Add(time.Hour * time.Duration(duration))
	case d == "d":
		expireAt = expireAt.AddDate(0, 0, duration)
	case d == "m":
		expireAt = expireAt.AddDate(0, duration, 0)
	case d == "y":
		expireAt = expireAt.AddDate(duration, 0, 0)
	}

	t, err := r.jwt.SignToken(expireAt)
	if err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	data := responses.Login{Token: t}
	responses.NewSuccessResponse(ctx, data)
}

func (r *AuthController) Validate(ctx *gin.Context) {
	var request requests.Validate

	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	isValid := r.jwt.Validate(request.Token)
	if isValid {
		data := responses.Validate{Token: responses.ValidateToken{IsValid: isValid}}
		responses.NewSuccessResponse(ctx, data)
	} else {
		responses.NewErrorResponse(ctx, http.StatusUnauthorized, responses.UnauthorizedMessage)
	}
}
