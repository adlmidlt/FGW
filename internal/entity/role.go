package entity

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Role роль сотрудника.
type Role struct {
	IdRole    uuid.UUID `json:"idRole" validate:"required"`
	Number    int       `json:"number" validate:"required,get>=1"`
	Name      string    `json:"name" validate:"required,max=55"`
	IsEditing bool      `json:"isEditing" validate:"required"`
}

var validate = validator.New()

// ValidateRole валидация полей структуры Role.
func ValidateRole(role *Role) error {
	return validate.Struct(role)
}
