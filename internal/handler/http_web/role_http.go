package http_web

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/convert"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

const (
	templateHtmlRoleList = "../web/html/role_list.html"
	fgwRolesStartUrl     = "/fgw/roles"
	paramIdRole          = "idRole"
)

type RoleHandlerHTTP struct {
	roleService service.RoleUseCase
	wLogg       *wlogger.CustomWLogg
}

func NewRoleHandlerHTTP(roleService service.RoleUseCase, wLogg *wlogger.CustomWLogg) *RoleHandlerHTTP {
	return &RoleHandlerHTTP{roleService, wLogg}
}

func (r *RoleHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwRolesStartUrl, r.RoleHandlerHTTPAll)
	mux.HandleFunc(fgwRolesStartUrl+"/update", r.RoleHandlerHTTPUpdate)
	mux.HandleFunc(fgwRolesStartUrl+"/delete", r.RoleHandlerHTTPDelete)
	mux.HandleFunc(fgwRolesStartUrl+"/add", r.RoleHandlerHTTPAdd)
}

func (r *RoleHandlerHTTP) RoleHandlerHTTPAll(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodGet, r.wLogg) {
		return
	}

	roles, err := r.roleService.All(request.Context())
	if err != nil {
		r.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}

	if roles == nil {
		roles = []*entity.Role{}
	}

	data := entity.RoleList{Roles: roles}

	if idStr := request.URL.Query().Get(paramIdRole); idStr != "" {
		r.markEditingRole(idStr, roles)
	}

	tmpl, ok := handler.ParseTemplateHTML(templateHtmlRoleList, writer, request, r.wLogg)
	if !ok {
		return
	}

	if !handler.ExecuteTemplate(tmpl, data, writer, request, r.wLogg) {
		return
	}
}

func (r *RoleHandlerHTTP) RoleHandlerHTTPUpdate(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		r.processUpdateFormRole(writer, request)
	case http.MethodGet:
		r.renderUpdateFormRole(writer, request)
	default:
		http.Error(writer, msg.H7002, http.StatusMethodNotAllowed)
	}
}

func (r *RoleHandlerHTTP) RoleHandlerHTTPDelete(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, r.wLogg) {
		return
	}

	idRole, err := handler.ParseStrToUUID(request.FormValue(paramIdRole), writer, request, r.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExists(request.Context(), idRole, writer, request, r.wLogg, r.roleService) {
		return
	}

	if err = r.roleService.Delete(request.Context(), idRole); err != nil {
		r.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7011, err)
		http.Error(writer, msg.H7011, http.StatusInternalServerError)

		return
	}
	http.Redirect(writer, request, fgwRolesStartUrl, http.StatusSeeOther)
}

func (r *RoleHandlerHTTP) RoleHandlerHTTPAdd(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, r.wLogg) {
		return
	}

	role := &entity.Role{
		Number: convert.ConvStrToInt(request.FormValue("number")),
		Name:   request.FormValue("name"),
	}

	if err := r.roleService.Add(request.Context(), role); err != nil {
		r.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7012, err)
		http.Error(writer, msg.H7012, http.StatusInternalServerError)

		return
	}
	http.Redirect(writer, request, fgwRolesStartUrl, http.StatusSeeOther)
}

func (r *RoleHandlerHTTP) renderUpdateFormRole(writer http.ResponseWriter, request *http.Request) {
	idRoleStr := request.URL.Query().Get(paramIdRole)
	http.Redirect(writer, request, fmt.Sprintf("%s?%s=%s", fgwRolesStartUrl, paramIdRole, idRoleStr), http.StatusSeeOther)
}

func (r *RoleHandlerHTTP) processUpdateFormRole(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		r.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7008, err)
		http.Error(writer, msg.H7008, http.StatusBadRequest)

		return
	}

	idRole, err := handler.ParseStrToUUID(request.FormValue(paramIdRole), writer, request, r.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExists(request.Context(), idRole, writer, request, r.wLogg, r.roleService) {
		return
	}

	role := &entity.Role{
		IdRole: idRole,
		Number: convert.ConvStrToInt(request.FormValue("number")),
		Name:   request.FormValue("name"),
	}

	if err = r.roleService.Update(request.Context(), idRole, role); err != nil {
		r.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7009, err)
		http.Error(writer, msg.H7009, http.StatusInternalServerError)

		return
	}
	http.Redirect(writer, request, fgwRolesStartUrl, http.StatusSeeOther)
}

// markEditingRole помечает роль как редактируемую по её UUID в строковом формате.
func (r *RoleHandlerHTTP) markEditingRole(idStr string, roles []*entity.Role) {
	if id, err := uuid.Parse(idStr); err == nil {
		for _, role := range roles {
			if role.IdRole == id {
				role.IsEditing = true
			}
		}
	}
}
