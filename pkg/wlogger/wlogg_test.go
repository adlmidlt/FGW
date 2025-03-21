package wlogger

import (
	"FGW/pkg"
	"fmt"
	"log"
	"os"
	"testing"
)

func createTempFile(t *testing.T) *os.File {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "tmp_fgw.json")
	if err != nil {
		t.Fatal(err)
	}

	defer tmpFile.Close()

	return tmpFile
}

func existFile(t *testing.T, fileName string) bool {
	t.Helper()
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		fmt.Printf("file %v does not exist", fileName)
		return false
	}

	return true
}

func Test_openLoggFile(t *testing.T) {
	type args struct {
		pathToLoggFile string
	}

	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		{
			name: "Открыть файл",
			args: args{
				pathToLoggFile: createTempFile(t).Name(),
			},
			want:    &os.File{},
			wantErr: false,
		},
		{
			name: "Файл не найден",
			args: args{
				pathToLoggFile: "log/tmp_fgw.json",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openLoggFile(tt.args.pathToLoggFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("openLoggFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				defer got.Close()
				if got.Name() != tt.args.pathToLoggFile {
					t.Errorf("openLoggFile() got = %v, want %v", got.Name(), tt.args.pathToLoggFile)
				}
			}
		})
	}
}

func Test_createLoggFile(t *testing.T) {
	type args struct {
		pathToLoggFile string
	}
	tests := []struct {
		name    string
		args    args
		want    *os.File
		wantErr bool
	}{
		{
			name: "Создать файл",
			args: args{
				pathToLoggFile: createTempFile(t).Name(),
			},
			want:    &os.File{},
			wantErr: false,
		},
		{
			name: "Ошибка создания файла",
			args: args{
				pathToLoggFile: "/log/tmp_fgw.json",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Запись в файл",
			args: args{
				pathToLoggFile: createTempFile(t).Name(),
			},
			want:    &os.File{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createLoggFile(tt.args.pathToLoggFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("createLoggFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				defer got.Close()
				if got.Name() != tt.args.pathToLoggFile {
					t.Errorf("createLoggFile() got = %v, want %v", got.Name(), tt.args.pathToLoggFile)
				}

				content, err := os.ReadFile(tt.args.pathToLoggFile)
				if err != nil {
					t.Errorf("os.ReadFile() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if string(content) != "[\n" {
					t.Errorf("createLoggFile() got = %v, want %v", string(content), "[\n")
					return
				}
			}
		})
	}
}

func Test_ensureLogFileExists(t *testing.T) {
	type args struct {
		pathToLoggFile string
	}
	tests := []struct {
		name      string
		args      args
		pathExist bool
		want      *os.File
		wantErr   bool
	}{
		{
			name: "Создать файл",
			args: args{
				pathToLoggFile: createTempFile(t).Name(),
			},
			want:    &os.File{},
			wantErr: false,
		},
		{
			name: "Файл не существует",
			args: args{
				pathToLoggFile: "tmp_fgw.json",
			},
			pathExist: existFile(t, "tmp_fgw.json"),
			want:      nil,
		},
		{
			name: "Файл существует",
			args: args{
				pathToLoggFile: "tmp_fgw.json",
			},
			pathExist: existFile(t, createTempFile(t).Name()),
			want:      &os.File{},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ensureLogFileExists(tt.args.pathToLoggFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ensureLogFileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				defer got.Close()
				if got.Name() != tt.args.pathToLoggFile {
					t.Errorf("createLoggFile() got = %v, want %v", got.Name(), tt.args.pathToLoggFile)
				}
			}
		})
		os.Remove(tt.args.pathToLoggFile)
	}
}

func TestCustomWLogg_Close(t *testing.T) {
	tmpFle := createTempFile(t)
	type fields struct {
		logI    *log.Logger
		logE    *log.Logger
		logW    *log.Logger
		infoPC  *pkg.InfoPC
		logFile *os.File
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Закрыть файл журнала логирования, если он nil",
			fields: fields{
				logI:    log.New(tmpFle, "", 0),
				logE:    log.New(tmpFle, "", 0),
				logW:    log.New(tmpFle, "", 0),
				infoPC:  &pkg.InfoPC{},
				logFile: nil,
			},
		},
		{
			name: "Закрыть существующий журнал логирования",
			fields: fields{
				logI:    log.New(tmpFle, "", 0),
				logE:    log.New(tmpFle, "", 0),
				logW:    log.New(tmpFle, "", 0),
				infoPC:  &pkg.InfoPC{},
				logFile: tmpFle,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &CustomWLogg{
				logI:    tt.fields.logI,
				logE:    tt.fields.logE,
				logW:    tt.fields.logW,
				infoPC:  tt.fields.infoPC,
				logFile: tt.fields.logFile,
			}
			l.Close()
		})
		os.Remove(tmpFle.Name())
	}
}

func TestNewCustomWLogger(t *testing.T) {
	tests := []struct {
		name          string
		logFilePath   string
		expectError   bool
		expectLogFile bool
	}{
		{
			name:          "Успешная инициализация файла",
			logFilePath:   "fgw_log.json",
			expectError:   false,
			expectLogFile: true,
		},
		{
			name:          "Ошибка при создании файла",
			logFilePath:   "/log/fgw_log.json",
			expectError:   true,
			expectLogFile: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathToLoggFile = tt.logFilePath

			logger, err := NewCustomWLogger()

			if (err != nil) != tt.expectError {
				t.Errorf("NewCustomWLogger() error = %v, expectError %v", err, tt.expectError)
			}

			if err == nil {
				defer logger.Close()

				if logger.logI == nil || logger.logW == nil || logger.logE == nil {
					t.Fatal("Один из логгеров не был инициализирован")
				}

				if tt.expectLogFile {
					if _, err := os.Stat(tt.logFilePath); os.IsNotExist(err) {
						t.Fatalf("Файл лога %s не был создан", tt.logFilePath)
					}
				}

			}
		})
		os.Remove(tt.logFilePath)
	}
}
