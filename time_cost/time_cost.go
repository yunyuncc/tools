package time_cost

import "time"

type TimeCost struct {
	start time.Time
}

func (t *TimeCost) Begin() {
	t.start = time.Now()
}
func (t *TimeCost) End() int64 {
	end := time.Now()
	latency := end.Sub(t.start)

	ms := latency.Nanoseconds() / 1000000
	return ms
}
