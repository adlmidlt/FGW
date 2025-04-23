package json_api

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"net/http"
)

const fgwCatalogStartUrl = "/api/fgw/catalogs"

type CatalogHandlerJSON struct {
	catalogService service.CatalogUseCase
	wLogg          *wlogger.CustomWLogg
}

func NewCatalogHandlerJSON(catalogService service.CatalogUseCase, wLogg *wlogger.CustomWLogg) *CatalogHandlerJSON {
	return &CatalogHandlerJSON{catalogService: catalogService, wLogg: wLogg}
}

func (c *CatalogHandlerJSON) ServeJSONRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwCatalogStartUrl, c.JSONAll)
}

func (c *CatalogHandlerJSON) JSONAll(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, c.wLogg) {
		return
	}

	catalogs, err := c.catalogService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, c.wLogg, msg.H7003, err)

		return
	}

	if catalogs == nil {
		catalogs = []*entity.Catalog{}
	}

	data := entity.CatalogList{Catalogs: catalogs}

	handler.WriteJSON(w, data, c.wLogg)
}
