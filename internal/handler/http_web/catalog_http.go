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

	// TODO: как фильтр в отдельный метод
	filteredIdParamId := c.filterByHandbook(r)

	var err error
	var catalogs []*entity.Catalog
	if filteredIdParamId >= 0 {
		catalogs, err = c.catalogService.AllFindByNumber(r.Context(), filteredIdParamId)
	} else {
		catalogs, err = c.catalogService.All(r.Context())
	}

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

	data := entity.CatalogList{Catalogs: catalogs, Handbooks: handbooks, SelectedHandbookId: filteredIdParamId}

	if idStr := r.URL.Query().Get("idCatalog"); idStr != "" {
		id := convert.ConvStrToInt(idStr)
		for _, catalog := range catalogs {
			if catalog.IdCatalog == id {
				catalog.IsEditing = true
			}
		}
	}

	tmpl, err := template.New("catalog_list.html").Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
		}).ParseFiles(templateHtmlCatalogList)
	if err != nil {
		handler.WriteServerError(w, r, c.wLogg, msg.H7006, err)

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

	//errors := make(map[string]string)

	catalogs, err := c.catalogService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, c.wLogg, msg.H7003, err)

		return
	}

	parentIdStr := r.FormValue("parentId")
	if parentIdStr == "" {
		parentIdStr = "0"
	}
	parentId := convert.ConvStrToInt(parentIdStr)

	handbookId := convert.ConvStrToInt(r.FormValue("handbookId"))
	if handbookId == 0 {
		handbookId = 0
	}

	recordIndexStr := r.FormValue("recordIndex")
	if recordIndexStr == "" {
		recordIndexStr = "0"
	}
	recordIndex := convert.ConvStrToInt(recordIndexStr)
	for _, catalog := range catalogs {
		if catalog.HandbookId == handbookId {
			if recordIndex == 0 {
				recordIndex = 0
			}
			recordIndex++
		}
	}

	ownerUserDateTime := r.FormValue("ownerUserDateTime")
	if ownerUserDateTime == "" {
		ownerUserDateTime = time.Now().Format("2006-01-02 15:04:05")
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
		HandbookValueInt1:     convert.ParseFormFieldInt(r, "handbookValueInt1"),
		HandbookValueInt2:     convert.ParseFormFieldInt(r, "handbookValueInt2"),
		HandbookValueDecimal1: convert.ParseFormFieldFloat(r, "handbookValueDecimal1"),
		HandbookValueDecimal2: convert.ParseFormFieldFloat(r, "handbookValueDecimal2"),
		HandbookValueBool1:    convert.ParseFormFieldBool(r, "handbookValueBool1"),
		HandbookValueBool2:    convert.ParseFormFieldBool(r, "handbookValueBool2"),
		IsArchive:             convert.ParseFormFieldBool(r, "isArchive"),
		AuditRecord: entity.AuditRecord{
			OwnerUser:         auth.UUIDEmployee,
			OwnerUserDateTime: ownerUserDateTime,
			LastUser:          auth.UUIDEmployee,
			LastUserDateTime:  lastUserDateTime,
		},
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
		handler.WriteMethodNotAllowed(w, r, c.wLogg, msg.H7002, nil)
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
		AuditRecord: entity.AuditRecord{
			LastUser:         auth.UUIDEmployee,
			LastUserDateTime: lastUserDateTime,
		},
	}

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

// filterByHandbook фильтр по справочнику.
func (c *CatalogHandlerHTTP) filterByHandbook(r *http.Request) int {
	var filteredIdParamId int
	handbookIdParam := r.URL.Query().Get("handbookId")
	if handbookIdParam != "" && handbookIdParam != "-1" {
		filteredIdParamId = convert.ConvStrToInt(handbookIdParam)
	} else {
		filteredIdParamId = -1
	}

	return filteredIdParamId
}
