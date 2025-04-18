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

func (e *EmployeeHandlerJSON) JSONAll(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodGet, e.wLogg) {
		return
	}

	employees, err := e.employeeService.All(request.Context())
	if err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}

	if employees == nil {
		employees = []*entity.Employee{}
	}

	roles, err := e.roleService.All(request.Context())
	if err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}

	data := entity.EmployeeList{Employees: employees, Roles: roles}

	handler.WriteJSON(writer, data, e.wLogg)
}

func (e *EmployeeHandlerJSON) JSONFindById(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodGet, e.wLogg) {
		return
	}

	idEmployee, err := handler.ParseStrToUUID(request.URL.Query().Get("idEmployee"), writer, request, e.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExists(request.Context(), idEmployee, writer, request, e.wLogg, e.employeeService) {
		return
	}

	employee, err := e.employeeService.FindById(request.Context(), idEmployee)
	if err != nil {
		e.wLogg.LogHttpE(http.StatusNotFound, request.Method, request.URL.Path, msg.H7005, err)
		http.Error(writer, msg.H7005, http.StatusNotFound)

		return
	}

	handler.WriteJSON(writer, employee, e.wLogg)
}

func (e *EmployeeHandlerJSON) JSONAdd(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, e.wLogg) {
		return
	}

	var employee entity.Employee
	if err := json.NewDecoder(request.Body).Decode(&employee); err != nil {
		e.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7004, err)
		http.Error(writer, msg.H7004, http.StatusBadRequest)

		return
	}

	if err := e.employeeService.Add(request.Context(), &employee); err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}

	handler.WriteJSON(writer, employee, e.wLogg)
}

func (e *EmployeeHandlerJSON) JSONUpdate(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPut, e.wLogg) {
		return
	}

	idEmployee, err := handler.ParseStrToUUID(request.URL.Query().Get("idEmployee"), writer, request, e.wLogg)
	if err != nil {
		return
	}

	var employee entity.Employee
	if err = json.NewDecoder(request.Body).Decode(&employee); err != nil {
		e.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7004, err)
		http.Error(writer, msg.H7004, http.StatusBadRequest)

		return
	}

	if !handler.EntityExists(request.Context(), idEmployee, writer, request, e.wLogg, e.employeeService) {
		return
	}

	if err = e.employeeService.Update(request.Context(), idEmployee, &employee); err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}
	handler.WriteJSON(writer, map[string]string{"message": msg.I2005}, e.wLogg)
}

func (e *EmployeeHandlerJSON) JSONDelete(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodDelete, e.wLogg) {
		return
	}

	idEmployee, err := handler.ParseStrToUUID(request.URL.Query().Get("idEmployee"), writer, request, e.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExists(request.Context(), idEmployee, writer, request, e.wLogg, e.employeeService) {
		return
	}

	if err = e.employeeService.Delete(request.Context(), idEmployee); err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}

	handler.WriteJSON(writer, map[string]string{"message": msg.I2004}, e.wLogg)
}
