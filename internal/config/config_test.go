package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"
)

func Test_ReadConfig(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get the current running file location")
	}
	dirName := path.Dir(filename)
	type args struct {
		filename string
	}

	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "well_formed",
			args: args{
				filename: "well_formed",
			},
			want: func() *Config {
				config := Config{}
				fileContents, err := ioutil.ReadFile(path.Join(dirName, "test_config", "well_formed.json"))
				if err != nil {
					t.Fatalf("Failed to read contents fo file: %v", err)
				}
				err = json.Unmarshal(fileContents, &config)
				if err != nil {
					t.Fatalf("failed unmarshaling contents: %v", err)
				}
				return &config
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearConfig()
			AddConfigDirectory(path.Join(dirName, "test_config"))
			os.Setenv("DELPHIS_UNITTEST", "1")
			os.Setenv("DELPHIS_ENV", "well_formed")

			got, err := ReadConfig()
			if err != nil {
				t.Errorf("Failed reading config with error: %w", err)
				return
			}

			if got == nil || *got != *tt.want {
				t.Errorf("ReadConfig = %+v, \nwant: %+v", got, tt.want)
			}
		})
	}
}
