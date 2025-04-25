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
	mux.HandleFunc(fgwCatalogStartUrl+"/find", c.JSONFindById)
	mux.HandleFunc(fgwCatalogStartUrl+"/add", c.JSONAdd)
	mux.HandleFunc(fgwCatalogStartUrl+"/update", c.JSONUpdate)
	mux.HandleFunc(fgwCatalogStartUrl+"/delete", c.JSONDelete)
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

func (c *CatalogHandlerJSON) JSONFindById(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, c.wLogg) {
		return
	}

	idCatalog, err := convert.ParseStrToID(r.URL.Query().Get("idCatalog"), w, r, c.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByID(r.Context(), idCatalog, w, r, c.wLogg, c.catalogService) {
		return
	}

	catalog, err := c.catalogService.FindById(r.Context(), idCatalog)
	if err != nil {
		handler.WriteNotFound(w, r, c.wLogg, msg.H7005, err)

		return
	}

	handler.WriteJSON(w, catalog, c.wLogg)
}

func (c *CatalogHandlerJSON) JSONAdd(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, c.wLogg) {
		return
	}

	var catalog entity.Catalog
	if err := json.NewDecoder(r.Body).Decode(&catalog); err != nil {
		handler.WriteBadRequest(w, r, c.wLogg, msg.H7004, err)

		return
	}

	if err := c.catalogService.Add(r.Context(), &catalog); err != nil {
		handler.WriteServerError(w, r, c.wLogg, msg.H7003, err)

		return
	}

	handler.WriteJSON(w, catalog, c.wLogg)
}

func (c *CatalogHandlerJSON) JSONUpdate(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPut, c.wLogg) {
		return
	}

	idCatalog, err := convert.ParseStrToID(r.URL.Query().Get("idCatalog"), w, r, c.wLogg)
	if err != nil {
		return
	}

	var catalog entity.Catalog
	if err = json.NewDecoder(r.Body).Decode(&catalog); err != nil {
		handler.WriteBadRequest(w, r, c.wLogg, msg.H7004, err)

		return
	}

	if !handler.EntityExistsByID(r.Context(), idCatalog, w, r, c.wLogg, c.catalogService) {
		return
	}

	if err = c.catalogService.Update(r.Context(), idCatalog, &catalog); err != nil {
		handler.WriteServerError(w, r, c.wLogg, msg.H7003, err)

		return
	}
	handler.WriteJSON(w, map[string]string{"message": msg.I2005}, c.wLogg)
}

func (c *CatalogHandlerJSON) JSONDelete(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodDelete, c.wLogg) {
		return
	}

	idCatalog, err := convert.ParseStrToID(r.URL.Query().Get("idCatalog"), w, r, c.wLogg)
	if err != nil {
		return
	}

	if !handler.EntityExistsByID(r.Context(), idCatalog, w, r, c.wLogg, c.catalogService) {
		return
	}

	if err = c.catalogService.Delete(r.Context(), idCatalog); err != nil {
		handler.WriteServerError(w, r, c.wLogg, msg.H7003, err)

		return
	}

	handler.WriteJSON(w, map[string]string{"message": msg.I2004}, c.wLogg)
}
