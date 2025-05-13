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
	"time"
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

	if idStr := r.URL.Query().Get(paramIdEmployee); idStr != "" {
		e.markEditingEmployee(idStr, employees)
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

func (e *EmployeeHandlerHTTP) Add(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, e.wLogg) {
		return
	}

	errors := make(map[string]string)

	employees, err := e.employeeService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7003, err)

		return
	}

	maxServiceNumber := 0
	for _, employee := range employees {
		if employee.ServiceNumber > maxServiceNumber {
			maxServiceNumber = employee.ServiceNumber
		}
	}
	serviceNumber := maxServiceNumber + 1

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	patronymic := r.FormValue("patronymic")
	passwd := r.FormValue("passwd")

	switch {
	case len(firstName) > 25:
		errors["firstName"] = msg.J1003
	case len(lastName) > 25:
		errors["lastName"] = msg.J1005
	case len(patronymic) > 25:
		errors["patronymic"] = msg.J1006
	case len(passwd) < 6:
		errors["passwd"] = msg.J1007
	}

	switch {
	case !handler.IsTextOnly(firstName):
		errors["firstName"] = msg.J1004
	case !handler.IsTextOnly(lastName):
		errors["lastName"] = msg.J1008
	case !handler.IsTextOnly(patronymic):
		errors["patronymic"] = msg.J1008
	}

	roleId, err := convert.ParseStrToUUID(r.FormValue(paramRoleId), w, r, e.wLogg)
	if err != nil {
		return
	}

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid.
	ownerUser := convert.ParseUUIDUnsafe(r.FormValue("ownerUser"))
	if ownerUser == uuid.Nil {
		ownerUser = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	}

	ownerUserDateTime := r.FormValue("ownerUserDateTime")
	if ownerUserDateTime == "" {
		ownerUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid.
	lastUser := convert.ParseUUIDUnsafe(r.FormValue("lastUser"))
	if lastUser == uuid.Nil {
		lastUser = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	}

	lastUserDateTime := r.FormValue("lastUserDateTime")
	if lastUserDateTime == "" {
		lastUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	if handler.SendErrorsJSON(w, errors, e.wLogg) {
		return
	}

	employee := &entity.Employee{
		ServiceNumber: serviceNumber,
		FirstName:     firstName,
		LastName:      lastName,
		Patronymic:    patronymic,
		Passwd:        passwd,
		RoleId:        roleId,
		AuditRecord: entity.AuditRecord{
			OwnerUser:         ownerUser,
			OwnerUserDateTime: ownerUserDateTime,
			LastUser:          lastUser,
			LastUserDateTime:  lastUserDateTime,
		},
	}

	if err = e.employeeService.Add(r.Context(), employee); err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7012, err)

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

	errors := make(map[string]string)

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	patronymic := r.FormValue("patronymic")
	passwd := r.FormValue("passwd")

	switch {
	case len(firstName) > 25:
		errors["firstName"] = msg.J1003
	case len(lastName) > 25:
		errors["lastName"] = msg.J1005
	case len(patronymic) > 25:
		errors["patronymic"] = msg.J1006
	case len(passwd) > 6:
		errors["passwd"] = msg.J1007
	}

	switch {
	case !handler.IsTextOnly(firstName):
		errors["firstName"] = msg.J1004
	case !handler.IsTextOnly(lastName):
		errors["lastName"] = msg.J1008
	case !handler.IsTextOnly(patronymic):
		errors["patronymic"] = msg.J1008
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

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid при изменении записи.
	lastUser := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	lastUserDateTime := time.Now().Format("2006-01-02 15:04:05")

	if handler.SendErrorsJSON(w, errors, e.wLogg) {
		return
	}

	employee := &entity.Employee{
		IdEmployee:    idEmployee,
		ServiceNumber: convert.ConvStrToInt(r.FormValue("serviceNumber")),
		FirstName:     firstName,
		LastName:      lastName,
		Patronymic:    patronymic,
		Passwd:        passwd,
		RoleId:        roleId,
		AuditRecord: entity.AuditRecord{
			LastUser:         lastUser,
			LastUserDateTime: lastUserDateTime,
		},
	}

	if err = e.employeeService.Update(r.Context(), idEmployee, employee); err != nil {
		handler.WriteServerError(w, r, e.wLogg, msg.H7009, err)

		return
	}

	fmt.Println(employee)
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
