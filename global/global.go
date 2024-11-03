package global

import (
	"runtime"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	OSS     OSSConfig `json:"oss" yaml:"oss"`
	DB      DB        `json:"db" yaml:"db"`
	Refresh Refresh   `json:"refresh" yaml:"refresh"`
	Query   Query     `json:"query" yaml:"query"`
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
	DSN          string `json:"dsn" yaml:"dsn"`
	MaxIdleConns int    `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns" yaml:"maxOpenConns"`
}

type Refresh struct {
	WorkerNum int `json:"workerNum" yaml:"workerNum"` // Number of concurrent workers if <=0 use cpu num
}

type Query struct {
	WorkerNum int `json:"workerNum" yaml:"workerNum"` // Number of concurrent workers if <=0 use cpu num
	Distance  int `json:"distance" yaml:"distance"`   // Distance threshold
}

var ZapLogger *zap.Logger
var Logger *zap.SugaredLogger
var AppConfig Config
var Db *gorm.DB
var OSSClient *oss.Client

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
	viper.SetDefault("query.workerNum", numCPU)
	viper.SetDefault("query.distance", 5)
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
	if AppConfig.Query.WorkerNum <= 0 {
		// Use the number of CPUs as the default number of workers
		AppConfig.Query.WorkerNum = numCPU
	}
	if AppConfig.Query.Distance <= 0 {
		AppConfig.Query.Distance = 0
	}
	return &AppConfig
}

func InitGlobal() {
	config := ReadConfig()
	db, err := gorm.Open(mysql.Open(config.DB.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(AppConfig.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.DB.MaxOpenConns)
	Db = db

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.OSS.AccessKey, config.OSS.SecretKey)).
		WithRegion(config.OSS.Region)

	client := oss.NewClient(cfg)
	OSSClient = client
}
