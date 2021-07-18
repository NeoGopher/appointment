package domain

import "time"

type Schedule struct {
	ID        int       `json:"scheduleId"`
	DoctorID  int       `json:"doctorId"`
	StartTime time.Time `json:"starttime"`
	EndTime   time.Time `json:"endtime"`
}
