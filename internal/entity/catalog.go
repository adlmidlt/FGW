package entity

import "github.com/google/uuid"

type CatalogList struct {
	Catalogs  []*Catalog
	Handbooks []*Handbook
}

type Catalog struct {
	IdCatalog             int       `json:"idCatalog"`
	ParentId              int       `json:"parentId" validate:"required"`
	HandbookId            int       `json:"handbookId" validate:"required"`
	RecordIndex           int       `json:"recordIndex" validate:"required"`
	Name                  string    `json:"name" validate:"required,max=255"`
	Comment               string    `json:"comment" validate:"required,max=5000"`
	HandbookValueInt1     int       `json:"handbookValueInt1" validate:"required"`
	HandbookValueInt2     int       `json:"handbookValueInt2" validate:"required"`
	HandbookValueDecimal1 float64   `json:"HandbookValueDecimal1" validate:"required"`
	HandbookValueDecimal2 float64   `json:"HandbookValueDecimal2" validate:"required"`
	HandbookValueBool1    bool      `json:"HandbookValueBool1" validate:"required"`
	HandbookValueBool2    bool      `json:"HandbookValueBool2" validate:"required"`
	IsArchive             bool      `json:"isArchive" validate:"required"`
	OwnerUser             uuid.UUID `json:"ownerUser" validate:"required"`
	OwnerUserDateTime     string    `json:"ownerUserDateTime" validate:"required"`
	LastUser              uuid.UUID `json:"lastUser" validate:"required"`
	LastUserDateTime      string    `json:"lastUserDateTime" validate:"required"`
	IsEditing             bool      `json:"isEditing"`
}
