package entity

import "github.com/google/uuid"

type EmployeeList struct {
	Employees []*Employee
	Roles     []*Role
}

type Employee struct {
	IdEmployee    uuid.UUID   `json:"idEmployee"`
	ServiceNumber int         `json:"serviceNumber"`
	FirstName     string      `json:"firstName"`
	LastName      string      `json:"lastName"`
	Patronymic    string      `json:"patronymic"`
	Passwd        string      `json:"passwd"`
	RoleId        uuid.UUID   `json:"roleId"` // Roles номер роли.
	AuditRecord   AuditRecord `json:"auditRecord"`
	IsEditing     bool        `json:"isEditing"`
}
