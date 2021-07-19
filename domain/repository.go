package domain

import (
	"appointment/errors"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var Repo repoInterface = &apptRepo{}

type repoInterface interface {
	CreateDoctorAccount(string) (int, errors.AppointmentErr)
	CreatePatientAccount(string) (int, errors.AppointmentErr)
	GetDoctorID(string) (int, errors.AppointmentErr)
	CheckScheduleExists(int, time.Time, time.Time) (bool, errors.AppointmentErr)
	CheckScheduleOverlaps(int, time.Time, time.Time) (bool, errors.AppointmentErr)
	AddSchedule(int, time.Time, time.Time) errors.AppointmentErr
	CheckSlotAvailable(int, time.Time) (bool, errors.AppointmentErr)
	CheckSlotWithinSchedule(int, time.Time) (bool, errors.AppointmentErr)
	BookSlot(int, int, time.Time) (int, errors.AppointmentErr)
	ListSchedule(int) ([]Appointment, errors.AppointmentErr)
	CancelAppointment(int, int, string) errors.AppointmentErr
	InitializeDB() *sql.DB
	CloseDB()
}

type apptRepo struct {
	db *sql.DB
}

func (ar *apptRepo) InitializeDB() *sql.DB {
	var err error
	ar.db, err = sql.Open("sqlite3", "./appointments.db")

	if err != nil {
		log.Fatal(err)
	}

	err = AutoMigrate(ar.db)
	if err != nil {
		log.Printf("Error occured while setting up DB : %s\n", err)
		ar.db.Close()
		os.Exit(1)
	}

	return ar.db
}

func AutoMigrate(db *sql.DB) error {
	query, err := ioutil.ReadFile("database/schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(query))
	if err != nil {
		return err
	}

	return nil
}

func (ar *apptRepo) CloseDB() {
	ar.db.Close()
}

func NewAppointmentRepository(db *sql.DB) repoInterface {
	return &apptRepo{
		db: db,
	}
}

func (ar *apptRepo) CreateDoctorAccount(name string) (int, errors.AppointmentErr) {
	var id int
	query := "SELECT COUNT(id) FROM doctor WHERE name=?;"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return id, errors.NewInternalServerError("error occured when preparing statement to check for existing doctor account", err)
	}
	defer stmt.Close()

	var count int

	result := stmt.QueryRow(name)
	if err = result.Scan(&count); err != nil {
		return id, errors.NewInternalServerError("error occured when executing statement to check for existing doctor account", err)
	}

	if count != 0 {
		return id, errors.NewGeneralError("account already exists", nil)
	}

	query = "INSERT INTO doctor(name) VALUES (?);"

	stmt, err = ar.db.Prepare(query)
	if err != nil {
		return id, errors.NewInternalServerError("error occured when preparing statement to create new doctor account", err)
	}
	defer stmt.Close()

	result2, err := stmt.Exec(name)
	if err != nil {
		return id, errors.NewInternalServerError("error occured when executing statement to create new doctor account", err)
	}

	newId, err := result2.LastInsertId()
	if err != nil {
		return id, errors.NewInternalServerError("error occured when getting doctor ID", err)
	}

	id = int(newId)

	return id, nil
}

func (ar *apptRepo) CreatePatientAccount(name string) (int, errors.AppointmentErr) {
	var id int
	query := "SELECT COUNT(id) FROM patient WHERE name=?;"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return id, errors.NewInternalServerError("error occured when preparing statement to check for existing patient account", err)
	}
	defer stmt.Close()

	var count int

	result := stmt.QueryRow(name)
	if err = result.Scan(&count); err != nil {
		return id, errors.NewInternalServerError("error occured when executing statement to check for existing patient account", err)
	}

	if count != 0 {
		return id, errors.NewGeneralError("account already exists", nil)
	}

	query = "INSERT INTO patient(name) VALUES (?);"

	stmt, err = ar.db.Prepare(query)
	if err != nil {
		return id, errors.NewInternalServerError("error occured when preparing statement to create new patient account", err)
	}
	defer stmt.Close()

	result2, err := stmt.Exec(name)
	if err != nil {
		return id, errors.NewInternalServerError("error occured when executing statement to create new patient account", err)
	}

	newId, err := result2.LastInsertId()
	if err != nil {
		return id, errors.NewInternalServerError("error occured when getting patient ID", err)
	}

	id = int(newId)

	return id, nil
}

