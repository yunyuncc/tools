package time_cost

import (
	"fmt"
	"testing"
	"time"
)

func costTimeFunc() {
	time.Sleep(time.Second * 3)
}
func TestRun(t *testing.T) {
	var cost TimeCost
	cost.Begin()
	costTimeFunc()
	fmt.Printf("%3d ms cost\n", cost.End())
}
