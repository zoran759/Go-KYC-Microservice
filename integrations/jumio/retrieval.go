package jumio

import "time"

var attempts = [10]time.Duration{
	40 * time.Second,
	60 * time.Second,
	100 * time.Second,
	160 * time.Second,
	240 * time.Second,
	340 * time.Second,
	460 * time.Second,
	600 * time.Second,
	760 * time.Second,
	940 * time.Second,
}