func (ar *apptRepo) GetDoctorID(doctorName string) (int, errors.AppointmentErr) {
	var doctorID int

	query := "SELECT id FROM doctor WHERE name=?;"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return doctorID, errors.NewInternalServerError("error occured when preparing statement to fetch doctor ID", err)
	}
	defer stmt.Close()

	result := stmt.QueryRow(doctorName)
	if err = result.Scan(&doctorID); err != nil {
		return doctorID, errors.NewNotFoundError(fmt.Sprintf("Doctor %s not found in database", doctorName), err)
	}

	return doctorID, nil
}

func (ar *apptRepo) CheckScheduleExists(doctorID int, startTime, endTime time.Time) (bool, errors.AppointmentErr) {
	query := "SELECT COUNT(id) FROM doctor_schedule WHERE doctor_id=? and start_time=? and end_time=?;"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return false, errors.NewInternalServerError("error occured when preparing statement to fetch Doctor schedule", err)
	}
	defer stmt.Close()

	var count int

	result := stmt.QueryRow(doctorID, startTime, endTime)
	if err = result.Scan(&count); err != nil {
		return false, errors.NewInternalServerError("error occured when executing statement to fetch Doctor schedule", err)
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (ar *apptRepo) CheckScheduleOverlaps(doctorID int, startTime, endTime time.Time) (bool, errors.AppointmentErr) {
	query := "SELECT start_time, end_time FROM doctor_schedule WHERE doctor_id=? AND start_time>=DATE('now') and end_time<=DATE('now', '+1 day');"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return false, errors.NewInternalServerError("error occured when preparing statement to fetch Doctor schedule", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(doctorID)
	if err != nil {
		return false, errors.NewInternalServerError("error occured when executing statement to fetch Doctor schedule", err)
	}
	defer rows.Close()

	for rows.Next() {
		var start_time time.Time
		var end_time time.Time

		err := rows.Scan(&start_time, &end_time)
		if err != nil {
			return false, errors.NewInternalServerError("error occured when parsing Doctor schedule", err)
		}

		if (startTime.Equal(start_time)) || (endTime.Equal(end_time)) || (startTime.After(start_time) && startTime.Before(end_time)) || (endTime.Before(end_time) && endTime.After(start_time)) {
			return true, nil
		}
	}

	return false, nil
}

func (ar *apptRepo) AddSchedule(doctorID int, startTime time.Time, endTime time.Time) errors.AppointmentErr {
	query := "INSERT INTO doctor_schedule (doctor_id, start_time, end_time) VALUES (?,?,?);"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return errors.NewInternalServerError("error occured when preparing statement to create Doctor schedule in database", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(doctorID, startTime, endTime)
	if err != nil {
		return errors.NewInternalServerError("error occured when executing statement to create Doctor schedule in database", err)
	}

	return nil
}

func (ar *apptRepo) CheckSlotAvailable(doctorID int, startTime time.Time) (bool, errors.AppointmentErr) {
	query := "SELECT count(id) FROM appointments WHERE doctor_id=? AND is_active=1 AND start_time=?;"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return false, errors.NewInternalServerError("error occured when preparing statement to check for available slots in database", err)
	}
	defer stmt.Close()

	var count int

	result := stmt.QueryRow(doctorID, startTime)
	if err = result.Scan(&count); err != nil {
		return false, errors.NewInternalServerError("error occured when executing statement to check for available slots in database", err)
	}

	if count != 0 {
		return false, nil
	}

	return true, nil
}

func (ar *apptRepo) CheckSlotWithinSchedule(doctorID int, startTime time.Time) (bool, errors.AppointmentErr) {
	query := "SELECT start_time, end_time FROM doctor_schedule WHERE doctor_id=? AND start_time>=DATE('now') and end_time<=DATE('now', '+1 day');"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return false, errors.NewInternalServerError("error occured when preparing statement to fetch Doctor schedule", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(doctorID)
	if err != nil {
		return false, errors.NewInternalServerError("error occured when executing statement to fetch Doctor schedule", err)
	}
	defer rows.Close()

	for rows.Next() {
		var start_time time.Time
		var end_time time.Time

		err := rows.Scan(&start_time, &end_time)
		if err != nil {
			return false, errors.NewInternalServerError("error occured when parsing Doctor schedule", err)
		}

		if (startTime.Equal(start_time)) || (startTime.After(start_time) && startTime.Before(end_time)) {
			return true, nil
		}
	}

	return false, nil
}

func (ar *apptRepo) BookSlot(doctorID int, userID int, startTime time.Time) (int, errors.AppointmentErr) {
	var appointmentID int

	query := "INSERT INTO appointments(doctor_id, patient_id, start_time, is_active) VALUES (?, ?, ?, 1);"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return appointmentID, errors.NewInternalServerError("error occured when preparing statement for booking slot in database", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(doctorID, userID, startTime)
	if err != nil {
		return appointmentID, errors.NewInternalServerError("error occured when executing statement for booking slot in database", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return appointmentID, errors.NewInternalServerError("error occured when getting appointment ID", err)
	}

	appointmentID = int(id)

	return appointmentID, nil
}

func (ar *apptRepo) ListSchedule(doctorID int) ([]Appointment, errors.AppointmentErr) {
	appointments := make([]Appointment, 0)

	// Get Booked Appointments
	query := "SELECT id, patient_id, start_time FROM appointments WHERE doctor_id=? AND is_active=1 AND start_time>=DATE('now') AND start_time<=DATE('now', '+1 day') ORDER BY start_time;"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		return appointments, errors.NewInternalServerError("error occured when preparing statement to fetch Booked Appointments", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(doctorID)
	if err != nil {
		return appointments, errors.NewInternalServerError("error occured when executing statement to fetch Booked Appointments", err)
	}
	defer rows.Close()

	bookedAppointments := make(map[time.Time]struct {
		AppointmentID int
		PatientID     int
	})

	for rows.Next() {
		var aptID, patID int
		var st time.Time

		err := rows.Scan(&aptID, &patID, &st)
		if err != nil {
			return appointments, errors.NewInternalServerError("error occured when parsing Booked Appointments", err)
		}

		bookedAppointments[st] = struct {
			AppointmentID int
			PatientID     int
		}{aptID, patID}
	}

	// Get Schedule
	query = "SELECT start_time, end_time FROM doctor_schedule WHERE doctor_id=? AND start_time>=DATE('now') AND end_time<=DATE('now', '+1 day') ORDER BY start_time;"

	stmt, err = ar.db.Prepare(query)
	if err != nil {
		return appointments, errors.NewInternalServerError("error occured when preparing statement to fetch Doctor schedule", err)
	}
	defer stmt.Close()

	rows, err = stmt.Query(doctorID)
	if err != nil {
		return appointments, errors.NewInternalServerError("error occured when executing statement to fetch Doctor schedule", err)
	}
	defer rows.Close()

	for rows.Next() {
		var start_time time.Time
		var end_time time.Time

		err := rows.Scan(&start_time, &end_time)
		if err != nil {
			return appointments, errors.NewInternalServerError("error occured when parsing Doctor schedule", err)
		}

		t := start_time
		for t.Before(end_time) {
			doctorID := strconv.Itoa(doctorID)

			appointment := Appointment{
				ID:        "",
				DoctorID:  doctorID,
				PatientID: "",
				StartTime: t,
				Booked:    false,
			}

			if data, ok := bookedAppointments[t]; ok {
				appointment.ID = strconv.Itoa(data.AppointmentID)
				appointment.PatientID = strconv.Itoa(data.PatientID)
				appointment.Booked = true
			}

			appointments = append(appointments, appointment)

			t = t.Add(15 * time.Minute)

		}
	}

	return appointments, nil
}

func (ar *apptRepo) CancelAppointment(appointmentID int, userID int, userType string) errors.AppointmentErr {
	query := "SELECT doctor_id, patient_id, is_active FROM appointments WHERE id=?;"

	stmt, err := ar.db.Prepare(query)
	if err != nil {
		errors.NewInternalServerError("error occured when preparing statement to check for appointment status", err)
	}
	defer stmt.Close()

	var doctorID int
	var patientID int
	var activeStatus int

	result := stmt.QueryRow(appointmentID)
	if err = result.Scan(&doctorID, &patientID, &activeStatus); err != nil {
		return errors.NewBadRequestError(fmt.Sprintf("appointment id %d does not exist in database", appointmentID), err)
	}

	if activeStatus == 0 {
		return errors.NewGeneralError(fmt.Sprintf("appointment id %d is already cancelled", appointmentID), nil)
	}

	if (userType == "doctor" && userID != doctorID) || (userType == "patient" && userID != patientID) {
		return errors.NewGeneralForbiddenError("unauthorised to perform this action", nil)
	}

	query = "UPDATE appointments SET deleted_at=CURRENT_TIMESTAMP, is_active=0 WHERE id=?;"

	stmt, err = ar.db.Prepare(query)
	if err != nil {
		return errors.NewInternalServerError("error occured when preparing statement to cancel slot", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(appointmentID)
	if err != nil {
		return errors.NewInternalServerError("error occured when executing statement to cancel slot", err)
	}

	return nil
}
