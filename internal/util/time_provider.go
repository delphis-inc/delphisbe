package util

import "time"

type TimeProvider interface {
	Now() time.Time
}

type RealTime struct{}

func (rt *RealTime) Now() time.Time {
	return time.Now()
}

type FrozenTime struct {
	NowTime time.Time
}

func (ft *FrozenTime) Now() time.Time {
	return ft.NowTime
}
