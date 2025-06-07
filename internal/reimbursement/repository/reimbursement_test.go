package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vnnyx/employee-management/internal/reimbursement/entity"
	"github.com/vnnyx/employee-management/internal/reimbursement/repository"
	"github.com/vnnyx/employee-management/pkg/optional"
)

func TestStoreNewReimbursement(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewReimbursementRepository(mock)
	now := time.Now()

	tests := []struct {
		name      string
		setupMock func()
		input     entity.Reimbursement
		expectErr bool
	}{
		{
			name: "success",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO reimbursements").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("rb-1"))
			},
			input: entity.Reimbursement{
				ID:                "rb-1",
				UserID:            "user-1",
				Amount:            5000,
				Description:       optional.NewString("Travel expenses"),
				ReimbursementDate: now,
				CreatedAt:         now,
				UpdatedAt:         now,
				CreatedBy:         "admin",
				UpdatedBy:         "admin",
				IPAddress:         "127.0.0.1",
			},
			expectErr: false,
		},
		{
			name: "empty returned ID",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO reimbursements").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(""))
			},
			input:     entity.Reimbursement{ID: "rb-2"},
			expectErr: true,
		},
		{
			name: "query error",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO reimbursements").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(errors.New("insert failed"))
			},
			input:     entity.Reimbursement{ID: "rb-3"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.StoreNewReimbursement(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindReimbursementByPeriod(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewReimbursementRepository(mock)

	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	startDate := now.AddDate(0, 0, -1)
	endDate := now.AddDate(0, 0, 1)

	tests := []struct {
		name       string
		opts       []entity.FindReimbursementOptions
		setupMock  func()
		expectErr  bool
		expectMap  bool
		expectedBy entity.MappedBy
	}{
		{
			name: "success - no mapping",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "amount", "description", "reimbursement_date",
					"created_at", "updated_at", "created_by", "updated_by", "ip_address",
				}).
					AddRow("rb-1", "user-1", 1000, "desc", now, now, now, "admin", "admin", "127.0.0.1").
					AddRow("rb-2", "user-2", 2000, "desc", now, now, now, "admin", "admin", "127.0.0.1")

				mock.ExpectQuery("SELECT (.+) FROM reimbursements").
					WithArgs(startDate, endDate).
					WillReturnRows(rows)
			},
			opts:      nil,
			expectErr: false,
			expectMap: false,
		},
		{
			name: "success - mapped by user ID",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "amount", "description", "reimbursement_date",
					"created_at", "updated_at", "created_by", "updated_by", "ip_address",
				}).
					AddRow("rb-1", "user-1", 1000, "desc", now, now, now, "admin", "admin", "127.0.0.1").
					AddRow("rb-2", "user-1", 2000, "desc", now, now, now, "admin", "admin", "127.0.0.1")

				mock.ExpectQuery("SELECT (.+) FROM reimbursements").
					WithArgs(startDate, endDate).
					WillReturnRows(rows)
			},
			opts: []entity.FindReimbursementOptions{
				{MappedOptions: &entity.MappedOptions{MappedBy: entity.MappedByUserID}},
			},
			expectErr:  false,
			expectMap:  true,
			expectedBy: entity.MappedByUserID,
		},
		{
			name: "error - scan fails",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "amount", "description", "reimbursement_date",
					"created_at", "updated_at", "created_by", "updated_by", "ip_address",
				}).
					AddRow("bad", nil, nil, nil, nil, nil, nil, nil, nil, nil).
					RowError(0, errors.New("scan error"))

				mock.ExpectQuery("SELECT (.+) FROM reimbursements").
					WithArgs(startDate, endDate).
					WillReturnRows(rows)
			},
			opts:      nil,
			expectErr: true,
		},
		{
			name: "error - unsupported mapped option",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "amount", "description", "reimbursement_date",
					"created_at", "updated_at", "created_by", "updated_by", "ip_address",
				}).
					AddRow("rb-1", "user-1", 1000, "desc", now, now, now, "admin", "admin", "127.0.0.1")

				mock.ExpectQuery("SELECT (.+) FROM reimbursements").
					WithArgs(startDate, endDate).
					WillReturnRows(rows)
			},
			opts: []entity.FindReimbursementOptions{
				{MappedOptions: &entity.MappedOptions{MappedBy: "unsupported"}},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := repo.FindReimbursementByPeriod(context.Background(), startDate, endDate, tt.opts...)

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

func TestFindReimbursementByUserIDPeriod(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewReimbursementRepository(mock)

	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	startDate := now.AddDate(0, 0, -1)
	endDate := now.AddDate(0, 0, 1)

	tests := []struct {
		name      string
		userID    string
		setupMock func()
		expectErr bool
		expectLen int
	}{
		{
			name:   "success",
			userID: "user-1",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM reimbursements").
					WithArgs("user-1", startDate, endDate).
					WillReturnRows(pgxmock.NewRows([]string{
						"id", "user_id", "amount", "description", "reimbursement_date",
						"created_at", "updated_at", "created_by", "updated_by", "ip_address",
					}).AddRow("rb-1", "user-1", 1000, "desc", now, now, now, "admin", "admin", "127.0.0.1"))
			},
			expectErr: false,
			expectLen: 1,
		},
		{
			name:   "error - scan fails",
			userID: "user-2",
			setupMock: func() {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "amount", "description", "reimbursement_date",
					"created_at", "updated_at", "created_by", "updated_by", "ip_address",
				}).
					AddRow("bad", nil, nil, nil, nil, nil, nil, nil, nil, nil).
					RowError(0, errors.New("scan error"))

				mock.ExpectQuery("SELECT (.+) FROM reimbursements").
					WithArgs("user-2", startDate, endDate).
					WillReturnRows(rows)
			},
			expectErr: true,
		},
		{
			name:   "error - query fails",
			userID: "user-3",
			setupMock: func() {
				mock.ExpectQuery("SELECT (.+) FROM reimbursements").
					WithArgs("user-3", startDate, endDate).
					WillReturnError(errors.New("query error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := repo.FindReimbursementByUserIDPeriod(context.Background(), tt.userID, startDate, endDate)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectLen)
			}
		})
	}
}
