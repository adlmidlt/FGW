package service

import (
	"FGW/internal/entity"
	"FGW/internal/repo"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type EmployeeService struct {
	employeeRepo   repo.EmployeeRepository
	wLogg          *wlogger.CustomWLogg
	validateStruct *validator.Validate
}

func NewEmployeeService(employeeRepo repo.EmployeeRepository, wLogg *wlogger.CustomWLogg, validateStruct *validator.Validate) *EmployeeService {
	return &EmployeeService{employeeRepo: employeeRepo, wLogg: wLogg, validateStruct: validateStruct}
}

type EmployeeUseCase interface {
	All(ctx context.Context) ([]*entity.Employee, error)
	FindById(ctx context.Context, idEmployee uuid.UUID) (*entity.Employee, error)
	Add(ctx context.Context, employee *entity.Employee) error
	Update(ctx context.Context, idEmployee uuid.UUID, employee *entity.Employee) error
	Delete(ctx context.Context, idEmployee uuid.UUID) error
	ExistsByUUID(ctx context.Context, idEmployee uuid.UUID) (bool, error)
}

func (e *EmployeeService) All(ctx context.Context) ([]*entity.Employee, error) {
	employees, err := e.employeeRepo.All(ctx)
	if err != nil {
		e.wLogg.LogE(msg.E3004, err)

		return nil, err
	}

	return employees, nil
}

func (e *EmployeeService) FindById(ctx context.Context, idEmployee uuid.UUID) (*entity.Employee, error) {
	employee, err := e.employeeRepo.FindById(ctx, idEmployee)
	if err != nil {
		e.wLogg.LogE(msg.E3005, err)

		return nil, err
	}

	return employee, nil
}

func (e *EmployeeService) Add(ctx context.Context, employee *entity.Employee) error {
	if err := e.validateStruct.Struct(employee); err != nil {
		e.wLogg.LogW(msg.W1001, err)

		return err
	}

	if employee.IdEmployee == uuid.Nil {
		employee.IdEmployee = uuid.New()
	}

	if err := e.employeeRepo.Add(ctx, employee); err != nil {
		e.wLogg.LogE(msg.E3006, err)

		return err
	}

	return nil
}

func (e *EmployeeService) Update(ctx context.Context, idEmployee uuid.UUID, employee *entity.Employee) error {
	if err := e.employeeRepo.Update(ctx, idEmployee, employee); err != nil {
		e.wLogg.LogE(msg.E3007, err)

		return err
	}

	return nil
}

func (e *EmployeeService) Delete(ctx context.Context, idEmployee uuid.UUID) error {
	if err := e.employeeRepo.Delete(ctx, idEmployee); err != nil {
		e.wLogg.LogE(msg.E3008, err)

		return err
	}

	return nil
}

func (e *EmployeeService) ExistsByUUID(ctx context.Context, idEmployee uuid.UUID) (bool, error) {
	return e.employeeRepo.ExistsByUUID(ctx, idEmployee)
}
