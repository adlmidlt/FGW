package entity

import "github.com/google/uuid"

type CatalogList struct {
	Catalogs           []*Catalog
	Handbooks          []*Handbook
	SelectedHandbookId int
}

type Catalog struct {
	IdCatalog             int       `json:"idCatalog"`
	ParentId              int       `json:"parentId"`
	HandbookId            int       `json:"handbookId"`
	RecordIndex           int       `json:"recordIndex" `
	Name                  string    `json:"name" `
	Comment               string    `json:"comment"`
	HandbookValueInt1     int       `json:"handbookValueInt1"`
	HandbookValueInt2     int       `json:"handbookValueInt2" `
	HandbookValueDecimal1 float64   `json:"HandbookValueDecimal1" `
	HandbookValueDecimal2 float64   `json:"HandbookValueDecimal2" `
	HandbookValueBool1    bool      `json:"HandbookValueBool1" `
	HandbookValueBool2    bool      `json:"HandbookValueBool2" `
	IsArchive             bool      `json:"isArchive"`
	OwnerUser             uuid.UUID `json:"ownerUser" `
	OwnerUserDateTime     string    `json:"ownerUserDateTime"`
	LastUser              uuid.UUID `json:"lastUser"`
	LastUserDateTime      string    `json:"lastUserDateTime"`
	IsEditing             bool      `json:"isEditing"`
}
