package common

import "time"

const (
	ReqTimeout             = time.Minute * 1
	SchedulerStopTimeout   = time.Second * 15
	RequestTimeout         = time.Minute * 5
	ConnTimeout            = time.Second * 10
	ReqTimeValidCoins      = time.Second * 30
	ReqTimePrice           = time.Second * 30
	PostgresDBQueryTimeout = time.Minute * 5

	statusOK        = "OK"
	statusNotFound  = "NotFound"
	statusError     = "ERROR"
	DefaultCurrency = "USD"
	Empty           = ""

	Precision = 8
	Zero      = 0
)
