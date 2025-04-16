package json_api

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type RoleHandlerJSON struct {
	roleService service.RoleUseCase
	wLogg       *wlogger.CustomWLogg
}

func NewRoleHandlerJson(roleService service.RoleUseCase, wLogg *wlogger.CustomWLogg) *RoleHandlerJSON {
	return &RoleHandlerJSON{roleService: roleService, wLogg: wLogg}
}

func (r *RoleHandlerJSON) ServeJSONRouters(mux *http.ServeMux) {
	mux.HandleFunc("/api/roles", r.RoleHandlerJSONAll)
	mux.HandleFunc("/api/roles/find", r.RoleHandlerJSONFindById)
	mux.HandleFunc("/api/roles/add", r.RoleHandlerJSONAdd)
	mux.HandleFunc("/api/roles/update", r.RoleHandlerJSONUpdate)
	mux.HandleFunc("/api/roles/delete", r.RoleHandlerJSONDelete)
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

	writeJSON(writer, roles, r.wLogg)
}

func (r *RoleHandlerJSON) RoleHandlerJSONFindById(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodGet, r.wLogg) {
		return
	}

	idRoleStr := request.URL.Query().Get("idRole")
	idRole, err := uuid.Parse(idRoleStr)
	if err != nil {
		r.wLogg.LogHttpE(http.StatusBadRequest, http.MethodGet, request.URL.Path, msg.H7004, err)
		http.Error(writer, msg.H7004, http.StatusBadRequest)

		return
	}

	role, err := r.roleService.FindById(request.Context(), idRole)
	if err != nil {
		r.wLogg.LogHttpE(http.StatusNotFound, request.Method, request.URL.Path, msg.H7005, err)
		http.Error(writer, msg.H7005, http.StatusNotFound)

		return
	}

	writeJSON(writer, role, r.wLogg)
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

	writeJSON(writer, role, r.wLogg)
}

func (r *RoleHandlerJSON) RoleHandlerJSONUpdate(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPut, r.wLogg) {
		return
	}

	idRoleStr := request.URL.Query().Get("idRole")
	idRole, err := uuid.Parse(idRoleStr)
	if err != nil {
		r.wLogg.LogHttpE(http.StatusBadRequest, http.MethodGet, request.URL.Path, msg.H7004, err)
		http.Error(writer, msg.H7004, http.StatusBadRequest)

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
}

func (r *RoleHandlerJSON) RoleHandlerJSONDelete(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodDelete, r.wLogg) {
		return
	}

	idRoleStr := request.URL.Query().Get("idRole")
	idRole, err := uuid.Parse(idRoleStr)
	if err != nil {
		r.wLogg.LogHttpE(http.StatusBadRequest, http.MethodGet, request.URL.Path, msg.H7004, err)
		http.Error(writer, msg.H7004, http.StatusBadRequest)

		return
	}

	if err = r.roleService.Delete(request.Context(), idRole); err != nil {
		r.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}

	writeJSON(writer, map[string]string{"message": "Удалено"}, r.wLogg)
}

func writeJSON(writer http.ResponseWriter, obj interface{}, wLogg *wlogger.CustomWLogg) {
	writer.Header().Set("Content-Type", "application/json_api; charset=UTF-8")
	if err := json.NewEncoder(writer).Encode(obj); err != nil {
		wLogg.LogE(msg.E3105, err)
		http.Error(writer, msg.E3105, http.StatusInternalServerError)

		return
	}
}
