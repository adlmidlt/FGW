package entity

type HandbookList struct {
	Handbooks []*Handbook
}

type Handbook struct {
	IdHandbook int    `json:"idHandbook"`
	Name       string `json:"name" validate:"required"`
	IsEditing  bool   `json:"isEditing"`
}
