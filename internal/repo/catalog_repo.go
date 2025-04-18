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

type CatalogRepo struct {
	mssql *sql.DB
	wLogg *wlogger.CustomWLogg
}

func NewCatalog(mssql *sql.DB, wLogg *wlogger.CustomWLogg) *CatalogRepo {
	return &CatalogRepo{mssql: mssql, wLogg: wLogg}
}

type CatalogRepository interface {
	All(ctx context.Context) ([]*entity.Catalog, error)
}

func (c *CatalogRepo) All(ctx context.Context) ([]*entity.Catalog, error) {
	rows, err := c.mssql.QueryContext(ctx, FGWCatalogAllQuery)
	if err != nil {
		c.wLogg.LogE(msg.E3000, err)

		return nil, err
	}
	defer db.CloseRows(rows)

	var catalogs []*entity.Catalog
	for rows.Next() {
		var catalog entity.Catalog
		if err = rows.Scan(
			&catalog.IdCatalog,
			&catalog.ParentId,
			&catalog.HandbookId,
			&catalog.RecordIndex,
			&catalog.Name,
			&catalog.Comment,
			&catalog.HandbookValueInt1,
			&catalog.HandbookValueInt2,
			&catalog.HandbookValueDecimal1,
			&catalog.HandbookValueDecimal2,
			&catalog.HandbookValueBool1,
			&catalog.HandbookValueBool2,
			&catalog.IsArchive,
			&catalog.OwnerUser,
			&catalog.OwnerUserDateTime,
			&catalog.LastUser,
			&catalog.LastUserDateTime,
		); err != nil {
			c.wLogg.LogE(msg.E3001, err)

			return nil, err
		}

		catalog.Name, _ = convert.Win1251ToUTF8(catalog.Name)
		catalog.Comment, _ = convert.Win1251ToUTF8(catalog.Comment)

		catalogs = append(catalogs, &catalog)
	}

	if len(catalogs) == 0 {
		c.wLogg.LogW(msg.W1000, err)
	}

	if err = rows.Err(); err != nil {
		c.wLogg.LogE(msg.E3002, err)

		return nil, err
	}

	return catalogs, nil
}
