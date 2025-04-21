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
	templateHtmlCatalogList = "../web/html/catalog_list.html"
	fgwCatalogStartUrl      = "/fgw/catalogs"
)

type CatalogHandlerHTTP struct {
	catalogService service.CatalogUseCase
	packVariant    service.PackVariantUseCase
	wLogg          *wlogger.CustomWLogg
}

func NewCatalogHandlerHTTP(catalogService service.CatalogUseCase, packVariant service.PackVariantUseCase, wLogg *wlogger.CustomWLogg) *CatalogHandlerHTTP {
	return &CatalogHandlerHTTP{catalogService: catalogService, packVariant: packVariant, wLogg: wLogg}
}

func (c *CatalogHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwCatalogStartUrl, c.All)
}

func (c *CatalogHandlerHTTP) All(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodGet, c.wLogg) {
		return
	}

	catalogs, err := c.catalogService.All(r.Context())
	if err != nil {
		c.wLogg.LogHttpE(http.StatusInternalServerError, r.Method, r.URL.Path, msg.H7003, err)
		http.Error(w, msg.H7003, http.StatusInternalServerError)

		return
	}

	if catalogs == nil {
		catalogs = []*entity.Catalog{}
	}

	data := entity.CatalogList{Catalogs: catalogs}

	tmpl, ok := handler.ParseTemplateHTML(templateHtmlCatalogList, w, r, c.wLogg)
	if !ok {
		return
	}

	if !handler.ExecuteTemplate(tmpl, data, w, r, c.wLogg) {
		return
	}
}

func (c *CatalogHandlerHTTP) Add(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, c.wLogg) {
		return
	}

	catalog := &entity.Catalog{
		ParentId:              convert.ConvStrToInt(r.FormValue("parentId")),
		HandbookId:            convert.ConvStrToInt(r.FormValue("handbookId")),
		RecordIndex:           convert.ConvStrToInt(r.FormValue("recordIndex")),
		Name:                  r.FormValue("name"),
		Comment:               r.FormValue("comment"),
		HandbookValueInt1:     convert.ConvStrToInt(r.FormValue("handbookValueInt1")),
		HandbookValueInt2:     convert.ConvStrToInt(r.FormValue("handbookValueInt2")),
		HandbookValueDecimal1: convert.ConvStrToFloat(r.FormValue("handbookValueDecimal1")),
		HandbookValueDecimal2: convert.ConvStrToFloat(r.FormValue("handbookValueDecimal2")),
		HandbookValueBool1:    convert.ConvStrToBool(r.FormValue("handbookValueBool1")),
		HandbookValueBool2:    convert.ConvStrToBool(r.FormValue("handbookValueBool2")),
		IsArchive:             convert.ConvStrToBool(r.FormValue("isArchive")),
		OwnerUser:             convert.ParseStrToUUID(r.FormValue("ownerUser")),
		OwnerUserDateTime:     r.FormValue("ownerUserDateTime"),
		LastUser:              convert.ParseStrToUUID(r.FormValue("lastUser")),
		LastUserDateTime:      r.FormValue("lastUserDateTime"),
	}

	// TODO: Доделать после создания справочника.
	fmt.Println(catalog)

}
