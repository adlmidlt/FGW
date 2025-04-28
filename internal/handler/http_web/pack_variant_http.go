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
	return &PackVariantHandlerHTTP{packVariantService: packVariantService, catalogService: catalogService, wLogg: wLogg}
}

func (p *PackVariantHandlerHTTP) ServeHTTPRouters(mux *http.ServeMux) {
	mux.HandleFunc(fgwPackVariantsStartUrl, p.All)
	mux.HandleFunc(fgwPackVariantsStartUrl+"/add", p.Add)
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

	fmt.Println(packVariants)

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

	tmpl, ok := handler.ParseTemplateHTML(templateHtmlPackVariantList, w, r, p.wLogg)
	if !ok {
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

	prodId := convert.ConvStrToInt(r.FormValue("prodId"))

	packName := r.FormValue("packName")
	for _, catalog := range catalogs {
		if catalog.Name == packName {
			prodId = catalog.IdCatalog
			break
		}
	}

	gl := convert.ConvStrToInt(r.FormValue("gl"))

	color := convert.ConvStrToInt(r.FormValue("color"))
	for _, catalog := range catalogs {
		if catalog.IdCatalog == color {
			gl = catalog.HandbookValueInt1
		}
	}
	isFood := convert.ConvStrToBool(r.FormValue("isFood"))
	if isFood == false {
		isFood = false
	}
	isAfraidMoisture := convert.ConvStrToBool(r.FormValue("isAfraidMoisture"))
	if isAfraidMoisture == false {
		isAfraidMoisture = false
	}
	isAfraidSun := convert.ConvStrToBool(r.FormValue("isAfraidSun"))
	if isAfraidSun == false {
		isAfraidSun = false
	}
	isEAC := convert.ConvStrToBool(r.FormValue("isEAC"))
	if isEAC == false {
		isEAC = false
	}
	isAccountingBatch := convert.ConvStrToBool(r.FormValue("isAccountingBatch"))
	if isAccountingBatch == false {
		isAccountingBatch = false
	}
	methodShip := convert.ConvStrToBool(r.FormValue("methodShip"))
	if methodShip == false {
		methodShip = false
	}

	isManufactured := convert.ConvStrToBool(r.FormValue("isManufactured"))
	if isManufactured == false {
		isManufactured = false
	}
	isArchive := convert.ConvStrToBool(r.FormValue("isArchive"))
	if isArchive == false {
		isArchive = false
	}

	currentDateBatch := r.FormValue("currentDateBatch")
	if currentDateBatch == "" {
		currentDateBatch = time.Now().Format("2006-01-02 15:04:05")
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

	packVariant := &entity.PackVariant{
		ProdId:            prodId,
		Article:           r.FormValue("article"),
		PackName:          packName,
		Color:             color,
		GL:                gl,
		QuantityRows:      convert.ConvStrToInt(r.FormValue("quantityRows")),
		QuantityPerRows:   convert.ConvStrToInt(r.FormValue("quantityPerRows")),
		Weight:            convert.ConvStrToInt(r.FormValue("weight")),
		Depth:             convert.ConvStrToInt(r.FormValue("depth")),
		Width:             convert.ConvStrToInt(r.FormValue("width")),
		Height:            convert.ConvStrToInt(r.FormValue("height")),
		IsFood:            isFood,
		IsAfraidMoisture:  isAfraidMoisture,
		IsAfraidSun:       isAfraidSun,
		IsEAC:             isEAC,
		IsAccountingBatch: isAccountingBatch,
		MethodShip:        methodShip,
		ShelfLifeMonths:   convert.ConvStrToInt(r.FormValue("shelfLifeMonths")),
		BathFurnace:       convert.ConvStrToInt(r.FormValue("bathFurnace")),
		MachineLine:       convert.ConvStrToInt(r.FormValue("machineLine")),
		IsManufactured:    isManufactured,
		CurrentDateBatch:  currentDateBatch,
		NumberingBatch:    convert.ConvStrToInt(r.FormValue("numberingBatch")),
		IsArchive:         isArchive,
		OwnerUser:         ownerUser,
		OwnerUserDateTime: ownerUserDateTime,
		LastUser:          lastUser,
		LastUserDateTime:  lastUserDateTime,
	}

	if err = p.packVariantService.Add(r.Context(), packVariant); err != nil {
		handler.WriteServerError(w, r, p.wLogg, msg.H7012, err)

		return
	}
	http.Redirect(w, r, fgwPackVariantsStartUrl, http.StatusSeeOther)
}
