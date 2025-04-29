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

type PackVariantUseCase interface {
	All(ctx context.Context) ([]*entity.PackVariant, error)
	FindById(ctx context.Context, idPackVariant int) (*entity.PackVariant, error)
	Add(ctx context.Context, packVariant *entity.PackVariant) error
	Update(ctx context.Context, idPackVariant int, packVariant *entity.PackVariant) error
	Delete(ctx context.Context, idPackVariant int) error
	ExistsByID(ctx context.Context, idPackVariant int) (bool, error)
}

func (p *PackVariantService) All(ctx context.Context) ([]*entity.PackVariant, error) {
	packVariants, err := p.packVariantRepo.All(ctx)
	if err != nil {
		p.wLogg.LogE(msg.E3004, err)

		return nil, err
	}

	return packVariants, nil
}

func (p *PackVariantService) FindById(ctx context.Context, idPackVariant int) (*entity.PackVariant, error) {
	packVariant, err := p.packVariantRepo.FindById(ctx, idPackVariant)
	if err != nil {
		p.wLogg.LogE(msg.E3005, err)

		return nil, err
	}

	return packVariant, nil
}

func (p *PackVariantService) Add(ctx context.Context, packVariant *entity.PackVariant) error {
	if err := p.validateStruct.Struct(packVariant); err != nil {
		p.wLogg.LogW(msg.W1001, err)

		return err
	}

	if err := p.packVariantRepo.Add(ctx, packVariant); err != nil {
		p.wLogg.LogW(msg.E3006, err)

		return err
	}

	return nil
}

func (p *PackVariantService) Update(ctx context.Context, idPackVariant int, packVariant *entity.PackVariant) error {
	if err := p.validateStruct.Struct(packVariant); err != nil {
		p.wLogg.LogW(msg.W1001, err)

		return err
	}

	if err := p.packVariantRepo.Update(ctx, idPackVariant, packVariant); err != nil {
		p.wLogg.LogW(msg.E3007, err)

		return err
	}

	return nil
}

func (p *PackVariantService) Delete(ctx context.Context, idPackVariant int) error {
	if err := p.packVariantRepo.Delete(ctx, idPackVariant); err != nil {
		p.wLogg.LogW(msg.E3008, err)

		return err
	}

	return nil
}

func (p *PackVariantService) ExistsByID(ctx context.Context, idPackVariant int) (bool, error) {
	return p.packVariantRepo.ExistsByID(ctx, idPackVariant)
}
