package json_api

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"net/http"
)

const fgwPackVariantStartUrl = "/api/fgw/pack_variants"

type PackVariantHandlerJSON struct {
	packVariantService service.PackVariantUseCase
	wLogg              *wlogger.CustomWLogg
}

func NewPackVariantHandlerJSON(packVariantService service.PackVariantUseCase, wLogg *wlogger.CustomWLogg) *PackVariantHandlerJSON {
	return &PackVariantHandlerJSON{packVariantService: packVariantService, wLogg: wLogg}
}

func (p *PackVariantHandlerJSON) ServeJSONRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwPackVariantStartUrl, p.JSONAll)
}

func (p *PackVariantHandlerJSON) JSONAll(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, p.wLogg) {
		return
	}

	packVariants, err := p.packVariantService.All(r.Context())
	if err != nil {
		p.wLogg.LogHttpE(http.StatusInternalServerError, r.Method, r.URL.Path, msg.H7003, err)
		http.Error(w, msg.H7003, http.StatusInternalServerError)

		return
	}

	if packVariants == nil {
		packVariants = []*entity.PackVariant{}
	}

	data := entity.PackVariantList{PackVariants: packVariants}

	handler.WriteJSON(w, data, p.wLogg)
}
