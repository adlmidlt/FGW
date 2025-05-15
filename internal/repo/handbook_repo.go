package repo

import (
	"FGW/internal/entity"
	"FGW/pkg/convert"
	"FGW/pkg/db"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
	"database/sql"
)

type HandbookRepo struct {
	mssql *sql.DB
	wLogg *wlogger.CustomWLogg
}

func NewHandbookRepo(mssql *sql.DB, wLogg *wlogger.CustomWLogg) *HandbookRepo {
	return &HandbookRepo{mssql: mssql, wLogg: wLogg}
}

type HandbookRepository interface {
	All(ctx context.Context) ([]*entity.Handbook, error)
	FindById(ctx context.Context, idHandbook int) (*entity.Handbook, error)
	Add(ctx context.Context, handbook *entity.Handbook) error
	Update(ctx context.Context, idHandbook int, handbook *entity.Handbook) error
	Delete(ctx context.Context, idHandbook int) error
	ExistsByID(ctx context.Context, idHandbook int) (bool, error)
	AddZeroObj(ctx context.Context) error
}

func (h *HandbookRepo) All(ctx context.Context) ([]*entity.Handbook, error) {
	rows, err := h.mssql.QueryContext(ctx, FGWHandbookAllQuery)
	if err != nil {
		h.wLogg.LogE(msg.E3000, err)

		return nil, err
	}
	defer db.CloseRows(rows)

	var handbooks []*entity.Handbook
	for rows.Next() {
		var handbook entity.Handbook
		if err = rows.Scan(
			&handbook.IdHandbook,
			&handbook.Name,
			&handbook.AuditRecord.OwnerUser,
			&handbook.AuditRecord.OwnerUserDateTime,
			&handbook.AuditRecord.LastUser,
			&handbook.AuditRecord.LastUserDateTime,
		); err != nil {
			h.wLogg.LogE(msg.E3001, err)

			return nil, err
		}
		handbook.Name, _ = convert.Win1251ToUTF8(handbook.Name)

		handbooks = append(handbooks, &handbook)
	}

	if len(handbooks) == 0 {
		h.wLogg.LogW(msg.W1000, nil)
	}

	if err = rows.Err(); err != nil {
		h.wLogg.LogE(msg.E3002, err)

		return nil, err
	}

	return handbooks, nil
}

func (h *HandbookRepo) FindById(ctx context.Context, idHandbook int) (*entity.Handbook, error) {
	var handbook entity.Handbook
	if err := h.mssql.QueryRowContext(ctx, FGWHandbookFindByIdQuery, idHandbook).Scan(
		&handbook.IdHandbook,
		&handbook.Name,
		&handbook.AuditRecord.OwnerUser,
		&handbook.AuditRecord.OwnerUserDateTime,
		&handbook.AuditRecord.LastUser,
		&handbook.AuditRecord.LastUserDateTime,
	); err != nil {
		h.wLogg.LogE(msg.E3000, err)

		return nil, err
	}
	handbook.Name, _ = convert.Win1251ToUTF8(handbook.Name)

	return &handbook, nil
}

func (h *HandbookRepo) Add(ctx context.Context, handbook *entity.Handbook) error {
	if _, err := h.mssql.ExecContext(ctx, FGWHandbookAddQuery,
		handbook.Name,
		handbook.AuditRecord.OwnerUser,
		handbook.AuditRecord.OwnerUserDateTime,
		handbook.AuditRecord.LastUser,
		handbook.AuditRecord.LastUserDateTime,
	); err != nil {
		h.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (h *HandbookRepo) Update(ctx context.Context, idHandbook int, handbook *entity.Handbook) error {
	if _, err := h.mssql.ExecContext(ctx, FGWHandbookUpdateQuery, idHandbook,
		handbook.Name,
		handbook.AuditRecord.LastUser,
		handbook.AuditRecord.LastUserDateTime,
	); err != nil {
		h.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (h *HandbookRepo) Delete(ctx context.Context, idHandbook int) error {
	if _, err := h.mssql.ExecContext(ctx, FGWHandbookDeleteQuery, idHandbook); err != nil {
		h.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (h *HandbookRepo) ExistsByID(ctx context.Context, idHandbook int) (bool, error) {
	var exists bool
	row := h.mssql.QueryRowContext(ctx, FGWHandbookExistsQuery, idHandbook)
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (h *HandbookRepo) AddZeroObj(ctx context.Context) error {
	if _, err := h.mssql.ExecContext(ctx, FGWHandbookAddZeroObjQuery); err != nil {
		h.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}
