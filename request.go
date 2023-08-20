package framework

import (
	"github.com/codeunic/golang-framework/functions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Request struct {
	Database   *Db
	Engine     *gin.Engine
	Context    *gin.Context
	Pagination *Pagination
	Services   any
}

func (r *Request) Ok(body any, message ...interface{}) {
	r.Context.JSON(http.StatusOK, NewResponse(http.StatusOK, body, true, message...))
}

func (r *Request) Error(body any, message ...interface{}) {
	r.Context.JSON(http.StatusInternalServerError, NewResponse(http.StatusInternalServerError, body, false, message...))
}

func (r *Request) BadRequest(body any, message ...interface{}) {
	r.Context.JSON(http.StatusBadRequest, NewResponse(http.StatusBadRequest, body, false, message...))
}

func (r *Request) NotFound(body any, message ...interface{}) {
	r.Context.JSON(http.StatusNotFound, NewResponse(http.StatusNotFound, body, false, message...))
}

func (r *Request) GetParam(param string) string {
	return r.Context.Params.ByName(param)
}

func (r *Request) GetParamNumber(s string) int64 {
	value := r.Context.Params.ByName(s)
	i, _ := strconv.ParseInt(value, 10, 64)

	return i
}

func NewRequest(database *Db, engine *gin.Engine, services any, context *gin.Context) *Request {
	limitPagination, _ := functions.StringToInt64(context.DefaultQuery("limit", "20"))
	pagePagination, _ := functions.StringToInt64(context.DefaultQuery("page", "1"))

	return &Request{
		Database:   database,
		Engine:     engine,
		Context:    context,
		Services:   services,
		Pagination: NewPagination(pagePagination, limitPagination),
	}
}
