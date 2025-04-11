package main

import (
	"FGW/internal/config"
	"FGW/pkg/db"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
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

	var cfg config.Config
	if err = cfg.ConfigMSSQL(definitionOfOS()); err != nil {
		logger.LogE(msg.E3102, err)
	}

	mssqlDBConn, err := db.MSSQLConn(cfg)
	if err != nil {
		logger.LogE(msg.E3003, err)
	}
	defer db.CloseDB(mssqlDBConn)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, new project \"FGW\"!!!")
	})

	server := &http.Server{
		Addr:         ":7777",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	logger.LogI(msg.I2003)

	if err = server.ListenAndServe(); err != nil {
		logger.LogE(msg.E3104, err)
		os.Exit(1)
	}
}

// definitionOfOS возвращает путь к конфигурации в зависимости от ОС.
func definitionOfOS() string {
	switch runtime.GOOS {
	case "windows":
		return pathToConfigWindows
	case "linux":
		return pathToConfigLinux
	}

	return ""
}
