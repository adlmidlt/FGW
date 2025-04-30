package json_api

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/convert"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"encoding/json"
	"net/http"
)

const (
	fgwPackVariantStartUrl = "/api/fgw/packVariants"
	paramIdPackVariant     = "idPackVariant"
)

type PackVariantHandlerJSON struct {
	packVariantService service.PackVariantUseCase
	wLogg              *wlogger.CustomWLogg
}

func NewPackVariantHandlerJSON(packVariantService service.PackVariantUseCase, wLogg *wlogger.CustomWLogg) *PackVariantHandlerJSON {
	return &PackVariantHandlerJSON{packVariantService: packVariantService, wLogg: wLogg}
}

func (p *PackVariantHandlerJSON) ServeJSONRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwPackVariantStartUrl, p.JSONAll)
	mux.HandleFunc(fgwPackVariantStartUrl+"/add", p.JSONAdd)
	mux.HandleFunc(fgwPackVariantStartUrl+"/find", p.JSONFindById)
	mux.HandleFunc(fgwPackVariantStartUrl+"/update", p.JSONUpdate)
	mux.HandleFunc(fgwPackVariantStartUrl+"/delete", p.JSONDelete)
}

func (p *PackVariantHandlerJSON) JSONAll(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, p.wLogg) {
		return
	}

	packVariants, err := p.packVariantService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7003, err)

		return
	}

	if packVariants == nil {
		packVariants = []*entity.PackVariant{}
	}

	data := entity.PackVariantList{PackVariants: packVariants}

	handler.WriteJSON(w, data, p.wLogg)
}

func (p *PackVariantHandlerJSON) JSONFindById(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, p.wLogg) {
		return
	}

	idPackVariant, err := convert.ParseStrToID(r.URL.Query().Get(paramIdPackVariant), w, r, p.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByID(r.Context(), idPackVariant, w, r, p.wLogg, p.packVariantService) {
		return
	}

	packVariant, err := p.packVariantService.FindById(r.Context(), idPackVariant)
	if err != nil {
		handler.WriteNotFound(w, r, p.wLogg, msg.H7005, err)

		return
	}

	handler.WriteJSON(w, packVariant, p.wLogg)
}

func (p *PackVariantHandlerJSON) JSONAdd(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, p.wLogg) {
		return
	}

	var packVariant entity.PackVariant
	if err := json.NewDecoder(r.Body).Decode(&packVariant); err != nil {
		handler.WriteBadRequest(w, r, p.wLogg, msg.H7004, err)

		return
	}

	if err := p.packVariantService.Add(r.Context(), &packVariant); err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7003, err)

		return
	}

	handler.WriteJSON(w, packVariant, p.wLogg)
}

func (p *PackVariantHandlerJSON) JSONUpdate(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPut, p.wLogg) {
		return
	}

	idPackVariant, err := convert.ParseStrToID(r.URL.Query().Get(paramIdPackVariant), w, r, p.wLogg)
	if err != nil {
		return
	}

	var packVariant entity.PackVariant
	if err = json.NewDecoder(r.Body).Decode(&packVariant); err != nil {
		handler.WriteBadRequest(w, r, p.wLogg, msg.H7004, err)

		return
	}

	if !handler.EntityExistsByID(r.Context(), idPackVariant, w, r, p.wLogg, p.packVariantService) {
		return
	}

	if err = p.packVariantService.Update(r.Context(), idPackVariant, &packVariant); err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7003, err)

		return
	}
	handler.WriteJSON(w, map[string]string{"message": msg.I2005}, p.wLogg)
}

func (p *PackVariantHandlerJSON) JSONDelete(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodDelete, p.wLogg) {
		return
	}

	idPackVariant, err := convert.ParseStrToID(r.URL.Query().Get(paramIdPackVariant), w, r, p.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByID(r.Context(), idPackVariant, w, r, p.wLogg, p.packVariantService) {
		return
	}

	if err = p.packVariantService.Delete(r.Context(), idPackVariant); err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7003, err)

		return
	}

	handler.WriteJSON(w, map[string]string{"message": msg.I2004}, p.wLogg)
}
