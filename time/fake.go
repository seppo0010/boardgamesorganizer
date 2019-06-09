package time

import "time"

type Fake struct {
	CurrentNow time.Time
}

func (f *Fake) Now() time.Time {
	return f.CurrentNow
}

func (f *Fake) AdvaceTime(duration time.Duration) {
	f.CurrentNow = f.CurrentNow.Add(duration)
}
