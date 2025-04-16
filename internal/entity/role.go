package entity

import (
	"github.com/google/uuid"
)

type RoleList struct {
	Roles []*Role
}

// Role роль сотрудника.
type Role struct {
	IdRole    uuid.UUID `json_api:"idRole"`
	Number    int       `json_api:"number" validate:"required"`
	Name      string    `json_api:"name" validate:"required,max=55"`
	IsEditing bool      `json_api:"isEditing"` // IsEditing флаг для редактирования поля.
}
