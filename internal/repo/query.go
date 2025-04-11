package repo

// Role query
var (
	FGWRoleAllQuery      = "exec dbo.fgw_role_all;"
	FGWRoleFindByIdQuery = "exec dbo.fgw_role_find_by_id ?;"
	FGWRoleAddQuery      = "exec dbo.fgw_role_add ?, ?, ?;"
	FGWRoleUpdateQuery   = "exec dbo.fgw_role_update ?, ?, ?;"
	FGWRoleDeleteQuery   = "exec dbo.fgw_role_delete_by_id ?;"
)
