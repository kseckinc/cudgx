package clients

import (
	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/internal/predict/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBClient *gorm.DB

func InitDBClient(config *config.Database) error {
	db, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		logger.GetLogger().Error("InitDBClient err", zap.Error(err))
		return err
	}
	DBClient = db
	return nil
}
