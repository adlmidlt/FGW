package auth

import (
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/convert"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

const templateHTMLAuth = "../web/html/auth.html"

var UUIDEmployee uuid.UUID

type AuthorizationHandlerHTTP struct {
	employeeService service.EmployeeUseCase
	wLogg           *wlogger.CustomWLogg
}

func NewAuthorizationHandlerHTTP(employeeService service.EmployeeUseCase, wLogg *wlogger.CustomWLogg) *AuthorizationHandlerHTTP {
	return &AuthorizationHandlerHTTP{employeeService, wLogg}
}

func (a *AuthorizationHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc("/", a.ShowAuthForm)
	mux.HandleFunc("/login", a.LogIn)
}
func (a *AuthorizationHandlerHTTP) ShowAuthForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(templateHTMLAuth)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !handler.ExecuteTemplate(tmpl, nil, w, r, a.wLogg) {
		return
	}
}
func (a *AuthorizationHandlerHTTP) LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	employees, err := a.employeeService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, a.wLogg, msg.H7003, err)
		return
	}

	serviceNumberStr, passwd := r.FormValue("serviceNumber"), r.FormValue("passwd")
	serviceNumber := convert.ConvStrToInt(serviceNumberStr)

	if serviceNumberStr == "" || passwd == "" {
		handler.WriteUnauthorized(w, r, a.wLogg, msg.H7003, nil)
		return
	}

	found := false
	for _, employee := range employees {
		if employee.ServiceNumber == serviceNumber && checkPasswd(employee.Passwd, passwd) {
			found = true
			UUIDEmployee = employee.IdEmployee
			break
		}
	}

	if found {
		http.Redirect(w, r, "/fgw/employees", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
}

func checkPasswd(hashedPasswd, plainPasswd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(plainPasswd))
	return err == nil
}
