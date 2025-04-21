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

type EmployeeRepo struct {
	mssql *sql.DB
	wLogg *wlogger.CustomWLogg
}

func NewEmployeeRepo(mssql *sql.DB, wLogg *wlogger.CustomWLogg) *EmployeeRepo {
	return &EmployeeRepo{mssql: mssql, wLogg: wLogg}
}

type EmployeeRepository interface {
	All(ctx context.Context) ([]*entity.Employee, error)
	FindById(ctx context.Context, idEmployee uuid.UUID) (*entity.Employee, error)
	Add(ctx context.Context, employee *entity.Employee) error
	Update(ctx context.Context, idEmployee uuid.UUID, employee *entity.Employee) error
	Delete(ctx context.Context, idEmployee uuid.UUID) error
	ExistsByUUID(ctx context.Context, idEmployee uuid.UUID) (bool, error)
}

func (e *EmployeeRepo) All(ctx context.Context) ([]*entity.Employee, error) {
	rows, err := e.mssql.QueryContext(ctx, FGWEmployeeAllQuery)
	if err != nil {
		e.wLogg.LogE(msg.E3000, err)

		return nil, err
	}
	defer db.CloseRows(rows)

	var employees []*entity.Employee
	for rows.Next() {
		var employee entity.Employee
		if err = rows.Scan(
			&employee.IdEmployee,
			&employee.ServiceNumber,
			&employee.FirstName,
			&employee.LastName,
			&employee.Patronymic,
			&employee.Passwd,
			&employee.RoleId,
		); err != nil {
			e.wLogg.LogE(msg.E3001, err)

			return nil, err
		}
		employee.FirstName, _ = convert.Win1251ToUTF8(employee.FirstName)
		employee.LastName, _ = convert.Win1251ToUTF8(employee.LastName)
		employee.Patronymic, _ = convert.Win1251ToUTF8(employee.Patronymic)

		employees = append(employees, &employee)
	}

	if len(employees) == 0 {
		e.wLogg.LogW(msg.W1000, nil)
	}

	if err = rows.Err(); err != nil {
		e.wLogg.LogE(msg.E3002, err)

		return nil, err
	}

	return employees, nil
}

func (e *EmployeeRepo) FindById(ctx context.Context, idEmployee uuid.UUID) (*entity.Employee, error) {
	var employee entity.Employee
	if err := e.mssql.QueryRowContext(ctx, FGWEmployeeFindByIdQuery, idEmployee).Scan(
		&employee.IdEmployee,
		&employee.ServiceNumber,
		&employee.FirstName,
		&employee.LastName,
		&employee.Patronymic,
		&employee.Passwd,
		&employee.RoleId,
	); err != nil {
		e.wLogg.LogE(msg.E3000, err)

		return nil, err
	}

	employee.FirstName, _ = convert.Win1251ToUTF8(employee.FirstName)
	employee.LastName, _ = convert.Win1251ToUTF8(employee.LastName)
	employee.Patronymic, _ = convert.Win1251ToUTF8(employee.Patronymic)

	return &employee, nil
}

func (e *EmployeeRepo) Add(ctx context.Context, employee *entity.Employee) error {
	if _, err := e.mssql.ExecContext(ctx, FGWEmployeeAddQuery, employee.IdEmployee,
		employee.ServiceNumber,
		employee.FirstName,
		employee.LastName,
		employee.Patronymic,
		employee.Passwd,
		employee.RoleId,
	); err != nil {
		e.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (e *EmployeeRepo) Update(ctx context.Context, idEmployee uuid.UUID, employee *entity.Employee) error {
	if _, err := e.mssql.ExecContext(ctx, FGWEmployeeUpdate, idEmployee, employee.ServiceNumber,
		employee.FirstName,
		employee.LastName,
		employee.Patronymic,
		employee.Passwd,
		employee.RoleId,
	); err != nil {
		e.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (e *EmployeeRepo) Delete(ctx context.Context, idEmployee uuid.UUID) error {
	if _, err := e.mssql.ExecContext(ctx, FGWEmployeeDelete, idEmployee); err != nil {
		e.wLogg.LogE(msg.E3000, err)

		return err
	}

	return nil
}

func (e *EmployeeRepo) ExistsByUUID(ctx context.Context, idEmployee uuid.UUID) (bool, error) {
	row := e.mssql.QueryRowContext(ctx, FGWEmployeeExistQuery, idEmployee)

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
