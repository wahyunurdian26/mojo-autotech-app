package attedance

import "time"

type Attendance struct {
	ID              uint       `json:"id"                 gorm:"primaryKey"`
	UserID          uint       `json:"user_id"            gorm:"index:idx_user_date,unique,priority:1"`
	WorkLocationID  *uint      `json:"work_location_id"`
	ShiftID         *uint      `json:"shift_id"`
	Date            time.Time  `json:"date"               gorm:"type:date;index:idx_user_date,unique,priority:2"`
	CheckInAt       *time.Time `json:"check_in_at"        gorm:"type:timestamptz"`
	CheckInLat      *float64   `json:"check_in_lat"`
	CheckInLng      *float64   `json:"check_in_lng"`
	CheckInPhotoURL *string    `json:"check_in_photo_url"`
	CheckInIP       *string    `json:"check_in_ip"        gorm:"type:inet"`
	CheckOutAt      *time.Time `json:"check_out_at"       gorm:"type:timestamptz"`
	CheckOutIP      *string    `json:"check_out_ip"       gorm:"type:inet"`
	TotalMinutes    int        `json:"total_minutes"`
	Status          string     `json:"status"` // PRESENT/LATE/ABSENT/ON_LEAVE
	Activity        string     `json:"activity"`
	CreatedAt       time.Time  `json:"created_at"        gorm:"type:timestamptz"`
	UpdatedAt       time.Time  `json:"updated_at"        gorm:"type:timestamptz"`
}
