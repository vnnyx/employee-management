package repository

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/auditlog"
	"github.com/vnnyx/employee-management/internal/auditlog/entity"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type auditLogRepo struct {
	db database.Queryer
}

func NewAuditLogRepository(db database.Queryer) auditlog.Repository {
	return &auditLogRepo{
		db: db,
	}
}

func (r *auditLogRepo) WithTx(tx database.DBTx) auditlog.Repository {
	return &auditLogRepo{
		db: tx,
	}
}

func (r *auditLogRepo) InsertNewAuditLog(ctx context.Context, auditLog entity.AuditLog) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"AuditLogRepository.InsertNewAuditLog()",
	)
	defer span.End()

	query, args, err := sqlx.Named(insertAuditLogQuery, auditLog)
	if err != nil {
		return errors.Wrap(err, "AuditLogRepository.InsertNewAuditLog().Named()")
	}
	query = database.Rebind(query)

	var returnedID string
	err = pgxscan.Get(ctx, r.db, &returnedID, query, args...)
	if err != nil {
		return errors.Wrap(err, "AuditLogRepository.InsertNewAuditLog().Get()")
	}

	if returnedID == "" {
		return errors.Wrap(errors.New("failed to insert audit log"), "AuditLogRepository.InsertNewAuditLog().Get()")
	}

	return nil
}
