package repo

import (
	"FGW/internal/entity"
	"FGW/pkg/convert"
	"FGW/pkg/db"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
)

type RoleRepo struct {
	mssql *sql.DB
	wLogg *wlogger.CustomWLogg
}

func NewRoleRepo(mssql *sql.DB, wLogg *wlogger.CustomWLogg) *RoleRepo {
	return &RoleRepo{mssql: mssql, wLogg: wLogg}
}

type RoleRepository interface {
	All(ctx context.Context) ([]*entity.Role, error)
	FindById(ctx context.Context, idRole uuid.UUID) (*entity.Role, error)
	Add(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, idRole uuid.UUID, role *entity.Role) error
	Delete(ctx context.Context, idRole uuid.UUID) error
	Exists(ctx context.Context, idRole uuid.UUID) (bool, error)
}

func (r *RoleRepo) All(ctx context.Context) ([]*entity.Role, error) {
	rows, err := r.mssql.QueryContext(ctx, FGWRoleAllQuery)
	if err != nil {
		r.wLogg.LogE(msg.E3000, err)

		return nil, err
	}
	defer db.CloseRows(rows)

	var roles []*entity.Role
	for rows.Next() {
		var role entity.Role
		if err = rows.Scan(&role.IdRole, &role.Number, &role.Name); err != nil {
			r.wLogg.LogE(msg.E3001, err)

			return nil, err
		}
		role.Name, _ = convert.Win1251ToUTF8(role.Name)

		roles = append(roles, &role)
	}

	if len(roles) == 0 {
		r.wLogg.LogW(msg.W1000, nil)
	}

	if err = rows.Err(); err != nil {
		r.wLogg.LogE(msg.E3002, err)

		return nil, err
	}

	return roles, nil
}

func (r *RoleRepo) FindById(ctx context.Context, idRole uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	if err := r.mssql.QueryRowContext(ctx, FGWRoleFindByIdQuery, idRole).Scan(&role.IdRole, &role.Number, &role.Name); err != nil {
		r.wLogg.LogE(msg.E3000, err)

		return nil, err
	}
	role.Name, _ = convert.Win1251ToUTF8(role.Name)

	return &role, nil
}

func (r *RoleRepo) Add(ctx context.Context, role *entity.Role) error {
	if _, err := r.mssql.ExecContext(ctx, FGWRoleAddQuery, role.IdRole, role.Number, role.Name); err != nil {
		r.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (r *RoleRepo) Update(ctx context.Context, idRole uuid.UUID, role *entity.Role) error {
	if _, err := r.mssql.ExecContext(ctx, FGWRoleUpdateQuery, idRole, role.Number, role.Name); err != nil {
		r.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (r *RoleRepo) Delete(ctx context.Context, idRole uuid.UUID) error {
	if _, err := r.mssql.ExecContext(ctx, FGWRoleDeleteQuery, idRole); err != nil {
		r.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (r *RoleRepo) Exists(ctx context.Context, idRole uuid.UUID) (bool, error) {
	row := r.mssql.QueryRowContext(ctx, FGWRoleExistsQuery, idRole)

	var exists int
	err := row.Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	if err != nil {
		return false, nil
	}

	return true, nil
}
