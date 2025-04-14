package service

import (
	"FGW/internal/entity"
	"FGW/internal/repo"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
	"github.com/google/uuid"
)

type RoleService struct {
	roleRepo repo.RoleRepository
	wLogg    *wlogger.CustomWLogg
}

func NewRoleService(roleRepo repo.RoleRepository, wLogg *wlogger.CustomWLogg) *RoleService {
	return &RoleService{roleRepo: roleRepo, wLogg: wLogg}
}

type RoleUseCase interface {
	All(ctx context.Context) ([]*entity.Role, error)
	FindById(ctx context.Context, idRole uuid.UUID) (*entity.Role, error)
	Add(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, idRole uuid.UUID, role *entity.Role) error
	Delete(ctx context.Context, idRole uuid.UUID) error
}

func (r *RoleService) All(ctx context.Context) ([]*entity.Role, error) {
	roles, err := r.roleRepo.All(ctx)
	if err != nil {
		r.wLogg.LogE(msg.E3004, err)

		return nil, err
	}

	return roles, nil
}

func (r *RoleService) FindById(ctx context.Context, idRole uuid.UUID) (*entity.Role, error) {
	role, err := r.roleRepo.FindById(ctx, idRole)
	if err != nil {
		r.wLogg.LogE(msg.E3005, err)

		return nil, err
	}

	return role, nil
}

func (r *RoleService) Add(ctx context.Context, role *entity.Role) error {
	if err := entity.ValidateRole(role); err != nil {
		r.wLogg.LogW(msg.W1001, err)

		return err
	}

	if role.IdRole == uuid.Nil {
		role.IdRole = uuid.New()

	}

	if err := r.roleRepo.Add(ctx, role); err != nil {
		r.wLogg.LogE(msg.E3006, err)

		return err
	}

	return nil
}

func (r *RoleService) Update(ctx context.Context, idRole uuid.UUID, role *entity.Role) error {
	if err := entity.ValidateRole(role); err != nil {
		r.wLogg.LogW(msg.W1001, err)

		return err
	}

	if err := r.roleRepo.Update(ctx, idRole, role); err != nil {
		r.wLogg.LogE(msg.E3007, err)

		return err
	}

	return nil
}

func (r *RoleService) Delete(ctx context.Context, idRole uuid.UUID) error {
	if err := r.roleRepo.Delete(ctx, idRole); err != nil {
		r.wLogg.LogE(msg.E3008, err)

		return err
	}

	return nil
}
