package repo

// Role query
var (
	FGWRoleAllQuery      = "exec dbo.fgw_role_all;"
	FGWRoleFindByIdQuery = "exec dbo.fgw_role_find_by_id ?;"
	FGWRoleAddQuery      = "exec dbo.fgw_role_add ?, ?, ?, ?, ?, ?, ?;"
	FGWRoleUpdateQuery   = "exec dbo.fgw_role_update ?, ?, ?, ?, ?;"
	FGWRoleDeleteQuery   = "exec dbo.fgw_role_delete_by_id ?;"
	FGWRoleExistsQuery   = "exec dbo.fgw_role_exist ?;"
)

var (
	FGWEmployeeAllQuery      = "exec dbo.fgw_employee_all;"
	FGWEmployeeFindByIdQuery = "exec dbo.fgw_employee_find_by_id ?;"
	FGWEmployeeAddQuery      = "exec dbo.fgw_employee_add ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?;"
	FGWEmployeeUpdateQuery   = "exec dbo.fgw_employee_update ?, ?, ?, ?, ?, ?, ?, ?, ?;"
	FGWEmployeeDeleteQuery   = "exec dbo.fgw_employee_delete_by_id ?;"
	FGWEmployeeExistQuery    = "exec dbo.fgw_employee_exist ?;"
)

var (
	FGWCatalogAllQuery             = "exec dbo.fgw_catalog_all;"
	FGWCatalogFindByIdQuery        = "exec dbo.fgw_catalog_find_by_id ?;"
	FGWCatalogAddQuery             = "exec dbo.fgw_catalog_add ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?;"
	FGWCatalogUpdateQuery          = "exec dbo.fgw_catalog_update ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?;"
	FGWCatalogDeleteQuery          = "exec dbo.fgw_catalog_delete_by_id ?;"
	FGWCatalogExistQuery           = "exec dbo.fgw_catalog_exist ?;"
	FGWCatalogAllFindByNumberQuery = "exec dbo.fgw_catalog_all_find_by_number ?;"
)

var (
	FGWHandbookAllQuery      = "exec dbo.fgw_handbook_all;"
	FGWHandbookFindByIdQuery = "exec dbo.fgw_handbook_find_by_id ?;"
	FGWHandbookAddQuery      = "exec dbo.fgw_handbook_add ?;"
	FGWHandbookUpdateQuery   = "exec dbo.fgw_handbook_update ?, ?;"
	FGWHandbookDeleteQuery   = "exec dbo.fgw_handbook_delete_by_id ?;"
	FGWHandbookExistsQuery   = "exec dbo.fgw_handbook_exist ?;"
)

var (
	FGWPackVariantAllQuery      = "exec dbo.fgw_pack_variant_all;"
	FGWPackVariantFindByIdQuery = "exec dbo.fgw_pack_variant_find_by_id ?;"
	FGWPackVariantAddQuery      = "exec dbo.fgw_pack_variant_add ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?;"
	FGWPackVariantUpdateQuery   = "exec dbo.fgw_pack_variant_update ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?;"
	FGWPackVariantDeleteQuery   = "exec dbo.fgw_pack_variant_delete_by_id ?;"
	FGWPackVariantExistQuery    = "exec dbo.fgw_pack_variant_exist ?;"
)
