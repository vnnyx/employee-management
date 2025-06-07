package auditlog

import (
	"context"

	"github.com/vnnyx/employee-management/internal/auditlog/entity"
	"github.com/vnnyx/employee-management/pkg/database"
)

type Repository interface {
	WithTx(tx database.DBTx) Repository

	InsertNewAuditLog(ctx context.Context, auditLog entity.AuditLog) error
}
