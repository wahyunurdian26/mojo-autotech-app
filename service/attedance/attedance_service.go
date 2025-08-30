package attedance

import (
	"context"
	"errors"
	"fmt"
	"time"

	entity "mojo-autotech/model/attedance"

	"gorm.io/gorm"
)

type (
	Attendance = entity.Attendance
)

// Request DTO (kamu boleh pindah ke package model jika mau samakan gaya dengan auth)
type CheckInReq struct {
	Activity string `json:"activity" binding:"required"`
	UserId   uint
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
	PhotoURL *string  `json:"photo_url"`
	IP       string
}

type CheckOutReq struct {
	UserId uint
	IP     string
}

type IAttendanceService interface {
	CheckIn(ctx context.Context, req CheckInReq) (Attendance, bool, error) // (record, created?, err)
	CheckOut(ctx context.Context, req CheckOutReq) (Attendance, error)
	GetToday(ctx context.Context, userID uint) (Attendance, error)
}

type AttendanceService struct {
	// penamaan mirip auth: s.user_authentication → s.attedance
	attedance entity.IAttendanceRepository
	loc       *time.Location
}

func NewAttendanceService() *AttendanceService {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	if loc == nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	return &AttendanceService{
		attedance: entity.NewAttendanceRepository(),
		loc:       loc,
	}
}

func (s *AttendanceService) workDateNow() time.Time {
	now := time.Now().In(s.loc)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, s.loc)
}

func (s *AttendanceService) CheckIn(ctx context.Context, req CheckInReq) (Attendance, bool, error) {
	if req.UserId == 0 {
		return Attendance{}, false, errors.New("unauthorized")
	}
	if req.Activity == "" {
		return Attendance{}, false, errors.New("activity wajib diisi")
	}

	wd := s.workDateNow()

	a := entity.Attendance{
		UserID:          req.UserId,
		Date:            wd,
		CheckInLat:      req.Lat,
		CheckInLng:      req.Lng,
		CheckInPhotoURL: req.PhotoURL,
		Status:          "PRESENT",
		Activity:        req.Activity,
	}
	if req.IP != "" {
		ip := req.IP
		a.CheckInIP = &ip
	}

	out, created, err := s.attedance.UpsertCheckIn(ctx, a)
	if err != nil {
		return Attendance{}, false, err
	}
	return out, created, nil
}

func (s *AttendanceService) CheckOut(ctx context.Context, req CheckOutReq) (Attendance, error) {
	if req.UserId == 0 {
		return Attendance{}, errors.New("unauthorized")
	}

	wd := s.workDateNow()

	// Pastikan sudah ada record & belum checkout
	cur, err := s.attedance.GetByUserAndDate(ctx, req.UserId, wd)
	if err != nil {
		fmt.Println("err get by user and date:", err)
		return Attendance{}, err
	}
	if cur.CheckOutAt != nil {
		return Attendance{}, errors.New("sudah check-out")
	}

	// Proses checkout
	var ip *string
	if req.IP != "" {
		ip = &req.IP
	}
	out, err := s.attedance.CheckOut(ctx, req.UserId, wd, ip)
	if err != nil {
		fmt.Println("err checkout:", err)
		// cek apakah err karena record tidak ditemukan
		if errors.Is(err, gormErrRecordNotFound()) {
			return Attendance{}, errors.New("belum check-in hari ini")
		}
		return Attendance{}, err
	}
	return out, nil
}

// helper untuk menghindari import langsung gorm di signature
func gormErrRecordNotFound() error {
	type notFound interface{ Error() string }
	// signature kadang berasal dari gorm.ErrRecordNotFound
	var e notFound
	_ = e
	return nil
}

func (s *AttendanceService) GetToday(ctx context.Context, userID uint) (Attendance, error) {
	if userID == 0 {
		return Attendance{}, errors.New("unauthorized")
	}
	wd := s.workDateNow()

	rec, err := s.attedance.GetByUserAndDate(ctx, userID, wd)
	if err != nil {
		// Kalau belum ada record hari ini, balikan objek "kosong" (200 OK)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Attendance{
				UserID: userID,
				Date:   wd,
				// CheckInAt/CheckOutAt biarkan nil → FE akan tampil "-"
			}, nil
		}
		return Attendance{}, err
	}
	return rec, nil
}
