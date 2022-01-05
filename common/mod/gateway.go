package mod

const GatewayModuleName = "monitoring/gateway"
const GatewayStatusSuccess = "success"

//GatewayPingResult 是Gateway探活接口，当程序启动时，用户可以 GET /ping接口
type GatewayPingResult struct {
	Status string `json:"status"`
	Module string `json:"module"`
}
