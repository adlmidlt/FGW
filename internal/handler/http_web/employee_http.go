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

const templateHtmlEmployeeList = "../web/html/employee_list.html"
const fgwEmployeesStartUrl = "/fgw/employees"

type EmployeeHandlerHTTP struct {
	roleService     service.RoleUseCase
	employeeService service.EmployeeUseCase
	wLogg           *wlogger.CustomWLogg
}

func NewEmployeeHandlerHTTP(roleService service.RoleUseCase, employeeService service.EmployeeUseCase, wLogg *wlogger.CustomWLogg) *EmployeeHandlerHTTP {
	return &EmployeeHandlerHTTP{roleService: roleService, employeeService: employeeService, wLogg: wLogg}
}

func (e *EmployeeHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc("/fgw/employees", e.EmployeeHandlerHTTPAll)
	mux.HandleFunc("/fgw/employees/update", e.EmployeeHandlerHTTPUpdate)
	mux.HandleFunc("/fgw/employees/delete", e.EmployeeHandlerHTTPDelete)
	mux.HandleFunc("/fgw/employees/add", e.EmployeeHandlerHTTPAdd)
}

func (e *EmployeeHandlerHTTP) EmployeeHandlerHTTPAll(writer http.ResponseWriter, request *http.Request) {
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

	if idEmployeeStr := request.URL.Query().Get("idEmployee"); idEmployeeStr != "" {
		if idEmployee, err := uuid.Parse(idEmployeeStr); err == nil {
			for _, employee := range employees {
				if employee.IdEmployee == idEmployee {
					employee.IsEditing = true
				}
			}
		}
	}

	tmpl, err := template.ParseFiles(templateHtmlEmployeeList)
	if err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7006, err)
		http.Error(writer, msg.H7006, http.StatusInternalServerError)

		return
	}

	if err = tmpl.Execute(writer, data); err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7007, err)
		http.Error(writer, msg.H7007, http.StatusInternalServerError)

		return
	}
}

func (e *EmployeeHandlerHTTP) EmployeeHandlerHTTPUpdate(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		e.processUpdateFormEmployee(writer, request)
	case http.MethodGet:
		e.renderUpdateFormEmployee(writer, request)
	default:
		http.Error(writer, msg.H7002, http.StatusMethodNotAllowed)
	}
}

func (e *EmployeeHandlerHTTP) EmployeeHandlerHTTPAdd(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, e.wLogg) {
		return
	}

	roleId, err := handler.ParseStrToUUID(request.FormValue("roleId"), writer, request, e.wLogg)
	if err != nil {
		return
	}

	employee := &entity.Employee{
		ServiceNumber: convert.ConvStrToInt(request.FormValue("serviceNumber")),
		FirstName:     request.FormValue("firstName"),
		LastName:      request.FormValue("lastName"),
		Patronymic:    request.FormValue("patronymic"),
		Passwd:        request.FormValue("passwd"),
		RoleId:        roleId,
	}

	if err = e.employeeService.Add(request.Context(), employee); err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7012, err)
		http.Error(writer, msg.H7012, http.StatusInternalServerError)

		return
	}
	http.Redirect(writer, request, fgwEmployeesStartUrl, http.StatusSeeOther)
}

func (e *EmployeeHandlerHTTP) EmployeeHandlerHTTPDelete(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, e.wLogg) {
		return
	}

	idEmployee, err := handler.ParseStrToUUID(request.FormValue("idEmployee"), writer, request, e.wLogg)
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
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7011, err)
		http.Error(writer, msg.H7011, http.StatusInternalServerError)

		return
	}
	http.Redirect(writer, request, fgwEmployeesStartUrl, http.StatusSeeOther)
}

func (e *EmployeeHandlerHTTP) renderUpdateFormEmployee(writer http.ResponseWriter, request *http.Request) {
	idEmployeeStr := request.URL.Query().Get("idEmployee")
	http.Redirect(writer, request, fmt.Sprintf("%s?idEmployee=%s", fgwEmployeesStartUrl, idEmployeeStr), http.StatusSeeOther)
}

func (e *EmployeeHandlerHTTP) processUpdateFormEmployee(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		e.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7008, err)
		http.Error(writer, msg.H7008, http.StatusBadRequest)

		return
	}

	idEmployee, err := handler.ParseStrToUUID(request.FormValue("idEmployee"), writer, request, e.wLogg)
	if err != nil {
		return
	}

	roleId, err := handler.ParseStrToUUID(request.FormValue("roleId"), writer, request, e.wLogg)
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

	employee := &entity.Employee{
		IdEmployee:    idEmployee,
		ServiceNumber: convert.ConvStrToInt(request.FormValue("serviceNumber")),
		FirstName:     request.FormValue("firstName"),
		LastName:      request.FormValue("lastName"),
		Patronymic:    request.FormValue("patronymic"),
		Passwd:        request.FormValue("passwd"),
		RoleId:        roleId,
	}

	if err = e.employeeService.Update(request.Context(), idEmployee, employee); err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, request.Method, request.URL.Path, msg.H7009, err)
		http.Error(writer, msg.H7009, http.StatusInternalServerError)

		return
	}
	http.Redirect(writer, request, fgwEmployeesStartUrl, http.StatusSeeOther)
}
