package service

import (
	"FGW/internal/entity"
	"FGW/internal/repo"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
	"github.com/go-playground/validator/v10"
)

type CatalogService struct {
	catalogRepo    repo.CatalogRepository
	wLogg          *wlogger.CustomWLogg
	validateStruct *validator.Validate
}

func NewCatalogService(catalogRepo repo.CatalogRepository, wLogg *wlogger.CustomWLogg, validateStruct *validator.Validate) *CatalogService {
	return &CatalogService{catalogRepo: catalogRepo, wLogg: wLogg, validateStruct: validateStruct}
}

type CatalogUseCase interface {
	All(ctx context.Context) ([]*entity.Catalog, error)
}

func (c *CatalogService) All(ctx context.Context) ([]*entity.Catalog, error) {
	catalogs, err := c.catalogRepo.All(ctx)
	if err != nil {
		c.wLogg.LogE(msg.E3004, err)

		return nil, err
	}

	return catalogs, nil
}
