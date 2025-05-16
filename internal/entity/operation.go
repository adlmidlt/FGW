package entity

type OperationList struct {
	Operations []*Operation `json:"operations"`
	Catalogs   []*Catalog   `json:"catalogs"`
}

type Operation struct {
	IdOperation       int         `json:"idOperation"`
	TypeOperation     int         `json:"typeOperation"`
	CreateDate        string      `json:"dateOperation"`
	CreateByEmployee  int         `json:"createByEmployee"`
	DateOrder         string      `json:"dateOrder"`
	ClosedByEmployee  int         `json:"closedByEmployee"`
	CodeAccountingObj int         `json:"codeAccountingObj"`
	Appoint           int         `json:"appoint"`
	AuditRecord       AuditRecord `json:"auditRecord"`
	IsEditing         bool        `json:"isEditing"`
}
