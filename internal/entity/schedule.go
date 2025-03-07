package entity

import (
	"time"

	"github.com/google/uuid"
)

type DailySchedule struct {
	ID        uuid.UUID
	DayOfWeek int
	FromTime  time.Time
	ToTime    time.Time
}
