package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, http.StatusOK, "ok", data)
}

func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, http.StatusOK, "ok", data)
}

func Error(c *gin.Context, httpStatus, code int, msg string) {
	JSON(c, httpStatus, code, msg, EmptyData())
}

func JSON(c *gin.Context, httpStatus, code int, msg string, data interface{}) {
	if data == nil {
		data = EmptyData()
	}
	c.JSON(httpStatus, Envelope{Code: code, Msg: msg, Data: data})
}

func EmptyData() map[string]interface{} {
	return map[string]interface{}{}
}
