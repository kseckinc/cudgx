package response

import "github.com/galaxy-future/cudgx/internal/predict/model"

type Pager struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
}

type ListPredictRuleResponse struct {
	PredictRuleList []*model.PredictRule `json:"predict_rule_list"`
	Pager           Pager                `json:"pager"`
}
