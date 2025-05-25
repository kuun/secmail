package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type SMTPConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	TLSPort int    `mapstructure:"tls_port"`
	TLS     struct {
		Enable   bool   `mapstructure:"enable"`
		CertFile string `mapstructure:"cert_file"`
		KeyFile  string `mapstructure:"key_file"`
	} `mapstructure:"tls"`
}

type Config struct {
	EmailDomain string         `mapstructure:"email_domain"`
	Database    DatabaseConfig `mapstructure:"database"`
	SMTP        SMTPConfig     `mapstructure:"smtp"`
}

var GlobalConfig Config

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

func InitConfig() error {
	viper.SetConfigName("secmail")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./etc")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return err
	}

	return nil
}
