package service

import (
	"FGW/internal/entity"
	"FGW/internal/repo"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
	"github.com/go-playground/validator/v10"
)

type HandbookService struct {
	handbookRepo   repo.HandbookRepository
	wLogg          *wlogger.CustomWLogg
	validateStruct *validator.Validate
}

func NewHandbookService(handbookRepo repo.HandbookRepository, wLogg *wlogger.CustomWLogg, validateStruct *validator.Validate) *HandbookService {
	return &HandbookService{handbookRepo: handbookRepo, wLogg: wLogg, validateStruct: validateStruct}
}

type HandbookUseCase interface {
	All(ctx context.Context) ([]*entity.Handbook, error)
	FindById(ctx context.Context, idHandbook int) (*entity.Handbook, error)
	Add(ctx context.Context, handbook *entity.Handbook) error
	Update(ctx context.Context, idHandbook int, handbook *entity.Handbook) error
	Delete(ctx context.Context, idHandbook int) error
	ExistsByID(ctx context.Context, idHandbook int) (bool, error)
	AddZeroObj(ctx context.Context) error
}

func (h *HandbookService) All(ctx context.Context) ([]*entity.Handbook, error) {
	handbooks, err := h.handbookRepo.All(ctx)
	if err != nil {
		h.wLogg.LogE(msg.E3004, err)

		return nil, err
	}

	return handbooks, nil
}

func (h *HandbookService) FindById(ctx context.Context, idHandbook int) (*entity.Handbook, error) {
	handbook, err := h.handbookRepo.FindById(ctx, idHandbook)
	if err != nil {
		h.wLogg.LogE(msg.E3005, err)

		return nil, err
	}

	return handbook, nil
}

func (h *HandbookService) Add(ctx context.Context, handbook *entity.Handbook) error {
	if err := h.validateStruct.Struct(handbook); err != nil {
		h.wLogg.LogW(msg.W1001, err)

		return err
	}

	if err := h.handbookRepo.Add(ctx, handbook); err != nil {
		h.wLogg.LogE(msg.E3006, err)

		return err
	}

	return nil
}

func (h *HandbookService) Update(ctx context.Context, idHandbook int, handbook *entity.Handbook) error {
	if err := h.validateStruct.Struct(handbook); err != nil {
		h.wLogg.LogW(msg.W1001, err)

		return err
	}

	if err := h.handbookRepo.Update(ctx, idHandbook, handbook); err != nil {
		h.wLogg.LogE(msg.E3007, err)

		return err
	}

	return nil
}

func (h *HandbookService) Delete(ctx context.Context, idHandbook int) error {
	if err := h.handbookRepo.Delete(ctx, idHandbook); err != nil {
		h.wLogg.LogE(msg.E3008, err)

		return err
	}

	return nil
}

func (h *HandbookService) ExistsByID(ctx context.Context, idHandbook int) (bool, error) {
	return h.handbookRepo.ExistsByID(ctx, idHandbook)
}

func (h *HandbookService) AddZeroObj(ctx context.Context) error {
	return h.handbookRepo.AddZeroObj(ctx)
}
