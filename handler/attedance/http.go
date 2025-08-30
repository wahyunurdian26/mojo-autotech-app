package handler

import (
	"net/http"

	mid "mojo-autotech/middleware"

	"github.com/gin-gonic/gin"

	"mojo-autotech/constant"
	"mojo-autotech/model"
	attSvc "mojo-autotech/service/attedance"
)

func HttpAttendanceHandler(router *gin.Engine) {
	h := NewAttendanceHandler()
	router.POST("/attendance/check-in", mid.Auth(), h.CheckIn)
	router.POST("/attendance/check-out", mid.Auth(), h.CheckOut)
	router.GET("/attendance/today", mid.Auth(), h.Today)
}

var attendance = func() attSvc.IAttendanceService {
	return attSvc.NewAttendanceService()
}

type AttendanceHandler struct {
	attendance attSvc.IAttendanceService
}

func NewAttendanceHandler() *AttendanceHandler {
	return &AttendanceHandler{
		attendance: attendance(),
	}
}

func (h *AttendanceHandler) CheckIn(ctx *gin.Context) {
	var param attSvc.CheckInReq
	if err := ctx.ShouldBind(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, model.Response{
			Code: http.StatusBadRequest,
			Msg:  constant.ReqParamInvalid,
			Err:  err.Error(),
		})
		return
	}

	uid, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
			Err:  "user_id tidak ditemukan di context",
		})
		return
	}
	userID, ok := toUint(uid)
	if !ok || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
			Err:  "tipe/isi user_id tidak valid",
		})
		return
	}

	param.UserId = userID
	param.IP = ctx.ClientIP()

	out, created, err := h.attendance.CheckIn(ctx.Request.Context(), param)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "unauthorized":
			status = http.StatusUnauthorized
		}
		ctx.JSON(status, model.Response{
			Code: status,
			Msg:  "Gagal check-in",
			Err:  err.Error(),
		})
		return
	}

	code := http.StatusCreated
	msg := "Check-in berhasil"
	if !created {
		code = http.StatusOK
		msg = "Check-in diperbarui"
	}
	ctx.JSON(code, model.Response{
		Code: code,
		Msg:  msg,
		Data: out,
	})
}

func (h *AttendanceHandler) CheckOut(ctx *gin.Context) {
	uid, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
			Err:  "user_id tidak ditemukan di context",
		})
		return
	}
	userID, ok := toUint(uid)
	if !ok || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
			Err:  "tipe/isi user_id tidak valid",
		})
		return
	}

	out, err := h.attendance.CheckOut(ctx.Request.Context(), attSvc.CheckOutReq{
		UserId: userID,
		IP:     ctx.ClientIP(),
	})
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "unauthorized":
			status = http.StatusUnauthorized
		case "belum check-in":
			status = http.StatusConflict
		case "sudah check-out":
			status = http.StatusConflict
		case "belum check-in atau sudah check-out":
			status = http.StatusConflict
		}
		ctx.JSON(status, model.Response{
			Code: status,
			Msg:  "Gagal check-out",
			Err:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: http.StatusOK,
		Msg:  "Check-out berhasil",
		Data: out,
	})
}

func toUint(v any) (uint, bool) {
	switch x := v.(type) {
	case uint:
		return x, true
	case int:
		if x >= 0 {
			return uint(x), true
		}
	case int64:
		if x >= 0 {
			return uint(x), true
		}
	case float64:
		if x >= 0 {
			return uint(x), true
		}
	}
	return 0, false
}

func (h *AttendanceHandler) Today(ctx *gin.Context) {
	uid, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
			Err:  "user_id tidak ditemukan di context",
		})
		return
	}
	userID, ok := toUint(uid)
	if !ok || userID == 0 {
		ctx.JSON(http.StatusUnauthorized, model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
			Err:  "tipe/isi user_id tidak valid",
		})
		return
	}

	out, err := h.attendance.GetToday(ctx.Request.Context(), userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "unauthorized" {
			status = http.StatusUnauthorized
		}
		ctx.JSON(status, model.Response{
			Code: status,
			Msg:  "Gagal memuat status hari ini",
			Err:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		Code: http.StatusOK,
		Msg:  "Status hari ini",
		Data: out,
	})
}
