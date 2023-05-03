package models

type Account struct {
	Model
	UserId           string `json:"user_id"`
	Balance          int64  `json:"balance"`
	Status           string `json:"status"`
	AvailableBalance int64  `json:"available_balance"`
	PendingBalance   int64  `json:"pending_balance"`
	LockedBalance    int64  `json:"locked_balance"`
}

type LockedBalance struct {
	Model
	LockDate     int64 `json:"lock_date"`
	ReleaseDate  int64 `json:"release_date"`
	AmountLocked int64 `json:"amount_locked"`
}
