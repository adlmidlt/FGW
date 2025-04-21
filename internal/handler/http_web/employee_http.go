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
)

const (
	templateHtmlEmployeeList = "../web/html/employee_list.html"
	fgwEmployeesStartUrl     = "/fgw/employees"
	paramIdEmployee          = "idEmployee"
	paramRoleId              = "roleId"
)

type EmployeeHandlerHTTP struct {
	roleService     service.RoleUseCase
	employeeService service.EmployeeUseCase
	wLogg           *wlogger.CustomWLogg
}

func NewEmployeeHandlerHTTP(roleService service.RoleUseCase, employeeService service.EmployeeUseCase, wLogg *wlogger.CustomWLogg) *EmployeeHandlerHTTP {
	return &EmployeeHandlerHTTP{roleService: roleService, employeeService: employeeService, wLogg: wLogg}
}

func (e *EmployeeHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwEmployeesStartUrl, e.All)
	mux.HandleFunc(fgwEmployeesStartUrl+"/update", e.Update)
	mux.HandleFunc(fgwEmployeesStartUrl+"/delete", e.Delete)
	mux.HandleFunc(fgwEmployeesStartUrl+"/add", e.Add)
}

func (e *EmployeeHandlerHTTP) All(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, e.wLogg) {
		return
	}

	employees, err := e.employeeService.All(r.Context())
	if err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, r.Method, r.URL.Path, msg.H7003, err)
		http.Error(w, msg.H7003, http.StatusInternalServerError)

		return
	}

	if employees == nil {
		employees = []*entity.Employee{}
	}

	roles, err := e.roleService.All(r.Context())
	if err != nil {
		e.wLogg.LogHttpE(http.StatusInternalServerError, r.Method, r.URL.Path, msg.H7003, err)
		http.Error(w, msg.H7003, http.StatusInternalServerError)

		return
	}

	data := entity.EmployeeList{Employees: employees, Roles: roles}

	if idEmployeeStr := r.URL.Query().Get(paramIdEmployee); idEmployeeStr != "" {
		e.markEditingEmployee(idEmployeeStr, employees)
	}

	tmpl, ok := handler.ParseTemplateHTML(templateHtmlEmployeeList, w, r, e.wLogg)
	if !ok {
		return
	}

	if !handler.ExecuteTemplate(tmpl, data, w, r, e.wLogg) {
		return
	}
}

func (e *EmployeeHandlerHTTP) Update(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		e.processUpdateFormEmployee(writer, request)
	case http.MethodGet:
		e.renderUpdateFormEmployee(writer, request)
	default:
		http.Error(writer, msg.H7002, http.StatusMethodNotAllowed)
	}
}

func (e *EmployeeHandlerHTTP) Add(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, e.wLogg) {
		return
	}

	roleId, err := handler.ParseStrToUUID(request.FormValue(paramRoleId), writer, request, e.wLogg)
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

func (e *EmployeeHandlerHTTP) Delete(writer http.ResponseWriter, request *http.Request) {
	if handler.MethodNotAllowed(writer, request, http.MethodPost, e.wLogg) {
		return
	}

	idEmployee, err := handler.ParseStrToUUID(request.FormValue(paramIdEmployee), writer, request, e.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(request.Context(), idEmployee, writer, request, e.wLogg, e.employeeService) {
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
	idEmployeeStr := request.URL.Query().Get(paramIdEmployee)
	http.Redirect(writer, request, fmt.Sprintf("%s?%s=%s", fgwEmployeesStartUrl, paramIdEmployee, idEmployeeStr), http.StatusSeeOther)
}

func (e *EmployeeHandlerHTTP) processUpdateFormEmployee(writer http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		e.wLogg.LogHttpE(http.StatusBadRequest, request.Method, request.URL.Path, msg.H7008, err)
		http.Error(writer, msg.H7008, http.StatusBadRequest)

		return
	}

	idEmployee, err := handler.ParseStrToUUID(request.FormValue(paramIdEmployee), writer, request, e.wLogg)
	if err != nil {
		return
	}

	roleId, err := handler.ParseStrToUUID(request.FormValue(paramRoleId), writer, request, e.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(request.Context(), idEmployee, writer, request, e.wLogg, e.employeeService) {
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

// markEditingRole помечает сотрудника как редактируемую по её UUID в строковом формате.
func (e *EmployeeHandlerHTTP) markEditingEmployee(idEmployeeStr string, employees []*entity.Employee) {
	if idEmployee, err := uuid.Parse(idEmployeeStr); err == nil {
		for _, employee := range employees {
			if employee.IdEmployee == idEmployee {
				employee.IsEditing = true
			}
		}
	}
}
