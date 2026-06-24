package response

import "ai-devops/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
