package handlers

import (
	"appointment/domain"
	"appointment/errors"
	"appointment/services"
	"appointment/utilities"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ScheduleForm struct {
	StartTime time.Time `form:"starttime" json:"starttime" binding:"required,bookabledate,multipleoffifteen" time_format:"2006-01-02 15:04:05"`
	EndTime   time.Time `form:"endtime" json:"endtime" binding:"required,bookabledate,multipleoffifteen,gtfield=StartTime" time_format:"2006-01-02 15:04:05"`
	Token     string    `form:"token" json:"token" binding:"required"`
}

type BookAppointmentForm struct {
	DoctorName string    `form:"doctorname" json:"doctorname" binding:"required"`
	StartTime  time.Time `form:"starttime" json:"starttime" binding:"required,bookabledate,multipleoffifteen" time_format:"2006-01-02 15:04:05"`
	Token      string    `form:"token" json:"token" binding:"required"`
}

type ListAppointmentsForm struct {
	DoctorName string `form:"doctorname" json:"doctorname" binding:"required"`
}

type CancelAppointmentForm struct {
	AppointmentID int    `form:"appointmentid" json:"appointmentid" binding:"required"`
	Token         string `form:"token" json:"token" binding:"required"`
}

type SignupForm struct {
	Name string `form:"name" json:"name" binding:"required"`
	Type string `form:"type" json:"usertype" binding:"required"`
}

func SetSchedule(c *gin.Context) {
	var form ScheduleForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("error occured while parsing input", err))

		return
	}

	id, userType, err := utilities.ParseToken(form.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("error occured while parsing token", err))

		return
	}

	switch strings.ToLower(userType) {
	case "patient":
		c.JSON(http.StatusForbidden, errors.NewGeneralForbiddenError("unauthorised to perform this action", nil))

		return
	case "doctor":
		break
	default:
		c.JSON(http.StatusBadRequest, errors.NewGeneralError("unknown usertype", nil))

		return
	}

	doctorID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("error occured while parsing userID", err))

		return
	}

	if err := services.AppointmentService.AddSchedule(doctorID, form.StartTime, form.EndTime); err != nil {
		c.JSON(err.GetStatus(), err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Schedule created"})
}

func BookAppointment(c *gin.Context) {
	var form BookAppointmentForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("error occured while parsing input", err))

		return
	}

	id, userType, err := utilities.ParseToken(form.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("error occured while parsing token", err))

		return
	}

	switch strings.ToLower(userType) {
	case "patient":
		break
	case "doctor":
		c.JSON(http.StatusForbidden, errors.NewGeneralForbiddenError("unauthorised to perform this action", nil))

		return
	default:
		c.JSON(http.StatusBadRequest, errors.NewGeneralError("unknown usertype", nil))

		return
	}

	var appointmentID int
	var err2 errors.AppointmentErr

	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("error occured while parsing userID", err))

		return
	}

	appointmentID, err2 = services.AppointmentService.Book(form.DoctorName, userID, form.StartTime)
	if err2 != nil {
		c.JSON(err2.GetStatus(), err2)

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Appointment booked", "appointmentid": appointmentID})
}

func ListAppointments(c *gin.Context) {
	var form ListAppointmentsForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewGeneralError("error occured while parsing input", err))

		return
	}

	var appointments []domain.Appointment

	appointments, err := services.AppointmentService.ListSchedule(form.DoctorName)
	if err != nil {
		c.JSON(err.GetStatus(), err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Appointments Listed", "appointments": appointments})
}

func CancelAppointment(c *gin.Context) {
	var form CancelAppointmentForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("error occured while parsing input", err))

		return
	}

	id, userType, err := utilities.ParseToken(form.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("error occured while parsing token", err))

		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("error occured while parsing userID", err))

		return
	}

	userType = strings.ToLower(userType)

	if err := services.AppointmentService.Cancel(form.AppointmentID, userID, userType); err != nil {
		c.JSON(err.GetStatus(), err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Appointment cancelled"})
}

func Signup(c *gin.Context) {
	var form SignupForm

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewBadRequestError("error occured while parsing input", err))

		return
	}

	var userID int

	// Create User
	switch strings.ToLower(form.Type) {
	case "patient":
		var err errors.AppointmentErr

		userID, err = services.AppointmentService.CreatePatientAccount(form.Name)
		if err != nil {
			c.JSON(err.GetStatus(), err)

			return
		}
	case "doctor":
		var err errors.AppointmentErr

		userID, err = services.AppointmentService.CreateDoctorAccount(form.Name)
		if err != nil {
			c.JSON(err.GetStatus(), err)

			return
		}
	default:
		c.JSON(http.StatusBadRequest, errors.NewGeneralError("Unknown usertype", nil))

		return
	}

	data := strings.Join([]string{strconv.Itoa(userID), form.Type}, "|")

	token := utilities.GetCode(data)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Account created", "token": token})
}

var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if ok {
		today := time.Now()
		if today.After(date) {
			return false
		}
	}
	return true
}

var multipleOfFifteen validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)

	if ok {
		offset := date.Minute()

		if offset%15 != 0 {
			return false
		}

		return true
	}

	return false
}

func RegisterValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("multipleoffifteen", multipleOfFifteen)
	}
}
