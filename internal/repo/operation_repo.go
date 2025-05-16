package repo

import (
	"FGW/internal/entity"
	"FGW/pkg/db"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
	"database/sql"
)

type OperationRepo struct {
	mssql *sql.DB
	wLogg *wlogger.CustomWLogg
}

func NewOperationRepo(mssql *sql.DB, wLogg *wlogger.CustomWLogg) *OperationRepo {
	return &OperationRepo{mssql: mssql, wLogg: wLogg}
}

type OperationRepository interface {
	All(ctx context.Context) ([]*entity.Operation, error)
	FindById(ctx context.Context, idOperation int) (*entity.Operation, error)
	Add(ctx context.Context, operation *entity.Operation) error
	Update(ctx context.Context, idOperation int, operation *entity.Operation) error
	Delete(ctx context.Context, idOperation int) error
	ExistsByID(ctx context.Context, idOperation int) (bool, error)
}

func (o *OperationRepo) All(ctx context.Context) ([]*entity.Operation, error) {
	rows, err := o.mssql.QueryContext(ctx, FGWOperationAllQuery)
	if err != nil {
		o.wLogg.LogE(msg.E3000, err)

		return nil, err
	}
	defer db.CloseRows(rows)

	var operations []*entity.Operation
	for rows.Next() {
		var operation entity.Operation
		if err = rows.Scan(
			&operation.TypeOperation,
			&operation.CreateDate,
			&operation.CreateByEmployee,
			&operation.DateOrder,
			&operation.ClosedByEmployee,
			&operation.CodeAccountingObj,
			&operation.Appoint,
			&operation.AuditRecord.OwnerUser,
			&operation.AuditRecord.OwnerUserDateTime,
			&operation.AuditRecord.LastUser,
			&operation.AuditRecord.LastUserDateTime,
		); err != nil {
			o.wLogg.LogE(msg.E3001, err)

			return nil, err
		}
		operations = append(operations, &operation)
	}

	if len(operations) == 0 {
		o.wLogg.LogW(msg.W1000, nil)
	}

	if err = rows.Err(); err != nil {
		o.wLogg.LogE(msg.E3002, err)

		return nil, err
	}

	return operations, nil
}

func (o *OperationRepo) FindById(ctx context.Context, idOperation int) (*entity.Operation, error) {
	var operation entity.Operation
	if err := o.mssql.QueryRowContext(ctx, FGWOperationFindByIdQuery, idOperation).Scan(
		&operation.IdOperation,
		&operation.TypeOperation,
		&operation.CreateDate,
		&operation.CreateByEmployee,
		&operation.DateOrder,
		&operation.ClosedByEmployee,
		&operation.CodeAccountingObj,
		&operation.Appoint,
		&operation.AuditRecord.OwnerUser,
		&operation.AuditRecord.OwnerUserDateTime,
		&operation.AuditRecord.LastUser,
		&operation.AuditRecord.LastUserDateTime,
	); err != nil {
		o.wLogg.LogE(msg.E3000, err)

		return nil, err
	}

	return &operation, nil
}

func (o *OperationRepo) Add(ctx context.Context, operation *entity.Operation) error {
	if _, err := o.mssql.ExecContext(ctx, FGWOperationAddQuery,
		operation.TypeOperation,
		operation.CreateDate,
		operation.CreateByEmployee,
		operation.DateOrder,
		operation.ClosedByEmployee,
		operation.CodeAccountingObj,
		operation.Appoint,
		operation.AuditRecord.OwnerUser,
		operation.AuditRecord.OwnerUserDateTime,
		operation.AuditRecord.LastUser,
		operation.AuditRecord.LastUserDateTime,
	); err != nil {
		o.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (o *OperationRepo) Update(ctx context.Context, idOperation int, operation *entity.Operation) error {
	if _, err := o.mssql.ExecContext(ctx, FGWOperationUpdateQuery,
		operation.TypeOperation,
		operation.CreateDate,
		operation.CreateByEmployee,
		operation.DateOrder,
		operation.ClosedByEmployee,
		operation.CodeAccountingObj,
		operation.Appoint,
		operation.AuditRecord.LastUser,
		operation.AuditRecord.LastUserDateTime,
	); err != nil {
		o.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (o *OperationRepo) Delete(ctx context.Context, idOperation int) error {
	if _, err := o.mssql.ExecContext(ctx, FGWOperationDeleteQuery, idOperation); err != nil {
		o.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (o *OperationRepo) ExistsByID(ctx context.Context, idOperation int) (bool, error) {
	var exists bool
	row := o.mssql.QueryRowContext(ctx, FGWOperationExistQuery, idOperation)
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
