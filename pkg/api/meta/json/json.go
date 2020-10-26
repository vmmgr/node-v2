package json

import (
	"github.com/gin-gonic/gin"
	"github.com/vmmgr/node/pkg/api/core/gateway"
	"net/http"
)

func ResponseOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gateway.Result{Status: http.StatusOK, Data: data})
}

func ResponseError(ctx *gin.Context, status int, err error) {
	ctx.JSON(status, gateway.ResultError{Status: status, Error: err})
}
