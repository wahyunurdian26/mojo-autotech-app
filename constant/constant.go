package constant

import "time"

const (
	ReqParamInvalid = "Request parameter is invalid"
	LoginError      = "Login Gagal"
	LoginSuccess    = "Login Sukses"

	AccessTTL  = 15 * time.Minute
	RefreshTTL = 7 * 24 * time.Hour
)
