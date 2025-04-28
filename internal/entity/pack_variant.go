package entity

import "github.com/google/uuid"

type PackVariantList struct {
	PackVariants []*PackVariant
}

type PackVariant struct {
	IdPackVariant     int       `json:"idPackVariant"`
	ProdId            int       `json:"prodId" validate:"required"`
	Article           string    `json:"article" validate:"required,max=5"`
	PackName          string    `json:"packName" validate:"required,max=255"`
	Color             int       `json:"color" validate:"required"`
	GL                int       `json:"gl" validate:"required"` // GL петля мёбиуса (переработка) 70-79
	QuantityRows      int       `json:"quantityRows" validate:"required"`
	QuantityPerRows   int       `json:"quantityPerRows" validate:"required"`
	Weight            int       `json:"weight" validate:"required"`
	Depth             int       `json:"depth" validate:"required"`
	Width             int       `json:"width" validate:"required"`
	Height            int       `json:"height" validate:"required"`
	IsFood            bool      `json:"isFood" validate:"required"`
	IsAfraidMoisture  bool      `json:"isAfraidMoisture" validate:"required"`
	IsAfraidSun       bool      `json:"isAfraidSun" validate:"required"`
	IsEAC             bool      `json:"isEAC" validate:"required"`
	IsAccountingBatch bool      `json:"isAccountingBatch" validate:"required"`
	MethodShip        bool      `json:"methodShip" validate:"required"`
	ShelfLifeMonths   int       `json:"shelfLifeMonths" validate:"required"`
	BathFurnace       int       `json:"bathFurnace" validate:"required"`
	MachineLine       int       `json:"machineLine" validate:"required"`
	IsManufactured    bool      `json:"isManufactured" validate:"required"`
	CurrentDateBatch  string    `json:"currentDateBatch" validate:"required"`
	NumberingBatch    int       `json:"numberingBatch" validate:"required"`
	IsArchive         bool      `json:"isArchive" validate:"required"`
	OwnerUserId       uuid.UUID `json:"ownerUserId" validate:"required"`
	OwnerUserDataTime string    `json:"ownerUserDataTime" validate:"required"`
	LastUser          uuid.UUID `json:"lastUser" validate:"required"`
	LastUserDateTime  string    `json:"lastUserDateTime" validate:"required"`
}
