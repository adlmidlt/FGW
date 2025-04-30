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

func (e *EmployeeHandlerHTTP) Update(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		e.processUpdateFormEmployee(w, r)
	case http.MethodGet:
		e.renderUpdateFormEmployee(w, r)
	default:
		handler.WriteMethodNotAllowed(w, r, e.wLogg, msg.H7002, nil)
	}
}

func (e *EmployeeHandlerHTTP) Add(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, e.wLogg) {
		return
	}

	roleId, err := convert.ParseStrToUUID(r.FormValue(paramRoleId), w, r, e.wLogg)
	if err != nil {
		return
	}

	employee := &entity.Employee{
		ServiceNumber: convert.ConvStrToInt(r.FormValue("serviceNumber")),
		FirstName:     r.FormValue("firstName"),
		LastName:      r.FormValue("lastName"),
		Patronymic:    r.FormValue("patronymic"),
		Passwd:        r.FormValue("passwd"),
		RoleId:        roleId,
	}

	if err = e.employeeService.Add(r.Context(), employee); err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7012, err)

		return
	}
	http.Redirect(w, r, fgwEmployeesStartUrl, http.StatusSeeOther)
}

func (e *EmployeeHandlerHTTP) Delete(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, e.wLogg) {
		return
	}

	idEmployee, err := convert.ParseStrToUUID(r.FormValue(paramIdEmployee), w, r, e.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(r.Context(), idEmployee, w, r, e.wLogg, e.employeeService) {
		return
	}

	if err = e.employeeService.Delete(r.Context(), idEmployee); err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7011, err)

		return
	}
	http.Redirect(w, r, fgwEmployeesStartUrl, http.StatusSeeOther)
}

func (e *EmployeeHandlerHTTP) renderUpdateFormEmployee(w http.ResponseWriter, r *http.Request) {
	idEmployeeStr := r.URL.Query().Get(paramIdEmployee)
	http.Redirect(w, r, fmt.Sprintf("%s?%s=%s", fgwEmployeesStartUrl, paramIdEmployee, idEmployeeStr), http.StatusSeeOther)
}

func (e *EmployeeHandlerHTTP) processUpdateFormEmployee(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handler.WriteBadRequest(w, r, e.wLogg, msg.H7008, err)

		return
	}

	idEmployee, err := convert.ParseStrToUUID(r.FormValue(paramIdEmployee), w, r, e.wLogg)
	if err != nil {
		return
	}

	roleId, err := convert.ParseStrToUUID(r.FormValue(paramRoleId), w, r, e.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByUUID(r.Context(), idEmployee, w, r, e.wLogg, e.employeeService) {
		return
	}

	employee := &entity.Employee{
		IdEmployee:    idEmployee,
		ServiceNumber: convert.ConvStrToInt(r.FormValue("serviceNumber")),
		FirstName:     r.FormValue("firstName"),
		LastName:      r.FormValue("lastName"),
		Patronymic:    r.FormValue("patronymic"),
		Passwd:        r.FormValue("passwd"),
		RoleId:        roleId,
	}

	if err = e.employeeService.Update(r.Context(), idEmployee, employee); err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7009, err)

		return
	}
	http.Redirect(w, r, fgwEmployeesStartUrl, http.StatusSeeOther)
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
