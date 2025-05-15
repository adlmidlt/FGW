package entity

type HandbookList struct {
	Handbooks []*Handbook
}

type Handbook struct {
	IdHandbook  int         `json:"idHandbook"`
	Name        string      `json:"name"`
	AuditRecord AuditRecord `json:"auditRecord"`
	IsEditing   bool        `json:"isEditing"`
}
