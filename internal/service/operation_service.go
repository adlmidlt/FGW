package service

import (
	"FGW/internal/entity"
	"FGW/internal/repo"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
)

type OperationService struct {
	operationRepo repo.OperationRepository
	wLogg         *wlogger.CustomWLogg
}

func NewOperationService(operationRepo repo.OperationRepository, wLogg *wlogger.CustomWLogg) *OperationService {
	return &OperationService{operationRepo: operationRepo, wLogg: wLogg}
}

type OperationUseCase interface {
	All(ctx context.Context) ([]*entity.Operation, error)
	FindById(ctx context.Context, idOperation int) (*entity.Operation, error)
	Add(ctx context.Context, operation *entity.Operation) error
	Update(ctx context.Context, idOperation int, operation *entity.Operation) error
	Delete(ctx context.Context, idOperation int) error
	ExistsByID(ctx context.Context, idOperation int) (bool, error)
}

func (o *OperationService) All(ctx context.Context) ([]*entity.Operation, error) {
	operations, err := o.operationRepo.All(ctx)
	if err != nil {
		o.wLogg.LogE(msg.E3004, err)

		return nil, err
	}

	return operations, nil
}

func (o *OperationService) FindById(ctx context.Context, idOperation int) (*entity.Operation, error) {
	operation, err := o.operationRepo.FindById(ctx, idOperation)
	if err != nil {
		o.wLogg.LogE(msg.E3005, err)

		return nil, err
	}

	return operation, nil
}

func (o *OperationService) Add(ctx context.Context, operation *entity.Operation) error {
	if err := o.operationRepo.Add(ctx, operation); err != nil {
		o.wLogg.LogE(msg.E3006, err)

		return err
	}

	return nil
}

func (o *OperationService) Update(ctx context.Context, idOperation int, operation *entity.Operation) error {
	if err := o.operationRepo.Update(ctx, idOperation, operation); err != nil {
		o.wLogg.LogE(msg.E3007, err)

		return err
	}

	return nil
}

func (o *OperationService) Delete(ctx context.Context, idOperation int) error {
	if err := o.operationRepo.Delete(ctx, idOperation); err != nil {
		o.wLogg.LogE(msg.E3008, err)

		return err
	}

	return nil
}

func (o *OperationService) ExistsByID(ctx context.Context, idOperation int) (bool, error) {
	return o.operationRepo.ExistsByID(ctx, idOperation)
}
