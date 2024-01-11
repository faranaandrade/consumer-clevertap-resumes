package location

import "time"

type Locater interface {
	DateNow() time.Time
}
