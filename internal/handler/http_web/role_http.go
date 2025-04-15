package http_web

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"html/template"
	"net/http"
)

const templateHtmlRoleList = "../web/html/role_list.html"

type RoleHandlerHTTP struct {
	roleService service.RoleUseCase
	wLogg       *wlogger.CustomWLogg
}

func NewRoleHandlerHTTP(roleService service.RoleUseCase, wLogg *wlogger.CustomWLogg) *RoleHandlerHTTP {
	return &RoleHandlerHTTP{roleService, wLogg}
}

func (r *RoleHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc("/roles/", r.RoleHandlerHTTPAll)
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
