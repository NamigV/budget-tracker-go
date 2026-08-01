package domain

import "time"

type Transaction struct {
	ID          int64
	UserID      int64
	CategoryID  int64
	AmountCents int64
	Currency    string
	Comment     *string
	OccurredOn  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
