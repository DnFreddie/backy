package config

import (
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	EmailCreds struct {
		Email string `mapstructure:"email"`
		Passw string `mapstructure:"passw"`
	} `mapstructure:"email_creds"`
	CronTime    string `mapstructure:"cron_time"`
	ConfigPath  string `mapstructure:"config_path"`
	BackupDir   string `mapstructure:"backup_dir"`
	Bconfig_Dir string `mapstructure:"config_dir"`
}

func LoadConfig() {
	var config Configuration

	//Location
	viper.SetConfigName(".backy")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath("$HOME")

	//Default Values
	viper.SetDefault("cron_time", "@daily")
	viper.SetDefault("config_path", ".config")
	viper.SetDefault("backup_dir", "backups")
	viper.SetDefault("config_dir", "back_conf")

	err := viper.ReadInConfig()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

		} else {
			log.Fatal(err)
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
