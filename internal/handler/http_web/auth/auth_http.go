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
const templateHTMLNotFound = "../web/html/404.html"
const templateHTMLIndex = "../web/html/index.html"

var UUIDEmployee uuid.UUID
var ServiceNumber int

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
	mux.HandleFunc("/404", a.NotFound)
	mux.HandleFunc("/fgw", a.StartPage)
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
			UUIDEmployee = employee.IdEmployee
			ServiceNumber = employee.ServiceNumber

			if employee.RoleId.String() == "943c699f-8fd3-4707-9db2-944c26ee2afc" {
				http.Redirect(w, r, "/fgw", http.StatusFound)
			} else {
				http.Redirect(w, r, "/покаещенепридумал", http.StatusFound)
			}
			found = true
			break
		}
	}

	if !found {
		a.NotFound(w, r)
	}
}

func (a *AuthorizationHandlerHTTP) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	tmpl, err := template.ParseFiles(templateHTMLNotFound)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !handler.ExecuteTemplate(tmpl, nil, w, r, a.wLogg) {
		return
	}
}

func (a *AuthorizationHandlerHTTP) StartPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(templateHTMLIndex)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !handler.ExecuteTemplate(tmpl, nil, w, r, a.wLogg) {
		return
	}
}

func checkPasswd(hashedPasswd, plainPasswd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(plainPasswd))
	return err == nil
}
