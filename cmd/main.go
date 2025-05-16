package main

import (
	"FGW/internal/config"
	"FGW/internal/handler/http_web"
	"FGW/internal/handler/http_web/auth"
	"FGW/internal/handler/json_api"
	"FGW/internal/repo"
	"FGW/internal/service"
	"FGW/pkg/db"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"context"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	pathToConfigLinux   = "../internal/config/config.yaml"
	pathToConfigWindows = "internal/config/config.yaml"
)

func main() {
	logger, err := wlogger.NewCustomWLogger()
	if err != nil {
		log.Printf("%s: %v", msg.E3103, err)
	}
	defer logger.Close()

	validateStruct := validator.New()

	ctxInit, cancelInit := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelInit()

	var cfg config.Config
	if err = cfg.ConfigMSSQL(definitionOfOS()); err != nil {
		logger.LogE(msg.E3102, err)
	}

	mssqlDBConn, err := db.MSSQLConn(ctxInit, cfg)
	if err != nil {
		logger.LogE(msg.E3003, err)
	}
	defer db.CloseDB(mssqlDBConn)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repoRole := repo.NewRoleRepo(mssqlDBConn, logger)
	serviceRole := service.NewRoleService(repoRole, logger, validateStruct)
	handlerRoleJSON := json_api.NewRoleHandlerJSON(serviceRole, logger)
	handlerRoleHTTP := http_web.NewRoleHandlerHTTP(serviceRole, logger)

	repoEmployee := repo.NewEmployeeRepo(mssqlDBConn, logger)
	serviceEmployee := service.NewEmployeeService(repoEmployee, logger, validateStruct)
	handlerEmployeeJSON := json_api.NewEmployeeHandlerJSON(serviceRole, serviceEmployee, logger)
	handlerEmployeeHTTP := http_web.NewEmployeeHandlerHTTP(serviceRole, serviceEmployee, logger)
	handlerAuthorizationHTTP := auth.NewAuthorizationHandlerHTTP(serviceEmployee, logger)

	repoHandbook := repo.NewHandbookRepo(mssqlDBConn, logger)
	serviceHandbook := service.NewHandbookService(repoHandbook, logger, validateStruct)
	handlerHandbookJSON := json_api.NewHandbookHandlerJSON(serviceHandbook, logger)
	handlerHandbookHTTP := http_web.NewHandbookHandlerHTTP(serviceHandbook, logger)

	repoCatalog := repo.NewCatalogRepo(mssqlDBConn, logger)
	serviceCatalog := service.NewCatalogService(repoCatalog, logger, validateStruct)
	handlerCatalogJSON := json_api.NewCatalogHandlerJSON(serviceCatalog, logger)
	handlerCatalogHTTP := http_web.NewCatalogHandlerHTTP(serviceCatalog, serviceHandbook, logger)

	repoPackVariant := repo.NewPackVariant(mssqlDBConn, logger)
	servicePackVariant := service.NewPackVariantService(repoPackVariant, logger, validateStruct)
	handlerPackVariantJSON := json_api.NewPackVariantHandlerJSON(servicePackVariant, logger)
	handlerPackVariantHTTP := http_web.NewPackVariantHandlerHTTP(servicePackVariant, serviceCatalog, logger)

	repoOperation := repo.NewOperationRepo(mssqlDBConn, logger)
	serviceOperation := service.NewOperationService(repoOperation, logger)
	handlerOperationHTTP := http_web.NewOperationHandlerHTTP(serviceOperation, serviceCatalog, logger)

	mux := http.NewServeMux()

	handlerRoleJSON.ServeJSONRouters(mux)
	handlerRoleHTTP.ServeHTTPRouters(mux)

	handlerEmployeeJSON.ServeJSONRouters(mux)
	handlerEmployeeHTTP.ServeHTTPRouters(mux)
	handlerAuthorizationHTTP.ServeHTTPRouters(mux)

	handlerHandbookJSON.ServeJSONRouters(mux)
	handlerHandbookHTTP.ServeHTTPRouters(mux)

	handlerCatalogJSON.ServeJSONRouters(mux)
	handlerCatalogHTTP.ServeHTTPRouters(mux)

	handlerPackVariantJSON.ServeJSONRouters(mux)
	handlerPackVariantHTTP.ServeHTTPRouters(mux)

	handlerOperationHTTP.ServeHTTPRouters(mux)

	// Подключение static (*.html, *.png/jpg *.css файлов, *.js)
	mux.Handle("/web/",
		http.StripPrefix("/web/", http.FileServer(http.Dir("../web"))))

	server := &http.Server{
		Addr:         ":7777",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	logger.LogI(msg.I2003)

	go func() {
		if err = server.ListenAndServe(); err != nil {
			logger.LogE(msg.E3104, err)
		}
	}()
	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.LogE(msg.E3106, err)
	}
}

// definitionOfOS возвращает путь к конфигурации в зависимости от ОС.
func definitionOfOS() string {
	if runtime.GOOS == "windows" {
		return pathToConfigWindows
	}

	if runtime.GOOS == "linux" {
		return pathToConfigLinux
	}

	return ""
}
