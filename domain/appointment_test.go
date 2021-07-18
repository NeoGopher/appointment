package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var created_at = time.Now()

func TestRepo_GetDoctorID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := NewAppointmentRepository(db)

	tests := []struct {
		name       string
		s          repoInterface
		doctorName string
		mock       func()
		want       int
		wantErr    bool
	}{
		{
			// When everything works as expected
			name:       "OK",
			s:          s,
			doctorName: "Doctor1",
			mock: func() {
				// We added one row
				rows := sqlmock.NewRows([]string{"Id"}).AddRow(1)
				mock.ExpectPrepare("SELECT (.+) FROM doctor").ExpectQuery().WithArgs("Doctor1").WillReturnRows(rows)
			},
			want: 1,
		},
		{
			// Doctor not present in database
			name:       "Not Found",
			s:          s,
			doctorName: "Doctor1",
			mock: func() {
				// We added no row
				rows := sqlmock.NewRows([]string{"Id"})
				mock.ExpectPrepare("SELECT (.+) FROM doctor").ExpectQuery().WithArgs("Doctor1").WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			// Invalid Prepare
			name:       "Not Found",
			s:          s,
			doctorName: "Doctor1",
			mock: func() {
				// Incorrect SQL statement
				rows := sqlmock.NewRows([]string{"Id"}).AddRow(1)
				mock.ExpectPrepare("SELECT (.+) FROM dummy").ExpectQuery().WithArgs("Doctor1").WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			tt.mock()
			got, err := tt.s.GetDoctorID(tt.doctorName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error new = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
