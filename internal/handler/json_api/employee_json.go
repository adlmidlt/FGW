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

type EmployeeHandlerJSON struct {
	roleService     service.RoleUseCase
	employeeService service.EmployeeUseCase
	wLogg           *wlogger.CustomWLogg
}

func NewEmployeeHandlerJSON(roleService service.RoleUseCase, employeeService service.EmployeeUseCase, wLogg *wlogger.CustomWLogg) *EmployeeHandlerJSON {
	return &EmployeeHandlerJSON{roleService: roleService, employeeService: employeeService, wLogg: wLogg}
}

func (e *EmployeeHandlerJSON) ServeJSONRouters(mux *http.ServeMux) {
	mux.HandleFunc("/api/fgw/employees", e.EmployeeHandlerJSONAll)
	mux.HandleFunc("/api/fgw/employees/find", e.EmployeeHandlerJSONFindById)
	mux.HandleFunc("/api/fgw/employees/add", e.EmployeeHandlerJSONAdd)
	mux.HandleFunc("/api/fgw/employees/update", e.EmployeeHandlerJSONUpdate)
	mux.HandleFunc("/api/fgw/employees/delete", e.EmployeeHandlerJSONDelete)
}

func (e *EmployeeHandlerJSON) EmployeeHandlerJSONAll(writer http.ResponseWriter, request *http.Request) {
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

func (e *EmployeeHandlerJSON) EmployeeHandlerJSONFindById(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodGet, e.wLogg) {
		return
	}

	idEmployee, err := handler.ParseStrToUUID(request.URL.Query().Get("idEmployee"), writer, request, e.wLogg)
	if err != nil {
		return
	}

	_, err = e.employeeService.Exists(request.Context(), idEmployee)
	if err != nil {
		e.wLogg.LogHttpW(http.StatusNotFound, request.Method, request.URL.Path, msg.H7005, err)
		http.Error(writer, msg.H7005, http.StatusNotFound)
		handler.WriteJSON(writer, map[string]string{"message": msg.W1002}, e.wLogg)

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

func (e *EmployeeHandlerJSON) EmployeeHandlerJSONAdd(writer http.ResponseWriter, request *http.Request) {
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

func (e *EmployeeHandlerJSON) EmployeeHandlerJSONUpdate(writer http.ResponseWriter, request *http.Request) {
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

	_, err = e.employeeService.Exists(request.Context(), idEmployee)
	if err != nil {
		e.wLogg.LogHttpW(http.StatusNotFound, request.Method, request.URL.Path, msg.H7005, err)
		http.Error(writer, msg.H7005, http.StatusNotFound)
		handler.WriteJSON(writer, map[string]string{"message": msg.W1002}, e.wLogg)

		return
	}

	if err = e.employeeService.Update(request.Context(), idEmployee, &employee); err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}
	handler.WriteJSON(writer, map[string]string{"message": msg.I2005}, e.wLogg)
}

func (e *EmployeeHandlerJSON) EmployeeHandlerJSONDelete(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodDelete, e.wLogg) {
		return
	}

	idEmployee, err := handler.ParseStrToUUID(request.URL.Query().Get("idEmployee"), writer, request, e.wLogg)
	if err != nil {
		return
	}

	_, err = e.employeeService.Exists(request.Context(), idEmployee)
	if err != nil {
		e.wLogg.LogHttpW(http.StatusNotFound, request.Method, request.URL.Path, msg.H7005, err)
		http.Error(writer, msg.H7005, http.StatusNotFound)
		handler.WriteJSON(writer, map[string]string{"message": msg.W1002}, e.wLogg)

		return
	}

	if err = e.employeeService.Delete(request.Context(), idEmployee); err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7003, err)
		http.Error(writer, msg.H7003, http.StatusInternalServerError)

		return
	}

	handler.WriteJSON(writer, map[string]string{"message": msg.I2004}, e.wLogg)
}
