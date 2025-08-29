package handler

import (
	"net/http"

	mid "mojo-autotech/middleware"

	"github.com/gin-gonic/gin"

	"mojo-autotech/model"
	attSvc "mojo-autotech/service/attedance"
)

func HttpAttendanceHandler(router *gin.Engine) {
	svc := attSvc.NewAttendanceService() // <-- langsung, tanpa repo param
	h := NewAttendanceHandler(svc)

	router.POST("/attendance/check-in", mid.Auth(), h.CheckIn)
}

type AttendanceHandler struct {
	svc attSvc.IAttendanceService
}

func NewAttendanceHandler(svc attSvc.IAttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: svc}
}

type checkInJSON struct {
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
	PhotoURL *string  `json:"photo_url"`
	Activity *string  `json:"activity"`
}

func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	var req checkInJSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Code: http.StatusBadRequest,
			Msg:  "Param request tidak valid",
			Err:  err.Error(),
		})
		return
	}

	uidVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
			Err:  "user_id tidak ditemukan di context",
		})
		return
	}

	var userID uint
	switch v := uidVal.(type) {
	case uint:
		userID = v
	case int:
		if v >= 0 {
			userID = uint(v)
		}
	case int64:
		if v >= 0 {
			userID = uint(v)
		}
	case float64:
		if v >= 0 {
			userID = uint(v)
		}
	default:
		c.JSON(http.StatusUnauthorized, model.Response{
			Code: http.StatusUnauthorized,
			Msg:  "Unauthorized",
			Err:  "tipe user_id tidak didukung",
		})
		return
	}

	out, created, err := h.svc.CheckIn(
		c.Request.Context(),
		attSvc.CheckInReq{
			UserId:   userID,
			Activity: *req.Activity,
			Lat:      req.Lat,
			Lng:      req.Lng,
			PhotoURL: req.PhotoURL,
			IP:       c.ClientIP(),
		},
	)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "sudah check-in hari ini" {
			status = http.StatusConflict
		}
		if err.Error() == "unauthorized" {
			status = http.StatusUnauthorized
		}
		c.JSON(status, model.Response{
			Code: status,
			Msg:  "Gagal check-in",
			Err:  err.Error(),
		})
		return
	}

	code := http.StatusOK
	msg := "Check-in diperbarui"
	if created {
		code = http.StatusCreated
		msg = "Check-in berhasil"
	}
	c.JSON(code, model.Response{
		Code: code,
		Msg:  msg,
		Data: out,
	})
}
