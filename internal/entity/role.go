package entity

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Role роль сотрудника.
type Role struct {
	IdRole    uuid.UUID `json_api:"idRole"`
	Number    int       `json_api:"number" validate:"required"`
	Name      string    `json_api:"name" validate:"required,max=55"`
	IsEditing bool      `json_api:"isEditing"` // IsEditing флаг для редактирования поля.
}

var validate = validator.New()

// ValidateRole валидация полей структуры Role.
func ValidateRole(role *Role) error {
	return validate.Struct(role)
}
