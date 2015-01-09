package booking

import "time"

type Service struct {
	dataPath       string
	processorToken string
}

func NewService(dataPath, processorToken string) Service {
	s := Service{}
	s.dataPath = dataPath
	s.processorToken = processorToken
	return s
}

// Returns list of available dates in the future
func (s Service) AvailableDays() ([]time.Time, error) {
	days := []time.Time{
		time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC),
	}
	return days, nil
}
