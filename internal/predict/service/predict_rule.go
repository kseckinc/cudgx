package service

import (
	"strings"
	"time"

	"github.com/galaxy-future/cudgx/internal/predict/model"
	"github.com/galaxy-future/cudgx/internal/request"
)

func CreatePredictRule(req *request.CreatePredictRuleRequest) error {
	predictRule := &model.PredictRule{
		Id:               0,
		Name:             req.Name,
		ServiceName:      req.ServiceName,
		ClusterName:      req.ClusterName,
		MetricName:       strings.ToLower(req.MetricName),
		BenchmarkQps:     req.BenchmarkQps,
		MinRedundancy:    req.MinRedundancy,
		MaxRedundancy:    req.MaxRedundancy,
		MinInstanceCount: req.MinInstanceCount,
		MaxInstanceCount: req.MaxInstanceCount,
		ExecuteRatio:     req.ExecuteRatio,
		Status:           req.Status,
		CreatedTime:      time.Now().Unix(),
	}
	if err := model.CreatePredictRule(predictRule); err != nil {
		return err
	}
	return nil
}

func DeletePredictRuleById(req *request.BatchDeletePredictRuleRequest) error {
	if err := model.DeletePredictRuleById(req.Ids); err != nil {
		return err
	}
	return nil
}

func UpdatePredictRuleById(req *request.UpdatePredictRuleRequest) error {
	if _, err := model.GetPredictRuleById(req.Id); err != nil {
		return err
	}
	predictRule := &model.PredictRule{
		Id:               req.Id,
		Name:             req.Name,
		ServiceName:      req.ServiceName,
		ClusterName:      req.ClusterName,
		MetricName:       strings.ToLower(req.MetricName),
		BenchmarkQps:     req.BenchmarkQps,
		MinRedundancy:    req.MinRedundancy,
		MaxRedundancy:    req.MaxRedundancy,
		MinInstanceCount: req.MinInstanceCount,
		MaxInstanceCount: req.MaxInstanceCount,
		ExecuteRatio:     req.ExecuteRatio,
		Status:           req.Status,
	}
	if err := model.UpdatePredictRule(predictRule); err != nil {
		return err
	}
	return nil
}

func GetPredictRuleByServiceNameAndClusterName(serviceName, clusterName string) (*model.PredictRule, error) {
	predictRule, err := model.GetPredictRuleByServiceNameAndClusterName(serviceName, clusterName)
	if err != nil {
		return nil, err
	}
	return predictRule, nil
}

func GetPredictRuleById(id int64) (*model.PredictRule, error) {
	predictRule, err := model.GetPredictRuleById(id)
	if err != nil {
		return nil, err
	}
	return predictRule, nil
}

func ListPredictRules(serviceName, clusterName string, pageNumber int, pageSize int) ([]*model.PredictRule, int, error) {
	predictRules, total, err := model.ListPredictRules(serviceName, clusterName, pageNumber, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return predictRules, total, nil
}

func UpdatePredictRuleStatus(id int64, status string) error {
	if _, err := model.GetPredictRuleById(id); err != nil {
		return err
	}
	if err := model.UpdatePredictRuleStatusById(id, status); err != nil {
		return err
	}
	return nil
}
