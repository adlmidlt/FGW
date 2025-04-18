package service

import (
	"FGW/internal/entity"
	"FGW/internal/repo"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
	"github.com/go-playground/validator/v10"
)

type PackVariantService struct {
	packVariantRepo repo.PackVariantRepository
	wLogg           *wlogger.CustomWLogg
	validateStruct  *validator.Validate
}

func NewPackVariantService(packVariantRepo repo.PackVariantRepository, wLogg *wlogger.CustomWLogg, validateStruct *validator.Validate) *PackVariantService {
	return &PackVariantService{packVariantRepo: packVariantRepo, wLogg: wLogg, validateStruct: validateStruct}
}

type PackVariantUserCase interface {
	All(ctx context.Context) ([]*entity.PackVariant, error)
}

func (p *PackVariantService) All(ctx context.Context) ([]*entity.PackVariant, error) {
	packVariants, err := p.packVariantRepo.All(ctx)
	if err != nil {
		p.wLogg.LogE(msg.E3004, err)

		return nil, err
	}

	return packVariants, nil
}
