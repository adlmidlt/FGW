package wlogger

import (
	"FGW/pkg"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	systemDateTime = time.Now().Format(time.DateTime)
	pathToLoggFile = os.Getenv("GOPATH") + "/fgw.json"
)

const (
	infoInConsoleInJSON = "\u001B[1;32m{\n" +
		"\u001B[1;32m  \"dataTime\":\u001B[1;38m \"%s\",\n" +
		"\u001B[1;32m  \"info\":{\n" +
		"\u001B[1;32m    \"hostname\":  \u001B[1;32m \"%s\",\n" +
		"\u001B[1;32m    \"ipAddr\":    \u001B[1;32m \"%s\",\n" +
		"\u001B[1;32m    \"pathToCode\":\u001B[1;32m %s ,\n" +
		"\u001B[1;32m    \"lineCode\":  \u001B[1;32m %s ,\n" +
		"\u001B[1;32m    \"code\":      \u001B[1;32m \"%s\",\n" +
		"\u001B[1;32m    \"message\":   \u001B[1;32m \"%s\"\n" +
		"\u001B[1;32m }\n" +
		"\u001B[1;32m},\u001B[0m\n"

	infoInJSON = "{\n" +
		" \"dataTime\": \"%s\",\n" +
		" \"info\":{\n" +
		"   \"hostname\":   \"%s\",\n" +
		"   \"ipAddr\":     \"%s\",\n" +
		"   \"pathToCode\": \"%s\",\n" +
		"   \"lineCode\":   \"%s\",\n" +
		"   \"code\":       \"%s\",\n" +
		"   \"message\":    \"%s\"\n" +
		" }\n" +
		"},\n"

	warningInConsoleInJson = "\u001B[1;33m{\n" +
		"\u001B[1;33m  \"dataTime\":\u001B[1;38m \"%s\",\n" +
		"\u001B[1;33m  \"warning\":{\n" +
		"\u001B[1;33m    \"hostname\":  \u001B[1;33m \"%s\",\n" +
		"\u001B[1;33m    \"ipAddr\":    \u001B[1;33m \"%s\",\n" +
		"\u001B[1;33m    \"pathToCode\":\u001B[1;33m %s ,\n" +
		"\u001B[1;33m    \"lineCode\":  \u001B[1;33m %s ,\n" +
		"\u001B[1;33m    \"code\":      \u001B[1;33m \"%s\",\n" +
		"\u001B[1;33m    \"message\":   \u001B[1;33m \"%s\",\n" +
		"\u001B[1;33m    \"error\":     \u001B[1;33m \"%s\"\n" +
		"\u001B[1;33m }\n" +
		"\u001B[1;33m},\u001B[0m\n"

	warningInJson = "{\n" +
		" \"dataTime\": \"%s\",\n" +
		" \"warning\":{\n" +
		"   \"hostname\":   \"%s\",\n" +
		"   \"ipAddr\":     \"%s\",\n" +
		"   \"pathToCode\": \"%s\",\n" +
		"   \"lineCode\":   \"%s\",\n" +
		"   \"code\":       \"%s\",\n" +
		"   \"message\":    \"%s\",\n" +
		"   \"error\":      \"%s\"\n" +
		" }\n" +
		"},\n"

	errorInConsoleInJSON = "\u001B[1;31m{\n" +
		"\u001B[1;31m  \"dataTime\":\u001B[1;38m \"%s\",\n" +
		"\u001B[1;31m  \"error\":{\n" +
		"\u001B[1;31m    \"hostname\":  \u001B[1;31m \"%s\",\n" +
		"\u001B[1;31m    \"ipAddr\":    \u001B[1;31m \"%s\",\n" +
		"\u001B[1;31m    \"pathToCode\":\u001B[1;31m %s ,\n" +
		"\u001B[1;31m    \"lineCode\":  \u001B[1;31m %s ,\n" +
		"\u001B[1;31m    \"code\":      \u001B[1;31m \"%s\",\n" +
		"\u001B[1;31m    \"message\":   \u001B[1;31m \"%s\",\n" +
		"\u001B[1;31m    \"error\":     \u001B[1;31m \"%s\"\n" +
		"\u001B[1;31m }\n" +
		"\u001B[1;31m},\u001B[0m\n"

	errorInJSON = "{\n" +
		" \"dataTime\": \"%s\",\n" +
		" \"error\":{\n" +
		"   \"hostname\":   \"%s\",\n" +
		"   \"ipAddr\":     \"%s\",\n" +
		"   \"pathToCode\": \"%s\",\n" +
		"   \"lineCode\":   \"%s\",\n" +
		"   \"code\":       \"%s\",\n" +
		"   \"message\":    \"%s\",\n" +
		"   \"error\":      \"%s\"\n" +
		" }\n" +
		"},\n"

	httpInfoInConsoleInJSON = "\u001B[1;32m{\n" +
		"\u001B[1;32m  \"dataTime\":\u001B[1;38m \"%s\",\n" +
		"\u001B[1;32m  \"httpInfo\":{\n" +
		"\u001B[1;32m    \"hostname\":  \u001B[1;32m \"%s\",\n" +
		"\u001B[1;32m    \"ipAddr\":    \u001B[1;32m \"%s\",\n" +
		"\u001B[1;32m    \"pathToCode\":\u001B[1;32m %s ,\n" +
		"\u001B[1;32m    \"lineCode\":  \u001B[1;32m %s ,\n" +
		"\u001B[1;32m    \"statusCode\":\u001B[1;32m \"%d\",\n" +
		"\u001B[1;32m    \"methodHttp\":\u001B[1;32m \"%s\",\n" +
		"\u001B[1;32m    \"url\":       \u001B[1;32m \"%s\",\n" +
		"\u001B[1;32m    \"code\":      \u001B[1;32m \"%s\",\n" +
		"\u001B[1;32m    \"message\":   \u001B[1;32m \"%s\"\n" +
		"\u001B[1;32m }\n" +
		"\u001B[1;32m},\u001B[0m\n"

	httpInfoInJSON = "{\n" +
		" \"dataTime\": \"%s\",\n" +
		" \"httpInfo\":{\n" +
		"   \"hostname\":   \"%s\",\n" +
		"   \"ipAddr\":     \"%s\",\n" +
		"   \"pathToCode\": \"%s\",\n" +
		"   \"lineCode\":   \"%s\",\n" +
		"   \"statusCode\": \"%d\",\n" +
		"   \"methodHttp\": \"%s\",\n" +
		"   \"url\":        \"%s\",\n" +
		"   \"code\":       \"%s\",\n" +
		"   \"message\":    \"%s\"\n" +
		" }\n" +
		"},\n"

	httpWarningInConsoleInJSON = "\u001B[1;33m{\n" +
		"\u001B[1;33m  \"dataTime\":\u001B[1;38m \"%s\",\n" +
		"\u001B[1;33m  \"httpWarning\":{\n" +
		"\u001B[1;33m    \"hostname\":  \u001B[1;33m \"%s\",\n" +
		"\u001B[1;33m    \"ipAddr\":    \u001B[1;33m \"%s\",\n" +
		"\u001B[1;33m    \"pathToCode\":\u001B[1;33m %s ,\n" +
		"\u001B[1;33m    \"lineCode\":  \u001B[1;33m %s ,\n" +
		"\u001B[1;33m    \"statusCode\":\u001B[1;33m \"%d\",\n" +
		"\u001B[1;33m    \"methodHttp\":\u001B[1;33m \"%s\",\n" +
		"\u001B[1;33m    \"url\":       \u001B[1;33m \"%s\",\n" +
		"\u001B[1;33m    \"code\":      \u001B[1;33m \"%s\",\n" +
		"\u001B[1;33m    \"message\":   \u001B[1;33m \"%s\",\n" +
		"\u001B[1;33m    \"error\":     \u001B[1;33m \"%s\"\n" +
		"\u001B[1;33m }\n" +
		"\u001B[1;33m},\u001B[0m\n"

	httpWarningInJSON = "{\n" +
		" \"dataTime\": \"%s\",\n" +
		" \"httpWarning\":{\n" +
		"   \"hostname\":   \"%s\",\n" +
		"   \"ipAddr\":     \"%s\",\n" +
		"   \"pathToCode\": \"%s\",\n" +
		"   \"lineCode\":   \"%s\",\n" +
		"   \"statusCode\": \"%d\",\n" +
		"   \"methodHttp\": \"%s\",\n" +
		"   \"url\":        \"%s\",\n" +
		"   \"code\":       \"%s\",\n" +
		"   \"message\":    \"%s\",\n" +
		"   \"error\":      \"%s\"\n" +
		" }\n" +
		"},\n"

	httpErrorInConsoleInJSON = "\u001B[1;31m{\n" +
		"\u001B[1;31m  \"dataTime\":\u001B[1;38m \"%s\",\n" +
		"\u001B[1;31m  \"httpError\":{\n" +
		"\u001B[1;31m    \"hostname\":  \u001B[1;31m \"%s\",\n" +
		"\u001B[1;31m    \"ipAddr\":    \u001B[1;31m \"%s\",\n" +
		"\u001B[1;31m    \"pathToCode\":\u001B[1;31m %s ,\n" +
		"\u001B[1;31m    \"lineCode\":  \u001B[1;31m %s ,\n" +
		"\u001B[1;31m    \"statusCode\":\u001B[1;31m \"%d\",\n" +
		"\u001B[1;31m    \"methodHttp\":\u001B[1;31m \"%s\",\n" +
		"\u001B[1;31m    \"url\":       \u001B[1;31m \"%s\",\n" +
		"\u001B[1;31m    \"code\":      \u001B[1;31m \"%s\",\n" +
		"\u001B[1;31m    \"message\":   \u001B[1;31m \"%s\",\n" +
		"\u001B[1;31m    \"error\":     \u001B[1;31m \"%s\"\n" +
		"\u001B[1;31m }\n" +
		"\u001B[1;31m},\u001B[0m\n"

	httpErrorInJSON = "{\n" +
		" \"dataTime\": \"%s\",\n" +
		" \"httpError\":{\n" +
		"   \"hostname\":   \"%s\",\n" +
		"   \"ipAddr\":     \"%s\",\n" +
		"   \"pathToCode\": \"%s\",\n" +
		"   \"lineCode\":   \"%s\",\n" +
		"   \"statusCode\": \"%d\",\n" +
		"   \"methodHttp\": \"%s\",\n" +
		"   \"url\":        \"%s\",\n" +
		"   \"code\":       \"%s\",\n" +
		"   \"message\":    \"%s\",\n" +
		"   \"error\":      \"%s\"\n" +
		" }\n" +
		"},\n"
)

