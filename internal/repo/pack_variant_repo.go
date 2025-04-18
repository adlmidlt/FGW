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
}

func (p *PackVariantRepo) All(ctx context.Context) ([]*entity.PackVariant, error) {
	rows, err := p.mssql.QueryContext(ctx, FGWPackVariantAllQuery)
	if err != nil {
		p.wLogg.LogE(msg.E3000, err)

		return nil, err
	}
	defer db.CloseRows(rows)

	var packVariants []*entity.PackVariant
	if rows.Next() {
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
			&packVariant.OwnerUserId,
			&packVariant.OwnerUserDataTime,
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
