package main

import (
	"FGW/internal/config"
	"FGW/pkg/convert"
	"FGW/pkg/db"
	"FGW/pkg/wlogger"
	"FGW/pkg/wlogger/msg"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type Role struct {
	idRole string
	number int
	name   string
}

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

	query, err := mssqlDBConn.Query("SELECT * FROM role")
	if err != nil {
		logger.LogE(msg.E3000, err)
	}

	var roles []Role
	if query.Next() {
		var role Role
		if err = query.Scan(&role.idRole, &role.number, &role.name); err != nil {
			logger.LogE(msg.E3002, err)
		}

		role.name, _ = convert.Win1251ToUTF8(role.name)

		roles = append(roles, role)
	}

	logger.LogI(fmt.Sprintf("%+v", roles))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, new project \"FGW\"!!!")
	})

	http.ListenAndServe(":7777", nil)
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
