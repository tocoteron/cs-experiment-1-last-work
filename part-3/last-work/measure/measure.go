package measure

import "time"

func MeasureFuncTime(f func() error) (float64, error) {
	start := time.Now()
	err := f()
	end := time.Now()

	return (end.Sub(start)).Seconds(), err
}