// CustomWLogg структура журнала логирования.
type CustomWLogg struct {
	logI    *log.Logger
	logE    *log.Logger
	logW    *log.Logger
	infoPC  *pkg.InfoPC
	logFile *os.File
}

// NewCustomWLogger создает и возвращает новый экземпляр структуры CustomWLogg.
func NewCustomWLogger() *CustomWLogg {
	infoPC := pkg.NewInfoPC()
	file, err := ensureLogFileExists(pathToLoggFile)
	if err != nil {
		log.Fatal(err)
	}

	logI := log.New(file, "", 0)
	logW := log.New(file, "", 0)
	logE := log.New(file, "", 0)

	return &CustomWLogg{logI: logI, logW: logW, logE: logE, infoPC: infoPC, logFile: file}
}

// LogI выводит информационное сообщение в консоль и записывает его в файл в формате JSON.
//
// Параметры:
//   - msg (string): Сообщение, которое нужно залогировать. В сообщении содержится код msg[:5] и описание msg[6:].
func (l *CustomWLogg) LogI(msg string) {
	fmt.Printf(infoInConsoleInJSON, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), msg[:5], msg[6:])
	l.logI.Printf(infoInJSON, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), msg[:5], msg[6:])
}

// LogW выводит предупреждающие сообщение в консоль и записывает его в файл в формате JSON.
//
// Параметры:
//   - msg (string): Сообщение о предупреждении, которое нужно залогировать. В сообщении содержится код msg[:5] и описание msg[6:].
//   - warn (error): Объект ошибки, содержащий дополнительную информацию об ошибке.
func (l *CustomWLogg) LogW(msg string, warn error) {
	fmt.Printf(warningInConsoleInJson, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), msg[:5], msg[6:], warn)
	l.logI.Printf(warningInJson, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), msg[:5], msg[6:], warn)
}

