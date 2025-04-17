package json_api

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"encoding/json"
	"net/http"
)

type RoleHandlerJSON struct {
	roleService service.RoleUseCase
	wLogg       *wlogger.CustomWLogg
}

func NewRoleHandlerJSON(roleService service.RoleUseCase, wLogg *wlogger.CustomWLogg) *RoleHandlerJSON {
	return &RoleHandlerJSON{roleService: roleService, wLogg: wLogg}
}

func (r *RoleHandlerJSON) ServeJSONRouters(mux *http.ServeMux) {
	mux.HandleFunc("/api/fgw/roles", r.RoleHandlerJSONAll)
	mux.HandleFunc("/api/fgw/roles/find", r.RoleHandlerJSONFindById)
	mux.HandleFunc("/api/fgw/roles/add", r.RoleHandlerJSONAdd)
	mux.HandleFunc("/api/fgw/roles/update", r.RoleHandlerJSONUpdate)
	mux.HandleFunc("/api/fgw/roles/delete", r.RoleHandlerJSONDelete)
}

func (r *RoleHandlerJSON) RoleHandlerJSONAll(writer http.ResponseWriter, request *http.Request) {
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

	handler.WriteJSON(writer, roles, r.wLogg)
}

func (r *RoleHandlerJSON) RoleHandlerJSONFindById(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodGet, r.wLogg) {
		return
	}

	idRole, err := handler.ParseStrToUUID(request.URL.Query().Get("idRole"), writer, request, r.wLogg)
	if err != nil {
		return
	}

	_, err = r.roleService.Exists(request.Context(), idRole)
	if err != nil {
		r.wLogg.LogHttpW(http.StatusNotFound, request.Method, request.URL.Path, msg.H7005, err)
		http.Error(writer, msg.H7005, http.StatusNotFound)
		handler.WriteJSON(writer, map[string]string{"message": msg.W1002}, r.wLogg)

		return
	}

	role, err := r.roleService.FindById(request.Context(), idRole)
	if err != nil {
		r.wLogg.LogHttpE(http.StatusNotFound, request.Method, request.URL.Path, msg.H7005, err)
		http.Error(writer, msg.H7005, http.StatusNotFound)

		return
	}

	handler.WriteJSON(writer, role, r.wLogg)
}

func (r *RoleHandlerJSON) RoleHandlerJSONAdd(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, r.wLogg) {
		return
	}

	var role entity.Role
	if err := json.NewDecoder(request.Body).Decode(&role); err != nil {
		r.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7004, err)
		http.Error(writer, msg.H7004, http.StatusBadRequest)

		return
	}

	if err := r.roleService.Add(request.Context(), &role); err != nil {
		r.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}

	handler.WriteJSON(writer, role, r.wLogg)
}

func (r *RoleHandlerJSON) RoleHandlerJSONUpdate(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPut, r.wLogg) {
		return
	}

	idRole, err := handler.ParseStrToUUID(request.URL.Query().Get("idRole"), writer, request, r.wLogg)
	if err != nil {
		return
	}

	_, err = r.roleService.Exists(request.Context(), idRole)
	if err != nil {
		r.wLogg.LogHttpW(http.StatusNotFound, request.Method, request.URL.Path, msg.H7005, err)
		http.Error(writer, msg.H7005, http.StatusNotFound)
		handler.WriteJSON(writer, map[string]string{"message": msg.W1002}, r.wLogg)

		return
	}

	var role entity.Role
	if err = json.NewDecoder(request.Body).Decode(&role); err != nil {
		r.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7004, err)
		http.Error(writer, msg.H7004, http.StatusBadRequest)

		return
	}

	if err = r.roleService.Update(request.Context(), idRole, &role); err != nil {
		r.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}
	handler.WriteJSON(writer, map[string]string{"message": msg.I2005}, r.wLogg)
}

func (r *RoleHandlerJSON) RoleHandlerJSONDelete(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodDelete, r.wLogg) {
		return
	}

	idRole, err := handler.ParseStrToUUID(request.URL.Query().Get("idRole"), writer, request, r.wLogg)
	if err != nil {
		return
	}

	_, err = r.roleService.Exists(request.Context(), idRole)
	if err != nil {
		r.wLogg.LogHttpW(http.StatusNotFound, request.Method, request.URL.Path, msg.H7005, err)
		http.Error(writer, msg.H7005, http.StatusNotFound)
		handler.WriteJSON(writer, map[string]string{"message": msg.W1002}, r.wLogg)

		return
	}

	if err = r.roleService.Delete(request.Context(), idRole); err != nil {
		r.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}

	handler.WriteJSON(writer, map[string]string{"message": msg.I2004}, r.wLogg)
}
