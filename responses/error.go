package responses

import "github.com/gin-gonic/gin"

var (
	BadRequestMessage   = "Invalid input data"             // 400
	ForbiddenMessage    = "Access Forbidden"               // 403
	UnauthorizedMessage = "Invalid credentials"            // 401
	ServerErrorMessage  = "An internal error has occurred" // 500
	NotFoundMessage     = "Not found"                      // 404
	ConflictMessage     = "Already exists"                 // 409
)

func NewErrorResponse(c *gin.Context, statusCode int, error string) {
	// TODO: тут почему-то отправляется заголовок "text/plain; charset=utf-8" - проверить критично ли это
	c.AbortWithStatusJSON(statusCode, Base{false, error, nil})
}