// LogE выводит сообщение об ошибке в консоль и записывает его в файл в формате JSON.
//
// Параметры:
//   - msg (string): Сообщение об ошибке, которое нужно залогировать. В сообщении содержится код msg[:5] и описание msg[6:].
//   - err (error): Объект ошибки, содержащий дополнительную информацию об ошибке.
func (l *CustomWLogg) LogE(msg string, err error) {
	fmt.Printf(errorInConsoleInJSON, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), msg[:5], msg[6:], err)
	l.logI.Printf(errorInJSON, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), msg[:5], msg[6:], err)
}

// LogHttpI выводит информационное сообщение HTTP в консоль и записывает его в файл в формате JSON.
//
// Параметры:
//   - statusCode (int): Код статуса HTTP-ответа, связанного с ошибкой.
//   - methodHttp (string): Метод HTTP запроса.
//   - url (string): Ссылка на HTTP.
//   - msg (string): Сообщение об ошибке, которое нужно залогировать. В сообщении содержится код msg[:5] и описание msg[6:].
func (l *CustomWLogg) LogHttpI(statusCode int, methodHttp, url, msg string) {
	fmt.Printf(httpInfoInConsoleInJSON, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), statusCode, methodHttp, url, msg[:5], msg[6:])
	l.logI.Printf(httpInfoInJSON, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), statusCode, methodHttp, url, msg[:5], msg[6:])
}

