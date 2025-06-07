package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vnnyx/employee-management/internal/attendance/entity"
	"github.com/vnnyx/employee-management/internal/attendance/repository"
)

func TestStoreNewAttendance(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewAttendanceRepository(mock)

	now := time.Now()

	tests := []struct {
		name      string
		setupMock func()
		input     entity.Attendance
		expectErr bool
	}{
		{
			name: "success - inserted",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO attendances").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("att-id-1"))
			},
			input: entity.Attendance{
				ID:             "att-id-1",
				UserID:         "user-1",
				AttendanceDate: now,
				CreatedAt:      now,
				UpdatedAt:      now,
				CreatedBy:      "admin",
				UpdatedBy:      "admin",
				IPAddress:      "127.0.0.1",
			},
			expectErr: false,
		},
		{
			name: "error - db returns empty ID",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO attendances").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(""))
			},
			input: entity.Attendance{
				ID: "att-id-2",
			},
			expectErr: true,
		},
		{
			name: "error - pgxscan.Get() fails",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO attendances").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnError(errors.New("scan error"))
			},
			input: entity.Attendance{
				ID: "att-id-3",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}
			err := repo.StoreNewAttendance(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpsertAttendance(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewAttendanceRepository(mock)
	now := time.Now()

	tests := []struct {
		name      string
		setupMock func()
		input     entity.Attendance
		expectErr bool
	}{
		{
			name: "success - upserted",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO attendances").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("att-id-1"))
			},
			input: entity.Attendance{
				ID:             "att-id-1",
				UserID:         "user-1",
				AttendanceDate: now,
				CreatedAt:      now,
				UpdatedAt:      now,
				CreatedBy:      "admin",
				UpdatedBy:      "admin",
				IPAddress:      "127.0.0.1",
			},
			expectErr: false,
		},
		{
			name: "error - db query fails",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO attendances").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnError(errors.New("db error"))
			},
			input: entity.Attendance{
				ID:             "att-id-2",
				UserID:         "user-2",
				AttendanceDate: now,
				CreatedAt:      now,
				UpdatedAt:      now,
				CreatedBy:      "admin",
				UpdatedBy:      "admin",
				IPAddress:      "127.0.0.1",
			},
			expectErr: true,
		},
		{
			name: "error - empty returned ID",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO attendances").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(""))
			},
			input: entity.Attendance{
				ID:             "att-id-3",
				UserID:         "user-3",
				AttendanceDate: now,
				CreatedAt:      now,
				UpdatedAt:      now,
				CreatedBy:      "admin",
				UpdatedBy:      "admin",
				IPAddress:      "127.0.0.1",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}
			err := repo.UpsertAttendance(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStoreNewAttendancePeriod(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewAttendanceRepository(mock)

	now := time.Now()

	tests := []struct {
		name      string
		setupMock func()
		input     entity.AttendancePeriod
		expectErr bool
	}{
		{
			name: "success - inserted",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO attendance_periods").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("period-1"))
			},
			input: entity.AttendancePeriod{
				ID:        "period-1",
				StartDate: now,
				EndDate:   now.Add(24 * time.Hour),
				CreatedAt: now,
				UpdatedAt: now,
				CreatedBy: "admin",
				UpdatedBy: "admin",
				IPAddress: "127.0.0.1",
			},
			expectErr: false,
		},
		{
			name: "error - db returns empty ID",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO attendance_periods").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(""))
			},
			input: entity.AttendancePeriod{
				ID:        "period-2",
				StartDate: now,
				EndDate:   now.Add(24 * time.Hour),
				CreatedAt: now,
				UpdatedAt: now,
				CreatedBy: "admin",
				UpdatedBy: "admin",
				IPAddress: "127.0.0.1",
			},
			expectErr: true,
		},
		{
			name: "error - pgxscan.Get() fails",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO attendance_periods").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnError(errors.New("scan error"))
			},
			input: entity.AttendancePeriod{
				ID:        "period-3",
				StartDate: now,
				EndDate:   now.Add(24 * time.Hour),
				CreatedAt: now,
				UpdatedAt: now,
				CreatedBy: "admin",
				UpdatedBy: "admin",
				IPAddress: "127.0.0.1",
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}
			err := repo.StoreNewAttendancePeriod(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindPeriodByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewAttendanceRepository(mock)

	tests := []struct {
		name      string
		setupMock func()
		inputID   string
		expectNil bool
		expectErr bool
	}{
		{
			name: "found",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM attendance_periods WHERE id =").
					WithArgs("period-1").
					WillReturnRows(
						pgxmock.NewRows([]string{
							"id", "start_date", "end_date", "created_at", "updated_at", "created_by", "updated_by", "ip_address",
						}).AddRow("period-1", time.Now(), time.Now(), time.Now(), time.Now(), "admin", "admin", "127.0.0.1"),
					)
			},
			inputID:   "period-1",
			expectNil: false,
			expectErr: false,
		},
		{
			name: "not found",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM attendance_periods WHERE id =").
					WithArgs("missing").
					WillReturnError(pgx.ErrNoRows)
			},
			inputID:   "missing",
			expectNil: true,
			expectErr: false,
		},
		{
			name: "query error",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM attendance_periods WHERE id =").
					WithArgs("fail").
					WillReturnError(errors.New("db fail"))
			},
			inputID:   "fail",
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := repo.FindPeriodByID(context.Background(), tt.inputID)

			if tt.expectErr {
				assert.Error(t, err)
			} else if tt.expectNil {
				assert.NoError(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestFindAttendanceByPeriod(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewAttendanceRepository(mock)

	now := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	startDate := now.AddDate(0, 0, -1)
	endDate := now.AddDate(0, 0, 1)

	tests := []struct {
		name       string
		opts       []entity.FindAttendanceOptions
		setupMock  func()
		expectErr  bool
		expectMap  bool
		expectedBy entity.MappedBy
	}{
		{
			name: "success - mapped by user ID",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "attendance_date", "created_at",
					"updated_at", "created_by", "updated_by", "ip_address",
				}).
					AddRow("1", "user-1", now, now, now, "admin", "admin", "127.0.0.1").
					AddRow("2", "user-1", now, now, now, "admin", "admin", "127.0.0.1")

				mock.ExpectQuery("SELECT (.+) FROM attendances").
					WithArgs(startDate, endDate).
					WillReturnRows(rows)
			},
			opts: []entity.FindAttendanceOptions{
				{MappedOptions: &entity.MappedOptions{MappedBy: entity.MappedByUserID}},
			},
			expectErr:  false,
			expectMap:  true,
			expectedBy: entity.MappedByUserID,
		},
		{
			name: "success - mapped by attendance_date",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "attendance_date", "created_at",
					"updated_at", "created_by", "updated_by", "ip_address",
				}).
					AddRow("1", "user-1", now, now, now, "admin", "admin", "127.0.0.1").
					AddRow("2", "user-2", now, now, now, "admin", "admin", "127.0.0.1")

				mock.ExpectQuery("SELECT (.+) FROM attendances").
					WithArgs(startDate, endDate).
					WillReturnRows(rows)
			},
			opts: []entity.FindAttendanceOptions{
				{MappedOptions: &entity.MappedOptions{MappedBy: entity.MappedByAttendanceDate}},
			},
			expectErr:  false,
			expectMap:  true,
			expectedBy: entity.MappedByAttendanceDate,
		},
		{
			name: "error - scan fails",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "attendance_date", "created_at",
					"updated_at", "created_by", "updated_by", "ip_address",
				}).
					AddRow("bad", "user-1", now, now, now, "admin", "admin", "127.0.0.1").
					RowError(0, errors.New("scan error"))

				mock.ExpectQuery("SELECT (.+) FROM attendances").
					WithArgs(startDate, endDate).
					WillReturnRows(rows)
			},
			opts:      nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := repo.FindAttendanceByPeriod(context.Background(), startDate, endDate, tt.opts...)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, result.List, 2)

			if tt.expectMap {
				assert.True(t, result.IsMapped)
				assert.Equal(t, tt.expectedBy, result.MappedBy)
				assert.NotEmpty(t, result.Mapped)
			} else {
				assert.False(t, result.IsMapped)
			}
		})
	}
}

