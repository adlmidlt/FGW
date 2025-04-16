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
	"html/template"
	"net/http"
)

const templateHtmlRoleList = "../web/html/role_list.html"
const fgwRolesStartUrl = "/fgw/roles"

type RoleHandlerHTTP struct {
	roleService service.RoleUseCase
	wLogg       *wlogger.CustomWLogg
}

func NewRoleHandlerHTTP(roleService service.RoleUseCase, wLogg *wlogger.CustomWLogg) *RoleHandlerHTTP {
	return &RoleHandlerHTTP{roleService, wLogg}
}

func (r *RoleHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc("/fgw/roles", r.RoleHandlerHTTPAll)
	mux.HandleFunc("/fgw/roles/update", r.RoleHandlerHTTPUpdate)
	mux.HandleFunc("/fgw/roles/delete", r.RoleHandlerHTTPDelete)
	mux.HandleFunc("/fgw/roles/add", r.RoleHandlerHTTPAdd)
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

	if idStr := request.URL.Query().Get("idRole"); idStr != "" {
		if id, err := uuid.Parse(idStr); err == nil {
			for _, role := range roles {
				if role.IdRole == id {
					role.IsEditing = true
				}
			}
		}
	}

	if roles == nil {
		roles = []*entity.Role{}
	}

	data := entity.RoleList{Roles: roles}

	tmpl, err := template.ParseFiles(templateHtmlRoleList)
	if err != nil {
		r.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7006, err)
		http.Error(writer, msg.H7006, http.StatusInternalServerError)

		return
	}

	if err = tmpl.Execute(writer, data); err != nil {
		r.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7007, err)
		http.Error(writer, msg.H7007, http.StatusInternalServerError)

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

	idRoleStr := request.FormValue("idRole")
	idRole, err := uuid.Parse(idRoleStr)
	if err != nil {
		r.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7004, err)
		http.Error(writer, msg.H7004, http.StatusBadRequest)

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
	idRoleStr := request.URL.Query().Get("idRole")
	idRole, err := uuid.Parse(idRoleStr)
	if err != nil {
		r.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7004, err)
		http.Error(writer, msg.H7004, http.StatusBadRequest)

		return
	}
	http.Redirect(writer, request, fmt.Sprintf("%s?idRole=%s", fgwRolesStartUrl, idRole), http.StatusSeeOther)
}

func (r *RoleHandlerHTTP) processUpdateFormRole(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		r.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7008, err)
		http.Error(writer, msg.H7008, http.StatusBadRequest)

		return
	}

	idRoleStr := request.FormValue("idRole")
	idRole, err := uuid.Parse(idRoleStr)
	if err != nil {
		r.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7004, err)
		http.Error(writer, msg.H7004, http.StatusBadRequest)

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
