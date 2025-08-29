package attedance

import (
	"context"
	"mojo-autotech/config"

	"gorm.io/gorm"
)

type IAttendanceRepository interface {
	InsertCheckIn(ctx context.Context, a Attendance) (Attendance, error)
}

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository() IAttendanceRepository {
	return &AttendanceRepository{db: config.NewDB()}
}

const (
	qInsertCheckIn = `
INSERT INTO attendances
  (user_id, date, check_in_at, check_in_lat, check_in_lng, check_in_photo_url, check_in_ip, status, activity, created_at, updated_at)
VALUES
  (?, ?::date, NOW(), ?, ?, ?, ?, ?, ?, NOW(), NOW())
RETURNING
  id, user_id, date,
  check_in_at, check_in_lat, check_in_lng, check_in_photo_url, check_in_ip,
  check_out_at, check_out_ip, total_minutes, status, activity, notes,
  created_at, updated_at;
`
)

func (r *AttendanceRepository) InsertCheckIn(ctx context.Context, a Attendance) (Attendance, error) {
	var out Attendance
	err := r.db.WithContext(ctx).Raw(
		qInsertCheckIn,
		a.UserID,
		a.Date.Format("2006-01-02"),
		a.CheckInLat,
		a.CheckInLng,
		a.CheckInPhotoURL,
		a.CheckInIP,
		a.Status,
		a.Activity,
	).Scan(&out).Error
	return out, err
}
