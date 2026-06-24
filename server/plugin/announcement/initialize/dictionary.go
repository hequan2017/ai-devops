package initialize

import (
	"context"
	model "ai-devops/server/model/system"
	"ai-devops/server/plugin/plugin-tool/utils"
)

func Dictionary(ctx context.Context) {
	entities := []model.SysDictionary{}
	utils.RegisterDictionaries(entities...)
}