// LogHttpW выводит предупреждающие сообщение HTTP в консоль и записывает его в файл в формате JSON.
//
// Параметры:
//   - statusCode (int): Код статуса HTTP-ответа, связанного с ошибкой.
//   - methodHttp (string): Метод HTTP запроса.
//   - url (string): Ссылка на HTTP.
//   - msg (string): Сообщение об ошибке, которое нужно залогировать. В сообщении содержится код msg[:5] и описание msg[6:].
//   - warn (error): Объект ошибки, содержащий дополнительную информацию об ошибке.
func (l *CustomWLogg) LogHttpW(statusCode int, methodHttp, url, msg string, warn error) {
	fmt.Printf(httpWarningInConsoleInJSON, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), statusCode, methodHttp, url, msg[:5], msg[6:], warn)
	l.logI.Printf(httpWarningInJSON, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), statusCode, methodHttp, url, msg[:5], msg[6:], warn)
}

// LogHttpE выводит сообщение об ошибке HTTP в консоль и записывает его в файл в формате JSON.
//
// Параметры:
//   - statusCode (int): Код статуса HTTP-ответа, связанного с ошибкой.
//   - methodHttp (string): Метод HTTP запроса.
//   - msg (string): Сообщение об ошибке, которое нужно залогировать. В сообщении содержится код msg[:5] и описание msg[6:].
//   - url (string): Ссылка на HTTP.
//   - err (error): Объект ошибки, содержащий дополнительную информацию об ошибке.
func (l *CustomWLogg) LogHttpE(statusCode int, methodHttp, url, msg string, err error) {
	fmt.Printf(httpErrorInConsoleInJSON, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), statusCode, methodHttp, url, msg[:5], msg[6:], err)
	l.logI.Printf(httpErrorInJSON, systemDateTime, l.infoPC.HostName(), l.infoPC.IPAddr(), fileWithFuncAndLineNum(), fileWithLineNum(), statusCode, methodHttp, url, msg[:5], msg[6:], err)
}

// Close закрывает файл логирования.
func (l *CustomWLogg) Close() {
	if l.logFile != nil {
		if err := l.logFile.Close(); err != nil {
			log.Fatal(err)
			return
		}
	}
	l.LogI("Ресурсы освобождены")
}

// ensureLogFileExists проверяет путь к файлу логирования и открывает его.
func ensureLogFileExists(pathToLoggFile string) (*os.File, error) {
	if _, err := os.Stat(pathToLoggFile); os.IsNotExist(err) {
		return createLoggFile(pathToLoggFile)
	} else if err != nil {
		return nil, err
	}

	return openLoggFile(pathToLoggFile)
}

// createLoggFile создает файл логирования и добавляет начальную строку "[".
func createLoggFile(pathToLoggFile string) (*os.File, error) {
	file, err := os.OpenFile(pathToLoggFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	if _, err = file.WriteString("[\n"); err != nil {
		if err = file.Close(); err != nil {
			return nil, err
		}
	}

	return file, nil
}

// openLoggFile открывает существующий файл логирования.
func openLoggFile(pathToLoggFile string) (*os.File, error) {
	file, err := os.OpenFile(pathToLoggFile, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// skipNumOfStackFrame Количество кадров стека, которые необходимо пропустить перед записью на ПК, где 0 идентифицирует
// кадр для самих вызывающих абонентов, а 1 идентифицирует вызывающего абонента. Возвращает количество записей,
// записанных на компьютер.
const skipNumOfStackFrame = 3

// fileWithLineNum возвращает имя файла и строку номера текущего файла.
func fileWithLineNum() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(skipNumOfStackFrame, pc)
	frame, _ := runtime.CallersFrames(pc[:n]).Next()
	idxFile := strings.LastIndexByte(frame.File, '/')
	return frame.File[idxFile+1:] + ":" + strconv.FormatInt(int64(frame.Line), 10)
}

// fileWithFuncAndLineNum возвращает имя файла с функцией и числовой строкой текущего файла.
func fileWithFuncAndLineNum() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(skipNumOfStackFrame, pc)
	frame, _ := runtime.CallersFrames(pc[:n]).Next()
	idxFile := strings.LastIndexByte(strconv.Itoa(frame.Line), '/')

	return "[" + frame.Function + "] - " + frame.File[idxFile+1:] + ":" + strconv.FormatInt(int64(frame.Line), 10)
}
