package attedance

import (
	"context"
	"time"

	model "mojo-autotech/model/attedance"
)

type CheckInReq struct {
	Activity string `json:"activity" binding:"required"`
	UserId   uint
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
	PhotoURL *string  `json:"photo_url"`
	IP       string
}

type IAttendanceService interface {
	CheckIn(ctx context.Context, req CheckInReq) (model.Attendance, bool, error)
}

type AttendanceService struct {
	repo model.IAttendanceRepository
	loc  *time.Location
}

func NewAttendanceService() *AttendanceService {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	if loc == nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	return &AttendanceService{
		repo: model.NewAttendanceRepository(),
		loc:  loc,
	}
}

// workDate in Asia/Jakarta (tanggal tanpa jam)
func (s *AttendanceService) workDateNow() (time.Time, string) {
	now := time.Now().In(s.loc)
	wd := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, s.loc)
	return wd, wd.Format("2006-01-02")
}

func (s *AttendanceService) CheckIn(ctx context.Context, req CheckInReq) (model.Attendance, bool, error) {
	wd, _ := s.workDateNow()

	a := model.Attendance{
		UserID:          req.UserId,
		Activity:        req.Activity,
		Date:            wd,
		CheckInLat:      req.Lat,
		CheckInLng:      req.Lng,
		CheckInPhotoURL: req.PhotoURL,
		Status:          "PRESENT",
	}
	ip := req.IP
	a.CheckInIP = &ip

	out, err := s.repo.InsertCheckIn(ctx, a)
	if err != nil {
		return model.Attendance{}, false, err
	}
	return out, true, nil
}
