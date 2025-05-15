package http_web

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/handler/http_web/auth"
	"FGW/internal/service"
	"FGW/pkg/convert"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

const (
	templateHtmlHandbookList = "../web/html/handbook_list.html"
	fgwHandbookStartUrl      = "/fgw/handbooks"
	paramIdHandbook          = "idHandbook"
)

type HandbookHandlerHTTP struct {
	handbookService service.HandbookUseCase
	wLogg           *wlogger.CustomWLogg
}

func NewHandbookHandlerHTTP(handbookService service.HandbookUseCase, wLogg *wlogger.CustomWLogg) *HandbookHandlerHTTP {
	return &HandbookHandlerHTTP{handbookService: handbookService, wLogg: wLogg}
}

func (h *HandbookHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwHandbookStartUrl, h.All)
	mux.HandleFunc(fgwHandbookStartUrl+"/update", h.Update)
	mux.HandleFunc(fgwHandbookStartUrl+"/delete", h.Delete)
	mux.HandleFunc(fgwHandbookStartUrl+"/add", h.Add)
}

func (h *HandbookHandlerHTTP) All(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, h.wLogg) {
		return
	}

	handbooks, err := h.handbookService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, h.wLogg, msg.H7003, err)

		return
	}

	if handbooks == nil {
		handbooks = []*entity.Handbook{}
		if err = h.handbookService.AddZeroObj(r.Context()); err != nil {
			handler.WriteServerError(w, r, h.wLogg, msg.H7010, err)

			return
		}
	}

	data := entity.HandbookList{Handbooks: handbooks}

	if idStr := r.URL.Query().Get(paramIdHandbook); idStr != "" {
		id := convert.ConvStrToInt(idStr)
		for _, handbook := range handbooks {
			if handbook.IdHandbook == id {
				handbook.IsEditing = true
			}
		}
	}

	tmpl, err := template.New("handbook_list.html").Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
		}).ParseFiles(templateHtmlHandbookList)
	if err != nil {
		handler.WriteServerError(w, r, h.wLogg, msg.H7006, err)

		return
	}

	if !handler.ExecuteTemplate(tmpl, data, w, r, h.wLogg) {
		return
	}
}

func (h *HandbookHandlerHTTP) Update(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.processUpdateFormHandbook(w, r)
	case http.MethodGet:
		h.renderUpdateFormHandbook(w, r)
	default:
		handler.WriteMethodNotAllowed(w, r, h.wLogg, msg.H7002, nil)
	}
}

func (h *HandbookHandlerHTTP) renderUpdateFormHandbook(w http.ResponseWriter, r *http.Request) {
	idHandbookStr := r.URL.Query().Get(paramIdHandbook)
	http.Redirect(w, r, fmt.Sprintf("%s?%s=%s", fgwHandbookStartUrl, paramIdHandbook, idHandbookStr), http.StatusSeeOther)
}

func (h *HandbookHandlerHTTP) processUpdateFormHandbook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handler.WriteBadRequest(w, r, h.wLogg, msg.H7008, err)

		return
	}

	idHandbook := convert.ConvStrToInt(r.FormValue(paramIdHandbook))

	if !handler.EntityExistsByID(r.Context(), idHandbook, w, r, h.wLogg, h.handbookService) {
		return
	}

	lastUserDateTime := time.Now().Format("2006-01-02 15:04:05")

	handbook := &entity.Handbook{
		IdHandbook: idHandbook,
		Name:       r.FormValue("name"),
		AuditRecord: entity.AuditRecord{
			LastUser:         auth.UUIDEmployee,
			LastUserDateTime: lastUserDateTime,
		},
	}

	if err := h.handbookService.Update(r.Context(), idHandbook, handbook); err != nil {
		handler.WriteServerError(w, r, h.wLogg, msg.H7009, err)

		return
	}
	http.Redirect(w, r, fgwHandbookStartUrl, http.StatusSeeOther)
}

func (h *HandbookHandlerHTTP) Add(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, h.wLogg) {
		return
	}

	ownerUserDateTime := r.FormValue("ownerUserDateTime")
	if ownerUserDateTime == "" {
		ownerUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	lastUserDateTime := r.FormValue("lastUserDateTime")
	if lastUserDateTime == "" {
		lastUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	handbook := &entity.Handbook{
		Name: r.PostFormValue("name"),
		AuditRecord: entity.AuditRecord{
			OwnerUser:         auth.UUIDEmployee,
			OwnerUserDateTime: ownerUserDateTime,
			LastUser:          auth.UUIDEmployee,
			LastUserDateTime:  lastUserDateTime,
		},
	}

	if err := h.handbookService.Add(r.Context(), handbook); err != nil {
		handler.WriteServerError(w, r, h.wLogg, msg.H7012, err)

		return
	}
	http.Redirect(w, r, fgwHandbookStartUrl, http.StatusSeeOther)
}

func (h *HandbookHandlerHTTP) Delete(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, h.wLogg) {
		return
	}

	idHandbook := convert.ConvStrToInt(r.FormValue(paramIdHandbook))

	if !handler.EntityExistsByID(r.Context(), idHandbook, w, r, h.wLogg, h.handbookService) {
		return
	}

	if err := h.handbookService.Delete(r.Context(), idHandbook); err != nil {
		handler.WriteServerError(w, r, h.wLogg, msg.H7011, err)

		return
	}
	http.Redirect(w, r, fgwHandbookStartUrl, http.StatusSeeOther)
}
