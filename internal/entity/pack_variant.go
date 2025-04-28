package entity

import "github.com/google/uuid"

type PackVariantList struct {
	PackVariants []*PackVariant
	Catalogs     []*Catalog
}

type PackVariant struct {
	IdPackVariant     int       `json:"idPackVariant"`
	ProdId            int       `json:"prodId" `
	Article           string    `json:"article" `
	PackName          string    `json:"packName" `
	Color             int       `json:"color" `
	GL                int       `json:"gl" ` // GL петля мёбиуса (переработка) 70-79
	QuantityRows      int       `json:"quantityRows" `
	QuantityPerRows   int       `json:"quantityPerRows" `
	Weight            int       `json:"weight" `
	Depth             int       `json:"depth" `
	Width             int       `json:"width" `
	Height            int       `json:"height" `
	IsFood            bool      `json:"isFood" `
	IsAfraidMoisture  bool      `json:"isAfraidMoisture" `
	IsAfraidSun       bool      `json:"isAfraidSun" `
	IsEAC             bool      `json:"isEAC" `
	IsAccountingBatch bool      `json:"isAccountingBatch" `
	MethodShip        bool      `json:"methodShip" `
	ShelfLifeMonths   int       `json:"shelfLifeMonths" `
	BathFurnace       int       `json:"bathFurnace" `
	MachineLine       int       `json:"machineLine" `
	IsManufactured    bool      `json:"isManufactured" `
	CurrentDateBatch  string    `json:"currentDateBatch" `
	NumberingBatch    int       `json:"numberingBatch" `
	IsArchive         bool      `json:"isArchive" `
	OwnerUser         uuid.UUID `json:"ownerUser" `
	OwnerUserDateTime string    `json:"ownerUserDateTime" `
	LastUser          uuid.UUID `json:"lastUser" `
	LastUserDateTime  string    `json:"lastUserDateTime" `
	IsEditing         bool      `json:"isEditing"`
}
