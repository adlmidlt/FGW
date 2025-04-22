package http_web

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/convert"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	templateHtmlCatalogList = "../web/html/catalog_list.html"
	fgwCatalogStartUrl      = "/fgw/catalogs"
)

type CatalogHandlerHTTP struct {
	catalogService  service.CatalogUseCase
	handbookService service.HandbookUseCase
	wLogg           *wlogger.CustomWLogg
}

func NewCatalogHandlerHTTP(catalogService service.CatalogUseCase, handbookService service.HandbookUseCase, wLogg *wlogger.CustomWLogg) *CatalogHandlerHTTP {
	return &CatalogHandlerHTTP{catalogService: catalogService, handbookService: handbookService, wLogg: wLogg}
}

func (c *CatalogHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwCatalogStartUrl, c.All)
	mux.HandleFunc(fgwCatalogStartUrl+"/add", c.Add)
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

	handbooks, err := c.handbookService.All(r.Context())
	if err != nil {
		c.wLogg.LogHttpE(http.StatusInternalServerError, r.Method, r.URL.Path, msg.H7003, err)
		http.Error(w, msg.H7003, http.StatusInternalServerError)

		return
	}

	data := entity.CatalogList{Catalogs: catalogs, Handbooks: handbooks}

	if idStr := r.URL.Query().Get("idCatalog"); idStr != "" {
		id := convert.ConvStrToInt(idStr)
		for _, catalog := range catalogs {
			if catalog.IdCatalog == id {
				catalog.IsEditing = true
			}
		}
	}

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

	catalogs, err := c.catalogService.All(r.Context())
	if err != nil {
		c.wLogg.LogHttpE(http.StatusInternalServerError, r.Method, r.URL.Path, msg.H7003, err)
		http.Error(w, msg.H7003, http.StatusInternalServerError)

		return
	}

	parentId := convert.ConvStrToInt(r.FormValue("parentId"))
	if parentId == 0 {
		parentId = 0
	}

	handbookId := convert.ConvStrToInt(r.FormValue("handbookId"))
	if handbookId == 0 {
		handbookId = 0
	}

	recordIndex := convert.ConvStrToInt(r.FormValue("recordIndex"))
	for _, catalog := range catalogs {
		if catalog.HandbookId == handbookId {
			if recordIndex == 0 {
				recordIndex = 0
			}
			recordIndex++
		}
	}

	handbookValueInt1 := convert.ConvStrToInt(r.FormValue("handbookValueInt1"))
	if handbookValueInt1 == 0 {
		handbookValueInt1 = 0
	}

	handbookValueInt2 := convert.ConvStrToInt(r.FormValue("handbookValueInt2"))
	if handbookValueInt2 == 0 {
		handbookValueInt2 = 0
	}

	handbookValueDecimal1 := convert.ConvStrToFloat(r.FormValue("handbookValueDecimal1"))
	if handbookValueDecimal1 == 0 {
		handbookValueDecimal1 = 0.0
	}

	handbookValueDecimal2 := convert.ConvStrToFloat(r.FormValue("handbookValueDecimal2"))
	if handbookValueDecimal2 == 0 {
		handbookValueDecimal2 = 0.0
	}

	handbookValueBool1 := convert.ConvStrToBool(r.FormValue("handbookValueBool1"))
	if handbookValueBool1 == false {
		handbookValueBool1 = false
	}

	handbookValueBool2 := convert.ConvStrToBool(r.FormValue("handbookValueBool2"))
	if handbookValueBool2 == false {
		handbookValueBool2 = false
	}

	isArchiveStr := r.FormValue("isArchive")
	isArchive := isArchiveStr == "on" || isArchiveStr == "true"

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid.
	ownerUser := convert.ParseStrToUUID(r.FormValue("ownerUser"))
	if ownerUser == uuid.Nil {
		ownerUser = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	}
	ownerUserDateTime := r.FormValue("ownerUserDateTime")
	if ownerUserDateTime == "" {
		ownerUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid.
	lastUser := convert.ParseStrToUUID(r.FormValue("lastUser"))
	if lastUser == uuid.Nil {
		lastUser = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	}

	lastUserDateTime := r.FormValue("lastUserDateTime")
	if lastUserDateTime == "" {
		lastUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	catalog := &entity.Catalog{
		ParentId:              parentId,
		HandbookId:            handbookId,
		RecordIndex:           recordIndex,
		Name:                  r.FormValue("name"),
		Comment:               r.FormValue("comment"),
		HandbookValueInt1:     handbookValueInt1,
		HandbookValueInt2:     handbookValueInt2,
		HandbookValueDecimal1: handbookValueDecimal1,
		HandbookValueDecimal2: handbookValueDecimal2,
		HandbookValueBool1:    handbookValueBool1,
		HandbookValueBool2:    handbookValueBool2,
		IsArchive:             isArchive,
		OwnerUser:             ownerUser,
		OwnerUserDateTime:     ownerUserDateTime,
		LastUser:              lastUser,
		LastUserDateTime:      lastUserDateTime,
	}
	if err = c.catalogService.Add(r.Context(), catalog); err != nil {
		c.wLogg.LogHttpE(http.StatusInternalServerError, r.Method, r.URL.Path, msg.H7012, err)
		http.Error(w, msg.H7012, http.StatusInternalServerError)

		return
	}
	http.Redirect(w, r, fgwCatalogStartUrl, http.StatusSeeOther)
}
