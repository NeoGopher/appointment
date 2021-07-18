package services

import (
	"appointment/domain"
	"appointment/errors"
	"fmt"
	"time"
)

var AppointmentService appointmentServiceInterface = &appointmentService{}

type appointmentServiceInterface interface {
	CreateDoctorAccount(string) (int, errors.AppointmentErr)
	CreatePatientAccount(string) (int, errors.AppointmentErr)
	AddSchedule(int, time.Time, time.Time) errors.AppointmentErr
	Book(string, int, time.Time) (int, errors.AppointmentErr)
	ListSchedule(string) ([]domain.Appointment, errors.AppointmentErr)
	Cancel(int, int, string) errors.AppointmentErr
}

type appointmentService struct{}

func (as *appointmentService) CreateDoctorAccount(name string) (int, errors.AppointmentErr) {
	var id int

	id, err := domain.Repo.CreateDoctorAccount(name)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (as *appointmentService) CreatePatientAccount(name string) (int, errors.AppointmentErr) {
	var id int

	id, err := domain.Repo.CreatePatientAccount(name)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (as *appointmentService) AddSchedule(doctorID int, startTime time.Time, endTime time.Time) errors.AppointmentErr {
	// Allow current days bookings only
	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 23, 59, 59, 0, &time.Location{})

	if startTime.After(today) {
		return errors.NewGeneralError("Schedule can be created for current day only", nil)
	}

	// Check If Schedule already exists for Doctor
	scheduleExists, err := domain.Repo.CheckScheduleExists(doctorID, startTime, endTime)
	if err != nil {
		return err
	}

	if scheduleExists {
		return errors.NewGeneralError("Schedule already exists ", nil)
	}

	// Check If Schedule overlaps with existing
	scheduleOverlaps, err := domain.Repo.CheckScheduleOverlaps(doctorID, startTime, endTime)
	if err != nil {
		return err
	}

	if scheduleOverlaps {
		return errors.NewGeneralError("Schedule overlaps with existing schedule", nil)
	}

	// Else Add Schedule
	err = domain.Repo.AddSchedule(doctorID, startTime, endTime)
	if err != nil {
		return err
	}

	return nil
}

func (as *appointmentService) Book(doctorName string, userID int, startTime time.Time) (int, errors.AppointmentErr) {
	var appointmentID int

	// Get DoctorID for given DoctorName
	doctorID, err := domain.Repo.GetDoctorID(doctorName)
	if err != nil {
		return appointmentID, err
	}

	// Check If Appointment slot is available
	slotAvailable, err := domain.Repo.CheckSlotAvailable(doctorID, startTime)
	if err != nil {
		return appointmentID, err
	}

	// Cant book
	if !slotAvailable {
		return appointmentID, errors.NewGeneralError(fmt.Sprintf("Slot already taken"), nil)
	}

	// Check If Appointment within Doctor schedule
	slotWithinSchedule, err := domain.Repo.CheckSlotWithinSchedule(doctorID, startTime)
	if err != nil {
		return appointmentID, err
	}

	if !slotWithinSchedule {
		return appointmentID, errors.NewGeneralError(fmt.Sprintf("Slot not within schedule"), nil)
	}

	// Book
	appointmentID, err = domain.Repo.BookSlot(doctorID, userID, startTime)
	if err != nil {
		return appointmentID, err
	}

	return appointmentID, nil
}

func (as *appointmentService) ListSchedule(doctorName string) ([]domain.Appointment, errors.AppointmentErr) {
	appointments := make([]domain.Appointment, 0)

	// Get DoctorID for given DoctorName
	doctorID, err := domain.Repo.GetDoctorID(doctorName)
	if err != nil {
		return appointments, err
	}

	// List
	appointments, err = domain.Repo.ListSchedule(doctorID)
	if err != nil {
		return appointments, err
	}

	return appointments, nil
}

func (as *appointmentService) Cancel(appointID int, userID int, userType string) errors.AppointmentErr {
	// Check If Appointment id exists and active
	err := domain.Repo.CancelAppointment(appointID, userID, userType)
	if err != nil {
		return err
	}

	return nil
}
