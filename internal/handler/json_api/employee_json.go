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

const fgwEmployeesStartUrl = "/api/fgw/employees"

type EmployeeHandlerJSON struct {
	roleService     service.RoleUseCase
	employeeService service.EmployeeUseCase
	wLogg           *wlogger.CustomWLogg
}

func NewEmployeeHandlerJSON(roleService service.RoleUseCase, employeeService service.EmployeeUseCase, wLogg *wlogger.CustomWLogg) *EmployeeHandlerJSON {
	return &EmployeeHandlerJSON{roleService: roleService, employeeService: employeeService, wLogg: wLogg}
}

func (e *EmployeeHandlerJSON) ServeJSONRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwEmployeesStartUrl, e.JSONAll)
	mux.HandleFunc(fgwEmployeesStartUrl+"/find", e.JSONFindById)
	mux.HandleFunc(fgwEmployeesStartUrl+"/add", e.JSONAdd)
	mux.HandleFunc(fgwEmployeesStartUrl+"/update", e.JSONUpdate)
	mux.HandleFunc(fgwEmployeesStartUrl+"/delete", e.JSONDelete)
}

func (e *EmployeeHandlerJSON) JSONAll(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, e.wLogg) {
		return
	}

	employees, err := e.employeeService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7003, err)

		return
	}

	if employees == nil {
		employees = []*entity.Employee{}
	}

	roles, err := e.roleService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7003, err)

		return
	}

	data := entity.EmployeeList{Employees: employees, Roles: roles}

	handler.WriteJSON(w, data, e.wLogg)
}

func (e *EmployeeHandlerJSON) JSONFindById(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, e.wLogg) {
		return
	}

	idEmployee, err := convert.ParseStrToUUID(r.URL.Query().Get("idEmployee"), w, r, e.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(r.Context(), idEmployee, w, r, e.wLogg, e.employeeService) {
		return
	}

	employee, err := e.employeeService.FindById(r.Context(), idEmployee)
	if err != nil {
		handler.WriteNotFound(w, r, e.wLogg, msg.H7005, err)

		return
	}

	handler.WriteJSON(w, employee, e.wLogg)
}

func (e *EmployeeHandlerJSON) JSONAdd(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, e.wLogg) {
		return
	}

	var employee entity.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		handler.WriteBadRequest(w, r, e.wLogg, msg.H7004, err)

		return
	}

	if err := e.employeeService.Add(r.Context(), &employee); err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7003, err)

		return
	}

	handler.WriteJSON(w, employee, e.wLogg)
}

func (e *EmployeeHandlerJSON) JSONUpdate(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPut, e.wLogg) {
		return
	}

	idEmployee, err := convert.ParseStrToUUID(r.URL.Query().Get("idEmployee"), w, r, e.wLogg)
	if err != nil {
		return
	}

	var employee entity.Employee
	if err = json.NewDecoder(r.Body).Decode(&employee); err != nil {
		handler.WriteBadRequest(w, r, e.wLogg, msg.H7004, err)

		return
	}

	if !handler.EntityExistsByUUID(r.Context(), idEmployee, w, r, e.wLogg, e.employeeService) {
		return
	}

	if err = e.employeeService.Update(r.Context(), idEmployee, &employee); err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7003, err)

		return
	}
	handler.WriteJSON(w, map[string]string{"message": msg.I2005}, e.wLogg)
}

func (e *EmployeeHandlerJSON) JSONDelete(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodDelete, e.wLogg) {
		return
	}

	idEmployee, err := convert.ParseStrToUUID(r.URL.Query().Get("idEmployee"), w, r, e.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(r.Context(), idEmployee, w, r, e.wLogg, e.employeeService) {
		return
	}

	if err = e.employeeService.Delete(r.Context(), idEmployee); err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7003, err)

		return
	}

	handler.WriteJSON(w, map[string]string{"message": msg.I2004}, e.wLogg)
}
