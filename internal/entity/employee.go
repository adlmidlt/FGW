package entity

import "github.com/google/uuid"

type EmployeeList struct {
	Employees []*Employee
	Roles     []*Role
}

type Employee struct {
	IdEmployee    uuid.UUID `json:"idEmployee"`
	ServiceNumber int       `json:"serviceNumber" validate:"required"`
	FirstName     string    `json:"firstName" validate:"required,max=50"`
	LastName      string    `json:"lastName" validate:"required,max=50"`
	Patronymic    string    `json:"patronymic" validate:"required,max=50"`
	Passwd        string    `json:"passwd" validate:"required,max=255"`
	RoleId        uuid.UUID `json:"roleId" validate:"required"` // Roles номер роли.
	IsEditing     bool      `json:"isEditing"`
}
