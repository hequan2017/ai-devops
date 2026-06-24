package initialize

import (
	_ "ai-devops/server/source/example"
	_ "ai-devops/server/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
