package repository_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vnnyx/employee-management/internal/auditlog/entity"
	"github.com/vnnyx/employee-management/internal/auditlog/repository"
)

func TestInsertNewAuditLog(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewAuditLogRepository(mock)
	now := time.Now().Format(time.RFC3339)

	tests := []struct {
		name      string
		setupMock func()
		input     entity.AuditLog
		expectErr bool
	}{
		{
			name: "success - audit log inserted",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO audit_logs").
					WithArgs(
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
					).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("audit-id-1"))
			},
			input: entity.AuditLog{
				ID:        "audit-id-1",
				TableName: "users",
				RecordID:  "user-1",
				Action:    "UPDATE",
				ChangedBy: "admin",
				IPAddress: "127.0.0.1",
				RequestID: "req-123",
				OldData:   json.RawMessage(`{}`),
				NewData:   json.RawMessage(`{"name":"New Name"}`),
				CreatedAt: now,
			},
			expectErr: false,
		},
		{
			name: "error - returned ID is empty",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO audit_logs").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(""))
			},
			input: entity.AuditLog{
				ID:        "audit-id-2",
				OldData:   json.RawMessage(`{}`),
				NewData:   json.RawMessage(`{}`),
				CreatedAt: now,
			},
			expectErr: true,
		},
		{
			name: "error - query fails",
			setupMock: func() {
				mock.ExpectQuery("INSERT INTO audit_logs").
					WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(),
						pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(assert.AnError)
			},
			input: entity.AuditLog{
				ID:        "audit-id-3",
				OldData:   json.RawMessage(`{}`),
				NewData:   json.RawMessage(`{}`),
				CreatedAt: now,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := repo.InsertNewAuditLog(context.Background(), tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
