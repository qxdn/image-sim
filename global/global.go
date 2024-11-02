package global

import (
	"runtime"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	OSS     OSSConfig `json:"oss" yaml:"oss"`
	DB      DB        `json:"db" yaml:"db"`
	Refresh Refresh   `json:"refresh" yaml:"refresh"`
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

type Refresh struct {
	WorkerNum int `json:"workerNum" yaml:"workerNum"` // Number of concurrent workers if <=0 use cpu num
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
	numCPU := runtime.NumCPU()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetDefault("refresh.workerNum", numCPU)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		panic(err)
	}
	if AppConfig.Refresh.WorkerNum <= 0 {
		// Use the number of CPUs as the default number of workers
		AppConfig.Refresh.WorkerNum = numCPU
	}
	return &AppConfig
}
