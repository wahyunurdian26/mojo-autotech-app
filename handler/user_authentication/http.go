package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mojo-autotech/constant"
	"mojo-autotech/model"
	ua "mojo-autotech/model/user_authentication"
	authsvc "mojo-autotech/service/user_authentication"
)

func HttpHandler(router *gin.Engine) {
	handler := NewHandler()
	{
		router.POST("/login", handler.Login)
		router.POST("/create", handler.CreateAccount)
	}
}

var authentication = func() authsvc.IAuthService {

	return authsvc.NewAuthService()
}

type Handler struct {
	authentication authsvc.IAuthService
}

func NewHandler() *Handler {
	return &Handler{
		authentication: authentication(),
	}
}

func (h *Handler) Login(ctx *gin.Context) {
	var param ua.LoginReq
	if err := ctx.ShouldBind(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Response{
			Code: http.StatusBadRequest,
			Msg:  constant.ReqParamInvalid,
			Err:  err.Error(),
		})
		return
	}
	res, err := h.authentication.Login(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Response{
			Code: http.StatusInternalServerError,
			Msg:  constant.LoginError,
			Err:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, model.Response{
		Code: http.StatusCreated,
		Msg:  constant.LoginSuccess,
		Data: res,
	})
}

func (h *Handler) CreateAccount(ctx *gin.Context) {
	var param ua.RegisterReq
	if err := ctx.ShouldBind(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Response{
			Code: http.StatusBadRequest,
			Msg:  constant.ReqParamInvalid,
			Err:  err.Error(),
		})
		return
	}
	res, err := h.authentication.CreateAccount(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, model.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Gagal membuat akun",
			Err:  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, model.Response{
		Code: http.StatusCreated,
		Msg:  "Akun berhasil dibuat",
		Data: res,
	})
}
