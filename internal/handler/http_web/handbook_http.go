package http_web

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/convert"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"fmt"
	"net/http"
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

	tmpl, ok := handler.ParseTemplateHTML(templateHtmlHandbookList, w, r, h.wLogg)
	if !ok {
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
		http.Error(w, msg.H7002, http.StatusMethodNotAllowed)
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

	handbook := &entity.Handbook{
		IdHandbook: idHandbook,
		Name:       r.FormValue("name"),
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

	handbook := &entity.Handbook{
		Name: r.PostFormValue("name"),
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
