package model

import (
	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/internal/clients"
	"go.uber.org/zap"
)

type PredictRule struct {
	Id               int64  `json:"id"`
	Name             string `json:"name"`
	ServiceName      string `json:"service_name"`
	ClusterName      string `json:"cluster_name"`
	MetricName       string `json:"metric_name"`
	BenchmarkQps     int    `json:"benchmark_qps"`
	MinRedundancy    int    `json:"min_redundancy"`
	MaxRedundancy    int    `json:"max_redundancy"`
	MinInstanceCount int    `json:"min_instance_count"`
	MaxInstanceCount int    `json:"max_instance_count"`
	ExecuteRatio     int    `json:"execute_ratio"`
	Status           string `json:"status"`
	CreatedTime      int64  `json:"created_time"`
}

func (PredictRule) TableName() string {
	return "predict_rules"
}

func CreatePredictRule(predictRule *PredictRule) error {
	if err := clients.DBClient.Create(predictRule).Error; err != nil {
		logger.GetLogger().Error("CreatePredictRule from db", zap.Error(err))
		return err
	}
	return nil
}

func DeletePredictRuleById(ids []int64) error {
	if err := clients.DBClient.Delete(&PredictRule{}, ids).Error; err != nil {
		logger.GetLogger().Error("DeletePredictRuleById from db", zap.Error(err))
		return err
	}
	return nil
}

func UpdatePredictRule(predictRule *PredictRule) error {
	updateMap := map[string]interface{}{
		"name":               predictRule.Name,
		"service_name":       predictRule.ServiceName,
		"cluster_name":       predictRule.ClusterName,
		"metric_name":        predictRule.MetricName,
		"benchmark_qps":      predictRule.BenchmarkQps,
		"min_redundancy":     predictRule.MinRedundancy,
		"max_redundancy":     predictRule.MaxRedundancy,
		"min_instance_count": predictRule.MinInstanceCount,
		"max_instance_count": predictRule.MaxInstanceCount,
		"execute_ratio":      predictRule.ExecuteRatio,
		"status":             predictRule.Status,
	}
	if err := clients.DBClient.Model(&PredictRule{}).Where("id", predictRule.Id).Updates(updateMap).Error; err != nil {
		logger.GetLogger().Error("UpdatePredictRule from db", zap.Error(err))
		return err
	}
	return nil
}

func GetPredictRuleById(id int64) (*PredictRule, error) {
	var predictRule PredictRule
	if err := clients.DBClient.Where("id = ?", id).First(&predictRule).Error; err != nil {
		logger.GetLogger().Error("GetPredictRuleById from db", zap.Error(err))
		return nil, err
	}
	return &predictRule, nil
}

func GetPredictRuleByServiceNameAndClusterName(serviceName, clusterName string) (*PredictRule, error) {
	var predictRule PredictRule
	if err := clients.DBClient.Where("service_name = ? and cluster_name = ? ", serviceName, clusterName).First(&predictRule).Error; err != nil {
		logger.GetLogger().Error("GetPredictRuleByServiceNameAndClusterName from db", zap.Error(err))
		return nil, err
	}
	return &predictRule, nil
}

func ListPredictRules(serviceName, clusterName string, pageNumber int, pageSize int) ([]*PredictRule, int, error) {
	theClient := clients.DBClient.Model(&PredictRule{}).Where("service_name = ?", serviceName)
	if clusterName != "" {
		theClient.Where("cluster_name = ?", clusterName)
	}
	var total int64
	if err := theClient.Count(&total).Error; err != nil {
		logger.GetLogger().Error("ListPredictRules from db", zap.Error(err))
		return nil, 0, err
	}
	var predictRules []*PredictRule
	if err := theClient.Order("id desc").Offset((pageNumber - 1) * pageSize).Limit(pageSize).Find(&predictRules).Error; err != nil {
		logger.GetLogger().Error("ListPredictRules from db", zap.Error(err))
		return nil, 0, err
	}
	return predictRules, int(total), nil
}

func ListAllPredictRules() ([]*PredictRule, error) {
	theClient := clients.DBClient.Model(&PredictRule{})
	var predictRules []*PredictRule
	if err := theClient.Find(&predictRules).Error; err != nil {
		logger.GetLogger().Error("ListAllPredictRules from db", zap.Error(err))
		return nil, err
	}
	return predictRules, nil
}

func UpdatePredictRuleStatusById(id int64, status string) error {
	if err := clients.DBClient.Model(&PredictRule{}).Where("id", id).Update("status", status).Error; err != nil {
		logger.GetLogger().Error("UpdatePredictRuleStatusById from write db", zap.Error(err))
		return err
	}
	return nil
}
