package entity

import "github.com/google/uuid"

// AuditRecord - отслеживает изменения записей пользователями.
type AuditRecord struct {
	OwnerUser         uuid.UUID `json:"ownerUser"`
	OwnerUserDateTime string    `json:"ownerUserDateTime"`
	LastUser          uuid.UUID `json:"lastUser"`
	LastUserDateTime  string    `json:"lastUserDateTime"`
}
