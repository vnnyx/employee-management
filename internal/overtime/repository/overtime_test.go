package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vnnyx/employee-management/internal/overtime/entity"
	"github.com/vnnyx/employee-management/internal/overtime/repository"
)

func TestStoreNewOvertime(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewOvertimeRepository(mock)
	now := time.Now()

	tests := []struct {
		name      string
		setupMock func()
		input     entity.Overtime
		expectErr bool
	}{
		{
			name: "success",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO overtimes").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("ot-1"))
			},
			input: entity.Overtime{
				ID:            "ot-1",
				UserID:        "user-1",
				OverTimeDate:  now,
				OvertimeHours: 3,
				CreatedAt:     now,
				UpdatedAt:     now,
				CreatedBy:     "admin",
				UpdatedBy:     "admin",
				IPAddress:     "127.0.0.1",
			},
			expectErr: false,
		},
		{
			name: "db error",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO overtimes").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("insert failed"))
			},
			input:     entity.Overtime{},
			expectErr: true,
		},
		{
			name: "empty returned ID",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO overtimes").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(""))
			},
			input:     entity.Overtime{ID: "ot-2"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.StoreNewOvertime(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpsertOvertime(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewOvertimeRepository(mock)
	now := time.Now()

	tests := []struct {
		name      string
		setupMock func()
		input     entity.Overtime
		expectErr bool
	}{
		{
			name: "success",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO overtimes").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("ot-1"))
			},
			input: entity.Overtime{
				ID:            "ot-1",
				UserID:        "user-1",
				OverTimeDate:  now,
				OvertimeHours: 3,
				CreatedAt:     now,
				UpdatedAt:     now,
				CreatedBy:     "admin",
				UpdatedBy:     "admin",
				IPAddress:     "127.0.0.1",
			},
			expectErr: false,
		},
		{
			name: "query error",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO overtimes").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("upsert failed"))
			},
			input:     entity.Overtime{ID: "ot-2"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.UpsertOvertime(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindOvertimeByUserIDDate(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewOvertimeRepository(mock)
	now := time.Now()

	tests := []struct {
		name      string
		setupMock func()
		userID    string
		date      time.Time
		expectNil bool
		expectErr bool
	}{
		{
			name: "found",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM overtimes WHERE user_id").
					WithArgs("user-1", now).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "user_id", "overtime_date", "overtime_hours",
						"created_at", "updated_at", "created_by", "updated_by", "ip_address",
					}).AddRow("ot-1", "user-1", now, 2, now, now, "admin", "admin", "127.0.0.1"))
			},
			userID:    "user-1",
			date:      now,
			expectNil: false,
			expectErr: false,
		},
		{
			name: "not found",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM overtimes WHERE user_id").
					WithArgs("user-2", now).
					WillReturnError(pgx.ErrNoRows)
			},
			userID:    "user-2",
			date:      now,
			expectNil: true,
			expectErr: false,
		},
		{
			name: "query error",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM overtimes WHERE user_id").
					WithArgs("user-3", now).
					WillReturnError(errors.New("db error"))
			},
			userID:    "user-3",
			date:      now,
			expectNil: true,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := repo.FindOvertimeByUserIDDate(context.Background(), tt.userID, tt.date)
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

func TestFindOvertimeByPeriod(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewOvertimeRepository(mock)
	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	startDate := now.AddDate(0, 0, -1)
	endDate := now.AddDate(0, 0, 1)

	sampleRows := pgxmock.NewRows([]string{
		"id", "user_id", "overtime_date", "overtime_hours",
		"created_at", "updated_at", "created_by", "updated_by", "ip_address",
	}).AddRow("ot-1", "user-1", now, 2, now, now, "admin", "admin", "127.0.0.1").
		AddRow("ot-2", "user-2", now, 3, now, now, "admin", "admin", "127.0.0.1")

	tests := []struct {
		name       string
		opts       []entity.FindOvertimeOptions
		setupMock  func()
		expectErr  bool
		expectMap  bool
		expectedBy entity.MappedBy
	}{
		{
			name: "success - no mapping",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM overtimes").
					WithArgs(startDate, endDate).
					WillReturnRows(sampleRows)
			},
			opts:      nil,
			expectErr: false,
			expectMap: false,
		},
		{
			name: "success - mapped by user ID",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM overtimes").
					WithArgs(startDate, endDate).
					WillReturnRows(
						pgxmock.NewRows([]string{
							"id", "user_id", "overtime_date", "overtime_hours",
							"created_at", "updated_at", "created_by", "updated_by", "ip_address",
						}).
							AddRow("ot-1", "user-1", now, 2*time.Hour, now, now, "admin", "admin", "127.0.0.1").
							AddRow("ot-2", "user-1", now, 3*time.Hour, now, now, "admin", "admin", "127.0.0.1"),
					)
			},

			opts: []entity.FindOvertimeOptions{
				{MappedOptions: &entity.MappedOptions{MappedBy: entity.MappedByUserID}},
			},
			expectErr:  false,
			expectMap:  true,
			expectedBy: entity.MappedByUserID,
		},
		{
			name: "success - mapped by overtime_date",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM overtimes").
					WithArgs(startDate, endDate).
					WillReturnRows(
						pgxmock.NewRows([]string{
							"id", "user_id", "overtime_date", "overtime_hours",
							"created_at", "updated_at", "created_by", "updated_by", "ip_address",
						}).
							AddRow("ot-1", "user-1", now, 2*time.Hour, now, now, "admin", "admin", "127.0.0.1").
							AddRow("ot-2", "user-1", now, 3*time.Hour, now, now, "admin", "admin", "127.0.0.1"),
					)
			},
			opts: []entity.FindOvertimeOptions{
				{MappedOptions: &entity.MappedOptions{MappedBy: entity.MappedByAttendanceDate}},
			},
			expectErr:  false,
			expectMap:  true,
			expectedBy: entity.MappedByAttendanceDate,
		},
		{
			name: "error - query fails",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM overtimes").
					WithArgs(startDate, endDate).
					WillReturnError(errors.New("query error"))
			},
			opts:      nil,
			expectErr: true,
		},
		{
			name: "error - scan fails",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "overtime_date", "overtime_hours",
					"created_at", "updated_at", "created_by", "updated_by", "ip_address",
				}).
					AddRow("bad-id", nil, nil, nil, nil, nil, nil, nil, nil).
					RowError(0, errors.New("scan error"))

				mock.ExpectQuery("SELECT (.+) FROM overtimes").
					WithArgs(startDate, endDate).
					WillReturnRows(rows)
			},
			opts:      nil,
			expectErr: true,
		},
		{
			name: "error - unsupported mapped option",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM overtimes").
					WithArgs(startDate, endDate).
					WillReturnRows(sampleRows)
			},
			opts: []entity.FindOvertimeOptions{
				{MappedOptions: &entity.MappedOptions{MappedBy: "unsupported"}},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := repo.FindOvertimeByPeriod(context.Background(), startDate, endDate, tt.opts...)

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
				assert.Nil(t, result.Mapped)
			}
		})
	}
}
