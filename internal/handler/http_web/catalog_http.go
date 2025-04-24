package http_web

import (
	"FGW/internal/entity"
	"FGW/internal/handler"
	"FGW/internal/service"
	"FGW/pkg/convert"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	templateHtmlCatalogList = "../web/html/catalog_list.html"
	fgwCatalogStartUrl      = "/fgw/catalogs"
	paramIdCatalog          = "idCatalog"
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
	mux.HandleFunc(fgwCatalogStartUrl+"/update", c.Update)
	mux.HandleFunc(fgwCatalogStartUrl+"/delete", c.Delete)
}

func (c *CatalogHandlerHTTP) All(w http.ResponseWriter, r *http.Request) {
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

	handbooks, err := c.handbookService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, c.wLogg, msg.H7003, err)

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
		handler.WriteServerError(w, r, c.wLogg, msg.H7003, err)

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

	isArchive := convert.ConvStrToBool(r.FormValue("isArchive"))
	if isArchive == false {
		isArchive = false
	}

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid.
	ownerUser := convert.ParseUUIDUnsafe(r.FormValue("ownerUser"))
	if ownerUser == uuid.Nil {
		ownerUser = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	}
	ownerUserDateTime := r.FormValue("ownerUserDateTime")
	if ownerUserDateTime == "" {
		ownerUserDateTime = time.Now().Format("2006-01-02 15:04:05")
	}

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid.
	lastUser := convert.ParseUUIDUnsafe(r.FormValue("lastUser"))
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
		handler.WriteServerError(w, r, c.wLogg, msg.H7012, err)

		return
	}
	http.Redirect(w, r, fgwCatalogStartUrl, http.StatusSeeOther)
}

func (c *CatalogHandlerHTTP) Update(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.processUpdateFormEmployee(w, r)
	case http.MethodGet:
		c.renderUpdateFormEmployee(w, r)
	default:
		http.Error(w, msg.H7002, http.StatusMethodNotAllowed)
	}
}

func (c *CatalogHandlerHTTP) renderUpdateFormEmployee(w http.ResponseWriter, r *http.Request) {
	idCatalogStr := r.URL.Query().Get(paramIdCatalog)
	http.Redirect(w, r, fmt.Sprintf("%s?%s=%s", fgwCatalogStartUrl, paramIdCatalog, idCatalogStr), http.StatusSeeOther)
}

func (c *CatalogHandlerHTTP) processUpdateFormEmployee(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handler.WriteBadRequest(w, r, c.wLogg, msg.H7008, err)

		return
	}

	idCatalog := convert.ConvStrToInt(r.FormValue(paramIdCatalog))

	if !handler.EntityExistsByID(r.Context(), idCatalog, w, r, c.wLogg, c.catalogService) {
		return
	}

	handbookValueBool1 := r.PostForm.Get("handbookValueBool1") != ""
	handbookValueBool2 := r.PostForm.Get("handbookValueBool2") != ""
	isArchive := r.PostForm.Get("isArchive") != ""

	// TODO: временная заглушка, после написания авторизации, будет заполняться uuid при изменении записи.
	lastUser := uuid.MustParse("10000000-0000-0000-0000-000000000000")

	lastUserDateTime := time.Now().Format("2006-01-02 15:04:05")

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
		HandbookValueBool1:    handbookValueBool1,
		HandbookValueBool2:    handbookValueBool2,
		IsArchive:             isArchive,
		OwnerUser:             convert.ParseUUIDUnsafe(r.FormValue("ownerUser")),
		OwnerUserDateTime:     r.FormValue("ownerUserDateTime"),
		LastUser:              lastUser,
		LastUserDateTime:      lastUserDateTime,
	}
	fmt.Println(catalog)
	if err := c.catalogService.Update(r.Context(), idCatalog, catalog); err != nil {
		handler.WriteServerError(w, r, c.wLogg, msg.H7012, err)

		return
	}
	http.Redirect(w, r, fgwCatalogStartUrl, http.StatusSeeOther)
}

func (c *CatalogHandlerHTTP) Delete(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, c.wLogg) {
		return
	}

	idCatalog := convert.ConvStrToInt(r.FormValue(paramIdCatalog))

	if !handler.EntityExistsByID(r.Context(), idCatalog, w, r, c.wLogg, c.catalogService) {
		return
	}

	if err := c.catalogService.Delete(r.Context(), idCatalog); err != nil {
		handler.WriteServerError(w, r, c.wLogg, msg.H7011, err)

		return
	}
	http.Redirect(w, r, fgwCatalogStartUrl, http.StatusSeeOther)
}
