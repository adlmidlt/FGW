package json_api

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/convert"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"encoding/json"
	"net/http"
)

const fgwRolesStartUrl = "/api/fgw/roles"

type RoleHandlerJSON struct {
	roleService service.RoleUseCase
	wLogg       *wlogger.CustomWLogg
}

func NewRoleHandlerJSON(roleService service.RoleUseCase, wLogg *wlogger.CustomWLogg) *RoleHandlerJSON {
	return &RoleHandlerJSON{roleService: roleService, wLogg: wLogg}
}

func (r *RoleHandlerJSON) ServeJSONRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwRolesStartUrl, r.JSONAll)
	mux.HandleFunc(fgwRolesStartUrl+"/find", r.JSONFindById)
	mux.HandleFunc(fgwRolesStartUrl+"/add", r.JSONAdd)
	mux.HandleFunc(fgwRolesStartUrl+"/update", r.JSONUpdate)
	mux.HandleFunc(fgwRolesStartUrl+"/delete", r.JSONDelete)
}

func (r *RoleHandlerJSON) JSONAll(writer http.ResponseWriter, request *http.Request) {
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

	handler.WriteJSON(writer, roles, r.wLogg)
}

func (r *RoleHandlerJSON) JSONFindById(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodGet, r.wLogg) {
		return
	}

	idRole, err := convert.ParseStrToUUID(request.URL.Query().Get("idRole"), writer, request, r.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(request.Context(), idRole, writer, request, r.wLogg, r.roleService) {
		return
	}

	role, err := r.roleService.FindById(request.Context(), idRole)
	if err != nil {
		handler.WriteNotFound(writer, request, r.wLogg, msg.H7005, err)

		return
	}

	handler.WriteJSON(writer, role, r.wLogg)
}

func (r *RoleHandlerJSON) JSONAdd(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, r.wLogg) {
		return
	}

	var role entity.Role
	if err := json.NewDecoder(request.Body).Decode(&role); err != nil {
		handler.WriteBadRequest(writer, request, r.wLogg, msg.H7004, err)

		return
	}

	if err := r.roleService.Add(request.Context(), &role); err != nil {
		handler.WriteServerError(writer, request, r.wLogg, msg.H7003, err)

		return
	}

	handler.WriteJSON(writer, role, r.wLogg)
}

func (r *RoleHandlerJSON) JSONUpdate(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPut, r.wLogg) {
		return
	}

	idRole, err := convert.ParseStrToUUID(request.URL.Query().Get("idRole"), writer, request, r.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(request.Context(), idRole, writer, request, r.wLogg, r.roleService) {
		return
	}

	var role entity.Role
	if err = json.NewDecoder(request.Body).Decode(&role); err != nil {
		handler.WriteBadRequest(writer, request, r.wLogg, msg.H7004, err)

		return
	}

	if err = r.roleService.Update(request.Context(), idRole, &role); err != nil {
		handler.WriteServerError(writer, request, r.wLogg, msg.H7003, err)

		return
	}
	handler.WriteJSON(writer, map[string]string{"message": msg.I2005}, r.wLogg)
}

func (r *RoleHandlerJSON) JSONDelete(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodDelete, r.wLogg) {
		return
	}

	idRole, err := convert.ParseStrToUUID(request.URL.Query().Get("idRole"), writer, request, r.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(request.Context(), idRole, writer, request, r.wLogg, r.roleService) {
		return
	}

	if err = r.roleService.Delete(request.Context(), idRole); err != nil {
		handler.WriteServerError(writer, request, r.wLogg, msg.H7003, err)

		return
	}

	handler.WriteJSON(writer, map[string]string{"message": msg.I2004}, r.wLogg)
}
