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
	templateHtmlPackVariantList = "../web/html/pack_variant_list.html"
	fgwPackVariantsStartUrl     = "/fgw/pack_variants"
	paramIdPackVariant          = "idPackVariant"
)

type PackVariantHandlerHTTP struct {
	packVariantService service.PackVariantUseCase
	catalogService     service.CatalogUseCase
	wLogg              *wlogger.CustomWLogg
}

func NewPackVariantHandlerHTTP(packVariantService service.PackVariantUseCase, catalogService service.CatalogUseCase, wLogg *wlogger.CustomWLogg) *PackVariantHandlerHTTP {
	return &PackVariantHandlerHTTP{
		packVariantService: packVariantService,
		catalogService:     catalogService,
		wLogg:              wLogg,
	}
}

func (p *PackVariantHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwPackVariantsStartUrl, p.All)
	mux.HandleFunc(fgwPackVariantsStartUrl+"/add", p.Add)
	mux.HandleFunc(fgwPackVariantsStartUrl+"/update", p.Update)
	mux.HandleFunc(fgwPackVariantsStartUrl+"/delete", p.Delete)
}

func (p *PackVariantHandlerHTTP) All(w http.ResponseWriter, r *http.Request) {
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

	catalogs, err := p.catalogService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7003, err)

		return
	}

	data := entity.PackVariantList{PackVariants: packVariants, Catalogs: catalogs}

	if idStr := r.URL.Query().Get(paramIdPackVariant); idStr != "" {
		id := convert.ConvStrToInt(idStr)
		for _, packVariant := range packVariants {
			if packVariant.IdPackVariant == id {
				packVariant.IsEditing = true
			}
		}
	}

	tmpl, err := template.New("pack_variant_list.html").Funcs(
		template.FuncMap{
			"formatDateTime": convert.FormatDateTime,
		}).ParseFiles(templateHtmlPackVariantList)
	if err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7006, err)

		return
	}

	if !handler.ExecuteTemplate(tmpl, data, w, r, p.wLogg) {
		return
	}
}

func (p *PackVariantHandlerHTTP) Add(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, p.wLogg) {
		return
	}

	catalogs, err := p.catalogService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7003, err)

		return
	}

	prodId := convert.ParseFormFieldInt(r, "prodId")

	packName := r.FormValue("packName")
	for _, catalog := range catalogs {
		if catalog.Name == packName {
			prodId = catalog.IdCatalog
			break
		}
	}

	gl := convert.ParseFormFieldInt(r, "gl")

	color := convert.ParseFormFieldInt(r, "color")
	for _, catalog := range catalogs {
		if catalog.IdCatalog == color {
			gl = catalog.HandbookValueInt1
			break
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

	packVariant := &entity.PackVariant{
		ProdId:            prodId,
		Article:           r.FormValue("article"),
		PackName:          packName,
		Color:             color,
		GL:                gl,
		QuantityRows:      convert.ParseFormFieldInt(r, "quantityRows"),
		QuantityPerRows:   convert.ParseFormFieldInt(r, "quantityPerRows"),
		Weight:            convert.ParseFormFieldInt(r, "weight"),
		Depth:             convert.ParseFormFieldInt(r, "depth"),
		Width:             convert.ParseFormFieldInt(r, "width"),
		Height:            convert.ParseFormFieldInt(r, "height"),
		IsFood:            convert.ParseFormFieldBool(r, "isFood"),
		IsAfraidMoisture:  convert.ParseFormFieldBool(r, "isAfraidMoisture"),
		IsAfraidSun:       convert.ParseFormFieldBool(r, "isAfraidSun"),
		IsEAC:             convert.ParseFormFieldBool(r, "isEAC"),
		IsAccountingBatch: convert.ParseFormFieldBool(r, "isAccountingBatch"),
		MethodShip:        convert.ParseFormFieldBool(r, "methodShip"),
		ShelfLifeMonths:   convert.ParseFormFieldInt(r, "shelfLifeMonths"),
		BathFurnace:       convert.ParseFormFieldInt(r, "bathFurnace"),
		MachineLine:       convert.ParseFormFieldInt(r, "machineLine"),
		IsManufactured:    convert.ParseFormFieldBool(r, "isManufactured"),
		CurrentDateBatch:  r.FormValue("currentDateBatch"),
		NumberingBatch:    convert.ParseFormFieldInt(r, "numberingBatch"),
		IsArchive:         convert.ParseFormFieldBool(r, "isArchive"),
		AuditRecord: entity.AuditRecord{
			OwnerUser:         auth.UUIDEmployee,
			OwnerUserDateTime: ownerUserDateTime,
			LastUser:          auth.UUIDEmployee,
			LastUserDateTime:  lastUserDateTime,
		},
	}

	if err = p.packVariantService.Add(r.Context(), packVariant); err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7012, err)

		return
	}
	http.Redirect(w, r, fgwPackVariantsStartUrl, http.StatusSeeOther)
}

