package repo

// Role query
var (
	FGWRoleAllQuery      = "exec dbo.fgw_role_all;"
	FGWRoleFindByIdQuery = "exec dbo.fgw_role_find_by_id ?;"
	FGWRoleAddQuery      = "exec dbo.fgw_role_add ?, ?, ?;"
	FGWRoleUpdateQuery   = "exec dbo.fgw_role_update ?, ?, ?;"
	FGWRoleDeleteQuery   = "exec dbo.fgw_role_delete_by_id ?;"
)

var (
	FGWEmployeeAllQuery      = "exec dbo.fgw_employee_all;"
	FGWEmployeeFindByIdQuery = "exec dbo.fgw_employee_find_by_id ?;"
	FGWEmployeeAddQuery      = "exec dbo.fgw_employee_add ?, ?, ?, ?, ?, ?, ?;"
	FGWEmployeeUpdate        = "exec dbo.fgw_employee_update ?, ?, ?, ?, ?, ?, ?"
	FGWEmployeeDelete        = "exec dbo.fgw_employee_delete ?"
)
