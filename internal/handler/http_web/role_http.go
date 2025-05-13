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
	"strings"
	"time"
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
	mux.HandleFunc(fgwRolesStartUrl, r.All)
	mux.HandleFunc(fgwRolesStartUrl+"/update", r.Update)
	mux.HandleFunc(fgwRolesStartUrl+"/delete", r.Delete)
	mux.HandleFunc(fgwRolesStartUrl+"/add", r.Add)
}

func (r *RoleHandlerHTTP) All(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodGet, r.wLogg) {
		return
	}

	roles, err := r.roleService.All(request.Context())
	if err != nil {
		handler.WriteServerError(writer, request, r.wLogg, msg.H7003, err)

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

func (r *RoleHandlerHTTP) Update(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		r.processUpdateFormRole(writer, request)
	case http.MethodGet:
		r.renderUpdateFormRole(writer, request)
	default:
		handler.WriteMethodNotAllowed(writer, request, r.wLogg, msg.H7002, nil)
	}
}

func (r *RoleHandlerHTTP) Delete(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, r.wLogg) {
		return
	}

	idRole, err := convert.ParseStrToUUID(request.FormValue(paramIdRole), writer, request, r.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(request.Context(), idRole, writer, request, r.wLogg, r.roleService) {
		return
	}

	if err = r.roleService.Delete(request.Context(), idRole); err != nil {
		handler.WriteServerError(writer, request, r.wLogg, msg.H7011, err)

		return
	}
	http.Redirect(writer, request, fgwRolesStartUrl, http.StatusSeeOther)
}

func (r *RoleHandlerHTTP) Add(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, r.wLogg) {
		return
	}

	errors := make(map[string]string)

	number := convert.ConvStrToInt(request.FormValue("number"))
	if number <= 0 {
		errors["number"] = msg.J1001
	}

	name := strings.TrimSpace(request.FormValue("name"))
	if len(name) > 55 {
		errors["name"] = msg.J1002
	}

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid.
	ownerUser := convert.ParseUUIDUnsafe(request.FormValue("ownerUser"))
	if ownerUser == uuid.Nil {
		ownerUser = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	}

	ownerUserDateTime := request.FormValue("ownerUserDateTime")
	if ownerUserDateTime == "" {
		ownerUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid.
	lastUser := convert.ParseUUIDUnsafe(request.FormValue("lastUser"))
	if lastUser == uuid.Nil {
		lastUser = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	}

	lastUserDateTime := request.FormValue("lastUserDateTime")
	if lastUserDateTime == "" {
		lastUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	if handler.SendErrorsJSON(writer, errors, r.wLogg) {
		return
	}

	role := &entity.Role{
		Number: number,
		Name:   name,
		AuditRecord: entity.AuditRecord{
			OwnerUser:         ownerUser,
			OwnerUserDateTime: ownerUserDateTime,
			LastUser:          lastUser,
			LastUserDateTime:  lastUserDateTime,
		},
	}

	if err := r.roleService.Add(request.Context(), role); err != nil {
		handler.WriteServerError(writer, request, r.wLogg, msg.H7012, err)

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
		handler.WriteBadRequest(writer, request, r.wLogg, msg.H7008, err)

		return
	}

	errors := make(map[string]string)

	number := convert.ConvStrToInt(request.FormValue("number"))
	if number <= 0 {
		errors["number"] = msg.J1001
	}

	name := strings.TrimSpace(request.FormValue("name"))
	if len(name) > 55 {
		errors["name"] = msg.J1002
	}

	idRole, err := convert.ParseStrToUUID(request.FormValue(paramIdRole), writer, request, r.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(request.Context(), idRole, writer, request, r.wLogg, r.roleService) {
		return
	}

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid при изменении записи.
	lastUser := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	lastUserDateTime := time.Now().Format("2006-01-02 15:04:05")

	if handler.SendErrorsJSON(writer, errors, r.wLogg) {
		return
	}

	role := &entity.Role{
		IdRole: idRole,
		Number: number,
		Name:   name,
		AuditRecord: entity.AuditRecord{
			LastUser:         lastUser,
			LastUserDateTime: lastUserDateTime,
		},
	}

	if err = r.roleService.Update(request.Context(), idRole, role); err != nil {
		handler.WriteServerError(writer, request, r.wLogg, msg.H7009, err)

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
