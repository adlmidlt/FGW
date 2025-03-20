package wlogger

import (
	"os"
	"testing"
)

func createTempFile(t *testing.T) *os.File {
	tmpFile, err := os.CreateTemp("", "tmp_fgw.json")
	if err != nil {
		t.Fatal(err)
	}

	defer tmpFile.Close()

	return tmpFile
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
				pathToLoggFile: "tmp_fgw.json",
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
