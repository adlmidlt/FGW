package json_api

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/convert"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"encoding/json"
	"fmt"
	"net/http"
)

const fgwHandbookStartUrl = "/api/fgw/handbooks"

type HandbookHandlerJSON struct {
	handbookService service.HandbookUseCase
	wLogg           *wlogger.CustomWLogg
}

func NewHandbookHandlerJSON(handbookService service.HandbookUseCase, wLogg *wlogger.CustomWLogg) *HandbookHandlerJSON {
	return &HandbookHandlerJSON{handbookService: handbookService, wLogg: wLogg}
}

func (h *HandbookHandlerJSON) ServeJSONRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwHandbookStartUrl, h.JSONAll)
	mux.HandleFunc(fgwHandbookStartUrl+"/find", h.JSONFindById)
	mux.HandleFunc(fgwHandbookStartUrl+"/add", h.JSONAdd)
	mux.HandleFunc(fgwHandbookStartUrl+"/update", h.JSONUpdate)
	mux.HandleFunc(fgwHandbookStartUrl+"/delete", h.JSONDelete)
}

func (h *HandbookHandlerJSON) JSONAll(w http.ResponseWriter, r *http.Request) {
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

	handler.WriteJSON(w, handbooks, h.wLogg)
}

func (h *HandbookHandlerJSON) JSONFindById(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, h.wLogg) {
		return
	}

	idHandbook, err := convert.ParseStrToID(r.URL.Query().Get("idHandbook"), w, r, h.wLogg)
	if err != nil {
		return
	}
	if !handler.EntityExistsByID(r.Context(), idHandbook, w, r, h.wLogg, h.handbookService) {
		return
	}

	handbook, err := h.handbookService.FindById(r.Context(), idHandbook)
	if err != nil {
		handler.WriteNotFound(w, r, h.wLogg, msg.H7005, err)

		return
	}

	handler.WriteJSON(w, handbook, h.wLogg)
}

func (h *HandbookHandlerJSON) JSONAdd(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, h.wLogg) {
		return
	}

	var handbook entity.Handbook
	if err := json.NewDecoder(r.Body).Decode(&handbook); err != nil {
		handler.WriteBadRequest(w, r, h.wLogg, msg.H7004, err)

		return
	}

	if err := h.handbookService.Add(r.Context(), &handbook); err != nil {
		handler.WriteServerError(w, r, h.wLogg, msg.H7003, err)

		return
	}

	handler.WriteJSON(w, handbook, h.wLogg)
}

func (h *HandbookHandlerJSON) JSONUpdate(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPut, h.wLogg) {
		return
	}

	idHandbook, err := convert.ParseStrToID(r.URL.Query().Get("idHandbook"), w, r, h.wLogg)
	if err != nil {
		return
	}
	if !handler.EntityExistsByID(r.Context(), idHandbook, w, r, h.wLogg, h.handbookService) {
		return
	}

	var handbook entity.Handbook
	if err = json.NewDecoder(r.Body).Decode(&handbook); err != nil {
		handler.WriteBadRequest(w, r, h.wLogg, msg.H7004, err)

		return
	}

	if err = h.handbookService.Update(r.Context(), idHandbook, &handbook); err != nil {
		handler.WriteServerError(w, r, h.wLogg, msg.H7003, err)

		return
	}
	handler.WriteJSON(w, map[string]string{"message": msg.I2005}, h.wLogg)
}

func (h *HandbookHandlerJSON) JSONDelete(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodDelete, h.wLogg) {
		return
	}

	idHandbook, err := convert.ParseStrToID(r.URL.Query().Get("idHandbook"), w, r, h.wLogg)
	if err != nil {
		return
	}
	fmt.Println(idHandbook)
	if !handler.EntityExistsByID(r.Context(), idHandbook, w, r, h.wLogg, h.handbookService) {
		return
	}

	if err = h.handbookService.Delete(r.Context(), idHandbook); err != nil {
		handler.WriteServerError(w, r, h.wLogg, msg.H7003, err)

		return
	}
	handler.WriteJSON(w, map[string]string{"message": msg.I2004}, h.wLogg)
}
