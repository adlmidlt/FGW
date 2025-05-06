package entity

import (
	"github.com/google/uuid"
)

type RoleList struct {
	Roles []*Role
}

// Role роль сотрудника.
type Role struct {
	IdRole      uuid.UUID   `json_api:"idRole"`
	Number      int         `json_api:"number"`
	Name        string      `json_api:"name"`
	AuditRecord AuditRecord `json:"auditRecord"`
	IsEditing   bool        `json_api:"isEditing"` // IsEditing флаг для редактирования поля.
}
