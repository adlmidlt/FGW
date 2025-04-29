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

type PackVariantRepo struct {
	mssql *sql.DB
	wLogg *wlogger.CustomWLogg
}

func NewPackVariant(mssql *sql.DB, wLogg *wlogger.CustomWLogg) *PackVariantRepo {
	return &PackVariantRepo{mssql: mssql, wLogg: wLogg}
}

type PackVariantRepository interface {
	All(ctx context.Context) ([]*entity.PackVariant, error)
	FindById(ctx context.Context, idPackVariant int) (*entity.PackVariant, error)
	Add(ctx context.Context, packVariant *entity.PackVariant) error
	Update(ctx context.Context, idPackVariant int, packVariant *entity.PackVariant) error
	Delete(ctx context.Context, idPackVariant int) error
	ExistsByID(ctx context.Context, idPackVariant int) (bool, error)
}

func (p *PackVariantRepo) All(ctx context.Context) ([]*entity.PackVariant, error) {
	rows, err := p.mssql.QueryContext(ctx, FGWPackVariantAllQuery)
	if err != nil {
		p.wLogg.LogE(msg.E3000, err)

		return nil, err
	}
	defer db.CloseRows(rows)

	var packVariants []*entity.PackVariant
	for rows.Next() {
		var packVariant entity.PackVariant
		if err = rows.Scan(
			&packVariant.IdPackVariant,
			&packVariant.ProdId,
			&packVariant.Article,
			&packVariant.PackName,
			&packVariant.Color,
			&packVariant.GL,
			&packVariant.QuantityRows,
			&packVariant.QuantityPerRows,
			&packVariant.Weight,
			&packVariant.Depth,
			&packVariant.Width,
			&packVariant.Height,
			&packVariant.IsFood,
			&packVariant.IsAfraidMoisture,
			&packVariant.IsAfraidSun,
			&packVariant.IsEAC,
			&packVariant.IsAccountingBatch,
			&packVariant.MethodShip,
			&packVariant.ShelfLifeMonths,
			&packVariant.BathFurnace,
			&packVariant.MachineLine,
			&packVariant.IsManufactured,
			&packVariant.CurrentDateBatch,
			&packVariant.NumberingBatch,
			&packVariant.IsArchive,
			&packVariant.OwnerUser,
			&packVariant.OwnerUserDateTime,
			&packVariant.LastUser,
			&packVariant.LastUserDateTime,
		); err != nil {
			p.wLogg.LogE(msg.E3001, err)

			return nil, err
		}

		packVariant.PackName, _ = convert.Win1251ToUTF8(packVariant.PackName)

		packVariants = append(packVariants, &packVariant)
	}

	if len(packVariants) == 0 {
		p.wLogg.LogW(msg.W1000, nil)
	}

	if err = rows.Err(); err != nil {
		p.wLogg.LogE(msg.E3002, err)

		return nil, err
	}

	return packVariants, nil
}

func (p *PackVariantRepo) FindById(ctx context.Context, idPackVariant int) (*entity.PackVariant, error) {
	var packVariant entity.PackVariant
	if err := p.mssql.QueryRowContext(ctx, FGWPackVariantFindByIdQuery, idPackVariant).Scan(
		&packVariant.IdPackVariant,
		&packVariant.ProdId,
		&packVariant.Article,
		&packVariant.PackName,
		&packVariant.Color,
		&packVariant.GL,
		&packVariant.QuantityRows,
		&packVariant.QuantityPerRows,
		&packVariant.Weight,
		&packVariant.Depth,
		&packVariant.Width,
		&packVariant.Height,
		&packVariant.IsFood,
		&packVariant.IsAfraidMoisture,
		&packVariant.IsAfraidSun,
		&packVariant.IsEAC,
		&packVariant.IsAccountingBatch,
		&packVariant.MethodShip,
		&packVariant.ShelfLifeMonths,
		&packVariant.BathFurnace,
		&packVariant.MachineLine,
		&packVariant.IsManufactured,
		&packVariant.CurrentDateBatch,
		&packVariant.NumberingBatch,
		&packVariant.IsArchive,
		&packVariant.OwnerUser,
		&packVariant.OwnerUserDateTime,
		&packVariant.LastUser,
		&packVariant.LastUserDateTime,
	); err != nil {
		p.wLogg.LogE(msg.E3000, err)

		return nil, err
	}

	packVariant.PackName, _ = convert.Win1251ToUTF8(packVariant.PackName)

	return &packVariant, nil
}

func (p *PackVariantRepo) Add(ctx context.Context, packVariant *entity.PackVariant) error {
	if _, err := p.mssql.ExecContext(ctx, FGWPackVariantAddQuery,
		packVariant.ProdId,
		packVariant.Article,
		packVariant.PackName,
		packVariant.Color,
		packVariant.GL,
		packVariant.QuantityRows,
		packVariant.QuantityPerRows,
		packVariant.Weight,
		packVariant.Depth,
		packVariant.Width,
		packVariant.Height,
		packVariant.IsFood,
		packVariant.IsAfraidMoisture,
		packVariant.IsAfraidSun,
		packVariant.IsEAC,
		packVariant.IsAccountingBatch,
		packVariant.MethodShip,
		packVariant.ShelfLifeMonths,
		packVariant.BathFurnace,
		packVariant.MachineLine,
		packVariant.IsManufactured,
		packVariant.CurrentDateBatch,
		packVariant.NumberingBatch,
		packVariant.IsArchive,
		packVariant.OwnerUser,
		packVariant.OwnerUserDateTime,
		packVariant.LastUser,
		packVariant.LastUserDateTime,
	); err != nil {
		p.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (p *PackVariantRepo) Update(ctx context.Context, idPackVariant int, packVariant *entity.PackVariant) error {
	if _, err := p.mssql.ExecContext(ctx, FGWPackVariantUpdateQuery, idPackVariant,
		packVariant.ProdId,
		packVariant.Article,
		packVariant.PackName,
		packVariant.Color,
		packVariant.GL,
		packVariant.QuantityRows,
		packVariant.QuantityPerRows,
		packVariant.Weight,
		packVariant.Depth,
		packVariant.Width,
		packVariant.Height,
		packVariant.IsFood,
		packVariant.IsAfraidMoisture,
		packVariant.IsAfraidSun,
		packVariant.IsEAC,
		packVariant.IsAccountingBatch,
		packVariant.MethodShip,
		packVariant.ShelfLifeMonths,
		packVariant.BathFurnace,
		packVariant.MachineLine,
		packVariant.IsManufactured,
		packVariant.CurrentDateBatch,
		packVariant.NumberingBatch,
		packVariant.IsArchive,
		packVariant.OwnerUser,
		packVariant.OwnerUserDateTime,
		packVariant.LastUser,
		packVariant.LastUserDateTime,
	); err != nil {
		p.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (p *PackVariantRepo) Delete(ctx context.Context, idPackVariant int) error {
	if _, err := p.mssql.ExecContext(ctx, FGWPackVariantDeleteQuery, idPackVariant); err != nil {
		p.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (p *PackVariantRepo) ExistsByID(ctx context.Context, idPackVariant int) (bool, error) {
	var exists bool
	row := p.mssql.QueryRowContext(ctx, FGWPackVariantExistQuery, idPackVariant)
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
