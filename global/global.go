package global

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	OSS OSSConfig `json:"oss" yaml:"oss"`
	DB  DB        `json:"db" yaml:"db"`
}

type OSSConfig struct {
	Region     string `json:"region" yaml:"region"`
	BucketName string `json:"bucketName" yaml:"bucketName"`
	Directory  string `json:"directory" yaml:"directory"`
	AccessKey  string `json:"accessKey" yaml:"accessKey"`
	SecretKey  string `json:"secretKey" yaml:"secretKey"`
	CustomUrl  string `json:"customUrl" yaml:"customUrl"`
}

type DB struct {
	DSN string `json:"dsn" yaml:"dsn"`
}

var ZapLogger *zap.Logger
var Logger *zap.SugaredLogger
var AppConfig Config

func init() {
	ZapLogger, _ = zap.NewProduction()
	Logger = ZapLogger.Sugar()
}

/**
 * ReadConfig reads the configuration file and returns the configuration object
 */
func ReadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		panic(err)
	}
	return &AppConfig
}
