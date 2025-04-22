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
	FindById(ctx context.Context, idCatalog int) (*entity.Catalog, error)
	Add(ctx context.Context, catalog *entity.Catalog) error
	Update(ctx context.Context, idCatalog int, catalog *entity.Catalog) error
	Delete(ctx context.Context, idCatalog int) error
	ExistsByID(ctx context.Context, idCatalog int) (bool, error)
}

func (c *CatalogService) All(ctx context.Context) ([]*entity.Catalog, error) {
	catalogs, err := c.catalogRepo.All(ctx)
	if err != nil {
		c.wLogg.LogE(msg.E3004, err)

		return nil, err
	}

	return catalogs, nil
}

func (c *CatalogService) FindById(ctx context.Context, idCatalog int) (*entity.Catalog, error) {
	catalog, err := c.catalogRepo.FindById(ctx, idCatalog)
	if err != nil {
		c.wLogg.LogE(msg.E3005, err)

		return nil, err
	}

	return catalog, nil
}

func (c *CatalogService) Add(ctx context.Context, catalog *entity.Catalog) error {
	if err := c.validateStruct.Struct(catalog); err != nil {
		c.wLogg.LogW(msg.W1001, err)

		return err
	}

	if err := c.catalogRepo.Add(ctx, catalog); err != nil {
		c.wLogg.LogE(msg.E3006, err)

		return err
	}

	return nil
}

func (c *CatalogService) Update(ctx context.Context, idCatalog int, catalog *entity.Catalog) error {
	if err := c.validateStruct.Struct(catalog); err != nil {
		c.wLogg.LogW(msg.W1001, err)

		return err
	}

	if err := c.catalogRepo.Update(ctx, idCatalog, catalog); err != nil {
		c.wLogg.LogE(msg.E3007, err)

		return err
	}

	return nil
}

func (c *CatalogService) Delete(ctx context.Context, idCatalog int) error {
	if err := c.catalogRepo.Delete(ctx, idCatalog); err != nil {
		c.wLogg.LogE(msg.E3008, err)

		return err
	}

	return nil
}

func (c *CatalogService) ExistsByID(ctx context.Context, idCatalog int) (bool, error) {
	return c.catalogRepo.ExistsByID(ctx, idCatalog)
}
