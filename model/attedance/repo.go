package attedance

import (
	"context"
	"mojo-autotech/config"
	"time"

	"gorm.io/gorm"
)

type IAttendanceRepository interface {
	UpsertCheckIn(ctx context.Context, a Attendance) (out Attendance, created bool, err error)
	GetByUserAndDate(ctx context.Context, userID uint, date time.Time) (Attendance, error)
	CheckOut(ctx context.Context, userID uint, date time.Time, ip *string) (Attendance, error)
}

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository() IAttendanceRepository {
	return &AttendanceRepository{db: config.NewDB()}
}

const (
	// UPSERT: kalau (user_id, date) sudah ada -> update nilai check_in_* dan activity; tidak error 23505 lagi
	qUpsertCheckIn = `
WITH ins AS (
  INSERT INTO attendances
    (user_id, date, check_in_at, check_in_lat, check_in_lng, check_in_photo_url, check_in_ip, status, activity, created_at, updated_at)
  VALUES
    (?, ?::date, NOW(), ?, ?, ?, ?, ?, ?, NOW(), NOW())
  ON CONFLICT (user_id, date) DO NOTHING
  RETURNING
    id, user_id, date,
    check_in_at, check_in_lat, check_in_lng, check_in_photo_url, check_in_ip,
    check_out_at, check_out_ip, total_minutes, status, activity, notes,
    created_at, updated_at,
    TRUE AS created
), upd AS (
  UPDATE attendances a
  SET
    check_in_at        = COALESCE(a.check_in_at, NOW()), -- tidak overwrite kalau sudah ada
    check_in_lat       = COALESCE(?, a.check_in_lat),
    check_in_lng       = COALESCE(?, a.check_in_lng),
    check_in_photo_url = COALESCE(?, a.check_in_photo_url),
    check_in_ip        = COALESCE(?, a.check_in_ip),
    status             = ?,   -- tetap set PRESENT
    activity           = ?,   -- update activity terakhir
    updated_at         = NOW()
  WHERE NOT EXISTS (SELECT 1 FROM ins)
    AND a.user_id = ?
    AND a.date = ?::date
  RETURNING
    id, user_id, date,
    check_in_at, check_in_lat, check_in_lng, check_in_photo_url, check_in_ip,
    check_out_at, check_out_ip, total_minutes, status, activity, notes,
    created_at, updated_at,
    FALSE AS created
)
SELECT * FROM ins
UNION ALL
SELECT * FROM upd
LIMIT 1;
`

	qGetByUserAndDate = `
SELECT
  id, user_id, date,
  check_in_at, check_in_lat, check_in_lng, check_in_photo_url, check_in_ip,
  check_out_at, check_out_ip, total_minutes, status, activity, notes,
  created_at, updated_at
FROM attendances
WHERE user_id = ? AND date = ?::date
LIMIT 1;
`

	// Check-out hanya kalau belum pernah check-out (check_out_at IS NULL)
	// total_minutes dihitung dari (now - check_in_at) dalam menit, dibatasi >= 0.
	qCheckOut = `
UPDATE attendances a
SET
  check_out_at = NOW(),
  check_out_ip = ?,
  total_minutes = CASE
                    WHEN a.check_in_at IS NULL THEN a.total_minutes
                    ELSE GREATEST(0, a.total_minutes + CAST(EXTRACT(EPOCH FROM (NOW() - a.check_in_at))/60 AS INT))
                  END,
  updated_at = NOW()
WHERE a.user_id = ? AND a.date = ?::date AND a.check_out_at IS NULL
RETURNING
  id, user_id, date,
  check_in_at, check_in_lat, check_in_lng, check_in_photo_url, check_in_ip,
  check_out_at, check_out_ip, total_minutes, status, activity, notes,
  created_at, updated_at;
`
)

func (r *AttendanceRepository) UpsertCheckIn(ctx context.Context, a Attendance) (Attendance, bool, error) {
	var row struct {
		Attendance
		Created bool `gorm:"column:created"`
	}

	dateStr := a.Date.Format("2006-01-02")
	res := r.db.WithContext(ctx).Raw(
		qUpsertCheckIn,
		// INS args
		a.UserID, dateStr, a.CheckInLat, a.CheckInLng, a.CheckInPhotoURL, a.CheckInIP, a.Status, a.Activity,
		// UPD args
		a.CheckInLat, a.CheckInLng, a.CheckInPhotoURL, a.CheckInIP, a.Status, a.Activity, a.UserID, dateStr,
	).Scan(&row)

	if res.Error != nil {
		return Attendance{}, false, res.Error
	}
	return row.Attendance, row.Created, nil
}

func (r *AttendanceRepository) GetByUserAndDate(ctx context.Context, userID uint, date time.Time) (Attendance, error) {
	var out Attendance
	res := r.db.WithContext(ctx).Raw(qGetByUserAndDate, userID, date.Format("2006-01-02")).Scan(&out)
	if res.Error != nil {
		return Attendance{}, res.Error
	}
	if res.RowsAffected == 0 {
		return Attendance{}, gorm.ErrRecordNotFound
	}
	return out, nil
}

func (r *AttendanceRepository) CheckOut(ctx context.Context, userID uint, date time.Time, ip *string) (Attendance, error) {
	var out Attendance
	dateStr := date.Format("2006-01-02")

	res := r.db.WithContext(ctx).Raw(qCheckOut, ip, userID, dateStr).Scan(&out)
	if res.Error != nil {
		return Attendance{}, res.Error
	}
	if res.RowsAffected == 0 {
		// tidak ada baris yang ter-update → belum check-in atau sudah check-out
		return Attendance{}, gorm.ErrRecordNotFound
	}
	return out, nil
}
