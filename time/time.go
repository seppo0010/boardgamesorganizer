package time

import "time"

type Factory interface {
	Now() time.Time
}