func (p *PackVariantHandlerHTTP) Update(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p.processUpdateFormPackVariant(w, r)
	case http.MethodGet:
		p.renderUpdateFormPackVariant(w, r)
	default:
		handler.WriteMethodNotAllowed(w, r, p.wLogg, msg.H7002, nil)
	}
}

func (p *PackVariantHandlerHTTP) renderUpdateFormPackVariant(w http.ResponseWriter, r *http.Request) {
	idPackVariantStr := r.URL.Query().Get(paramIdPackVariant)
	http.Redirect(w, r, fmt.Sprintf("%s?%s=%s", fgwPackVariantsStartUrl, paramIdPackVariant, idPackVariantStr), http.StatusSeeOther)
}

func (p *PackVariantHandlerHTTP) processUpdateFormPackVariant(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handler.WriteBadRequest(w, r, p.wLogg, msg.H7008, err)

		return
	}

	idPackVariant := convert.ConvStrToInt(r.FormValue(paramIdPackVariant))

	if !handler.EntityExistsByID(r.Context(), idPackVariant, w, r, p.wLogg, p.packVariantService) {
		return
	}

	catalogs, err := p.catalogService.All(r.Context())
	if err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7003, err)

		return
	}

	prodId := convert.ParseFormFieldInt(r, "prodId")
	packName := r.FormValue("packName")
	for _, catalog := range catalogs {
		if catalog.Name == packName {
			prodId = catalog.IdCatalog
			break
		}
	}

	gl := convert.ParseFormFieldInt(r, "gl")
	color := convert.ParseFormFieldInt(r, "color")
	for _, catalog := range catalogs {
		if catalog.IdCatalog == color {
			gl = catalog.HandbookValueInt1
			break
		}
	}

	isFood := r.PostForm.Get("isFood") != ""
	isAfraidMoisture := r.PostForm.Get("isAfraidMoisture") != ""
	isAfraidSun := r.PostForm.Get("isAfraidSun") != ""
	isEAC := r.PostForm.Get("isEAC") != ""
	isAccountingBatch := r.PostForm.Get("isAccountingBatch") != ""
	methodShip := r.PostForm.Get("methodShip") != ""
	isManufactured := r.PostForm.Get("isManufactured") != ""
	isArchive := r.PostForm.Get("isArchive") != ""

	lastUserDateTime := time.Now().Format("2006-01-02 15:04:05")

	packVariant := &entity.PackVariant{
		ProdId:            prodId,
		Article:           r.FormValue("article"),
		PackName:          r.FormValue("packName"),
		Color:             convert.ParseFormFieldInt(r, "color"),
		GL:                gl,
		QuantityRows:      convert.ParseFormFieldInt(r, "quantityRows"),
		QuantityPerRows:   convert.ParseFormFieldInt(r, "quantityPerRows"),
		Weight:            convert.ParseFormFieldInt(r, "weight"),
		Depth:             convert.ParseFormFieldInt(r, "depth"),
		Width:             convert.ParseFormFieldInt(r, "width"),
		Height:            convert.ParseFormFieldInt(r, "height"),
		IsFood:            isFood,
		IsAfraidMoisture:  isAfraidMoisture,
		IsAfraidSun:       isAfraidSun,
		IsEAC:             isEAC,
		IsAccountingBatch: isAccountingBatch,
		MethodShip:        methodShip,
		ShelfLifeMonths:   convert.ParseFormFieldInt(r, "shelfLifeMonths"),
		BathFurnace:       convert.ParseFormFieldInt(r, "bathFurnace"),
		MachineLine:       convert.ParseFormFieldInt(r, "machineLine"),
		IsManufactured:    isManufactured,
		CurrentDateBatch:  r.FormValue("currentDateBatch"),
		NumberingBatch:    convert.ParseFormFieldInt(r, "numberingBatch"),
		IsArchive:         isArchive,
		AuditRecord: entity.AuditRecord{
			LastUser:         auth.UUIDEmployee,
			LastUserDateTime: lastUserDateTime,
		},
	}

	if err = p.packVariantService.Update(r.Context(), idPackVariant, packVariant); err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7012, err)

		return
	}
	http.Redirect(w, r, fgwPackVariantsStartUrl, http.StatusSeeOther)
}

func (p *PackVariantHandlerHTTP) Delete(w http.ResponseWriter, r *http.Request) {
	if handler.MethodNotAllowed(w, r, http.MethodPost, p.wLogg) {
		return
	}

	idPackVariant := convert.ConvStrToInt(r.FormValue(paramIdPackVariant))

	if !handler.EntityExistsByID(r.Context(), idPackVariant, w, r, p.wLogg, p.packVariantService) {
		return
	}

	if err := p.packVariantService.Delete(r.Context(), idPackVariant); err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7011, err)

		return
	}
	http.Redirect(w, r, fgwPackVariantsStartUrl, http.StatusSeeOther)
}
