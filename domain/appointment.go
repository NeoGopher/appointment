package domain

import (
	"time"
)

type Appointment struct {
	ID        string    `json:"appointmentid"`
	DoctorID  string    `json:"doctorid"`
	PatientID string    `json:"patientid"`
	StartTime time.Time `json:"starttime"`
	Booked    bool      `json:"booked"`
}
