package http_web

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/handler/http_web/auth"
	"FGW/internal/service"
	"FGW/pkg/convert"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"html/template"
	"net/http"
	"time"
)

const (
	templateHtmlOperationList = "../web/html/operation_list.html"
	fgwOperationStartUrl      = "/fgw/operations"
	paramIdOperation          = "idOperation"
)

type OperationHandlerHTTP struct {
	operationService service.OperationUseCase
	catalogService   service.CatalogUseCase
	wLogg            *wlogger.CustomWLogg
}

func NewOperationHandlerHTTP(operationService service.OperationUseCase, catalogService service.CatalogUseCase, wLogg *wlogger.CustomWLogg) *OperationHandlerHTTP {
	return &OperationHandlerHTTP{operationService, catalogService, wLogg}
}

func (o *OperationHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwOperationStartUrl, o.All)
	mux.HandleFunc(fgwOperationStartUrl+"/add", o.Add)
}

func (o *OperationHandlerHTTP) All(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, o.wLogg) {
		return
	}

	operations, err := o.operationService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, o.wLogg, msg.H7003, err)

		return
	}

	if operations == nil {
		operations = []*entity.Operation{}
	}

	catalogs, err := o.catalogService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, o.wLogg, msg.H7003, err)

		return
	}

	data := entity.OperationList{Operations: operations, Catalogs: catalogs}

	if idStr := r.URL.Query().Get(paramIdOperation); idStr != "" {
		id := convert.ConvStrToInt(idStr)
		for _, operation := range operations {
			if operation.IdOperation == id {
				operation.IsEditing = true
			}
		}
	}

	tmpl, err := template.New("operation_list.html").Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
		}).ParseFiles(templateHtmlOperationList)
	if err != nil {
		handler.WriteServerError(w, r, o.wLogg, msg.H7006, err)

		return
	}

	if !handler.ExecuteTemplate(tmpl, data, w, r, o.wLogg) {
		return
	}
}

func (o *OperationHandlerHTTP) Add(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, o.wLogg) {
		return
	}

	typeOperation := convert.ParseFormFieldInt(r, "typeOperation")

	createDate := r.FormValue("createDate")
	if createDate == "" {
		createDate = time.Now().Format("2006-01-02 15:04:05")
	}

	ownerUserDateTime := r.FormValue("ownerUserDateTime")
	if ownerUserDateTime == "" {
		ownerUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	lastUserDateTime := r.FormValue("lastUserDateTime")
	if lastUserDateTime == "" {
		lastUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	operation := &entity.Operation{
		TypeOperation:     typeOperation,
		CreateDate:        createDate,
		CreateByEmployee:  auth.ServiceNumber,
		DateOrder:         createDate,         // TODO: Когда спецификация будет готова, переписать дату создания ордера.
		ClosedByEmployee:  auth.ServiceNumber, // TODO: Когда спецификация будет готова, переписать табельный номер.
		CodeAccountingObj: convert.ParseFormFieldInt(r, "codeAccountingObj"),
		Appoint:           convert.ParseFormFieldInt(r, "appoint"),
		AuditRecord: entity.AuditRecord{
			OwnerUser:         auth.UUIDEmployee,
			OwnerUserDateTime: ownerUserDateTime,
			LastUser:          auth.UUIDEmployee,
			LastUserDateTime:  lastUserDateTime,
		},
	}

	if err := o.operationService.Add(r.Context(), operation); err != nil {
		handler.WriteServerError(w, r, o.wLogg, msg.H7012, err)

		return
	}
	http.Redirect(w, r, fgwOperationStartUrl, http.StatusSeeOther)
}
