package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vnnyx/employee-management/internal/payroll/entity"
	"github.com/vnnyx/employee-management/internal/payroll/repository"
)

func TestStoreNewPayroll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewPayrollRepository(mock)
	now := time.Now()

	tests := []struct {
		name      string
		setupMock func()
		input     entity.Payroll
		expectErr bool
	}{
		{
			name: "success",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO payrolls").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("payroll-1"))
			},
			input:     entity.Payroll{ID: "payroll-1", CreatedAt: now, UpdatedAt: now},
			expectErr: false,
		},
		{
			name: "error - empty returned ID",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO payrolls").
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(""))
			},
			input:     entity.Payroll{ID: "payroll-2"},
			expectErr: true,
		},
		{
			name: "error - query fails",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO payrolls").
					WillReturnError(errors.New("query error"))
			},
			input:     entity.Payroll{ID: "payroll-3"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.StoreNewPayroll(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStoreNewPayslips(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewPayrollRepository(mock)
	now := time.Now()

	tests := []struct {
		name      string
		setupMock func()
		input     []entity.Payslip
		expectErr bool
	}{
		{
			name: "success - multiple rows",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO payslips").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("ps-1").AddRow("ps-2"))
			},
			input: []entity.Payslip{
				{ID: "ps-1", CreatedAt: now, UpdatedAt: now},
				{ID: "ps-2", CreatedAt: now, UpdatedAt: now},
			},
			expectErr: false,
		},
		{
			name: "error - mismatched return count",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO payslips").
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("ps-1"))
			},
			input:     []entity.Payslip{{ID: "ps-1"}, {ID: "ps-2"}},
			expectErr: true,
		},
		{
			name: "error - query fails",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO payslips").
					WillReturnError(assert.AnError)
			},
			input:     []entity.Payslip{{ID: "ps-1"}},
			expectErr: true,
		},
		{
			name:      "empty input",
			setupMock: func() {}, // No DB call expected
			input:     []entity.Payslip{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.StoreNewPayslips(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStoreNewPayrollSummary(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewPayrollRepository(mock)
	now := time.Now()

	tests := []struct {
		name      string
		setupMock func()
		input     entity.PayrollSummary
		expectErr bool
	}{
		{
			name: "success",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO payroll_summaries").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("summary-1"))
			},
			input:     entity.PayrollSummary{ID: "summary-1", CreatedAt: now, UpdatedAt: now},
			expectErr: false,
		},
		{
			name: "error - empty returned ID",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO payroll_summaries").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(""))
			},
			input:     entity.PayrollSummary{ID: "summary-2"},
			expectErr: true,
		},
		{
			name: "error - query fails",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO payroll_summaries").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnError(assert.AnError)
			},
			input:     entity.PayrollSummary{ID: "summary-3"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.StoreNewPayrollSummary(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
