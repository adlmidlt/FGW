package main

import (
	"FGW/pkg/wlogger"
	"errors"
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, new project \"FGW\"!!!")
	})

	logger, _ := wlogger.NewCustomWLogger()
	defer logger.Close()
	logger.LogI("20001 Сервер запущен!")
	logger.LogE("30001 Сервер не запущен", errors.New("server not starting"))
	logger.LogW("30011 Сервер не запущен", errors.New("server not starting"))
	logger.LogHttpI(http.StatusOK, http.MethodGet, "", "10001 Добавление пользователя")
	logger.LogHttpE(http.StatusNotFound, http.MethodPost, "api/json/add", "50001 Ошибка добавления пользователя", errors.New("user adding error"))
	logger.LogHttpW(http.StatusNotFound, http.MethodPost, "api/json/upd", "50011 Ошибка добавления пользователя", errors.New("user adding error"))

	http.ListenAndServe(":7000", nil)
}