func TestFindAttendancePeriodByPayrollID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewAttendanceRepository(mock)

	tests := []struct {
		name      string
		setupMock func()
		inputID   string
		expectNil bool
		expectErr bool
	}{
		{
			name: "found",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM attendance_periods ad JOIN payrolls p ON ad.id = p.period_id WHERE p.id =").
					WithArgs("payroll-1").
					WillReturnRows(
						pgxmock.NewRows([]string{
							"id", "start_date", "end_date", "created_at", "updated_at", "created_by", "updated_by", "ip_address",
						}).AddRow("period-1", time.Now(), time.Now(), time.Now(), time.Now(), "admin", "admin", "127.0.0.1"),
					)
			},
			inputID:   "payroll-1",
			expectNil: false,
			expectErr: false,
		},
		{
			name: "not found",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM attendance_periods ad JOIN payrolls p ON ad.id = p.period_id WHERE p.id =").
					WithArgs("missing").
					WillReturnError(pgx.ErrNoRows)
			},
			inputID:   "missing",
			expectNil: true,
			expectErr: false,
		},
		{
			name: "query error",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM attendance_periods ad JOIN payrolls p ON ad.id = p.period_id WHERE p.id =").
					WithArgs("fail").
					WillReturnError(errors.New("db fail"))
			},
			inputID:   "fail",
			expectNil: true,
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := repo.FindAttendancePeriodByPayrollID(context.Background(), tt.inputID)

			if tt.expectErr {
				assert.Error(t, err)
			} else if tt.expectNil {
				assert.NoError(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
