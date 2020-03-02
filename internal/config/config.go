package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string   `json:"env"`
	DBConfig    DBConfig `json:"db"`
}

type DBConfig struct {
	Host         string       `json:"host"`
	Port         int          `json:"port"`
	TablesConfig TablesConfig `json:"tables_config"`
}

type TablesConfig struct {
	Discussions   TableConfig `json:"discussions"`
	Participants  TableConfig `json:"participants"`
	PostBookmarks TableConfig `json:"post_bookmarks"`
	Posts         TableConfig `json:"posts"`
	Users         TableConfig `json:"users"`
	Viewers       TableConfig `json:"viewers"`
}

type TableConfig struct {
	TableName string `json:"table_name"`
}

// addConfigDirectory used for testing to add the test config files directories.
func addConfigDirectory(dir string) {
	viper.AddConfigPath(dir)
}

func clearConfig() {
	viper.Reset()
}

func ReadConfig() (*Config, error) {
	env := os.Getenv("DELPHIS_ENV")
	if env == "" {
		env = "local"
	}
	viper.SetConfigName(fmt.Sprintf("%s.json", env))
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
