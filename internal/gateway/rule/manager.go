package rule

import (
	"github.com/galaxy-future/cudgx/common/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Manager struct {
	rules      map[string]map[string]*Rule
	connection *gorm.DB
}

func NewRuleManager(option *MysqlOption) (*Manager, error) {
	manager := &Manager{
		rules:      make(map[string]map[string]*Rule),
		connection: nil,
	}
	return manager, manager.InitManager(option)
}

func (manager *Manager) InitManager(option *MysqlOption) (err error) {

	connection, sqlErr := createConnection(option)
	if sqlErr != nil {
		err = sqlErr
		return
	}
	manager.connection = connection

	err = manager.refreshRules()
	if err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(time.Duration(option.RefreshSeconds) * time.Second)
		for {
			select {
			case <-ticker.C:
				err := manager.refreshRules()
				if err != nil {
					logger.GetLogger().Error("failed to refresh rules", zap.Error(err))
				}
			}
		}
	}()

	return err
}

func (manager *Manager) MatchRule(serviceName, metricName string) *Rule {

	serviceRulesMap, ok := manager.rules[metricName]
	if !ok {
		return nil
	}
	for serviceName != "" {
		rule, exists := serviceRulesMap[serviceName]
		if exists {
			return rule
		}
		tokenIndex := strings.LastIndex(serviceName, ".")
		if tokenIndex > 0 {
			serviceName = serviceName[0:tokenIndex]
		} else {
			return nil
		}
	}
	return nil
}

func (manager *Manager) refreshRules() error {
	rules, err := getRulesByMysql(manager.connection)
	if err != nil {
		return err
	}
	mapRules := make(map[string]map[string]*Rule)
	for _, rule := range rules {
		metricRules, mapExists := manager.rules[rule.MetricName]
		if mapExists {
			metricRules[rule.ServiceName] = &rule
		} else {
			metricRules = make(map[string]*Rule)
			metricRules[rule.ServiceName] = &rule
			mapRules[rule.MetricName] = metricRules
		}
	}
	manager.rules = mapRules
	return nil
}

func getRulesByMysql(connection *gorm.DB) ([]Rule, error) {
	var metricRules []Rule

	if err := connection.Find(&metricRules).Error; err != nil {
		return nil, err
	}

	return metricRules, nil
}

func createConnection(config *MysqlOption) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil

}
