package global

import "go.uber.org/zap"

var ZapLogger *zap.Logger
var Logger *zap.SugaredLogger

func init() {
	ZapLogger, _ = zap.NewProduction()
	Logger = ZapLogger.Sugar()
}
