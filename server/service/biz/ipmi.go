package biz

import (
	"os/exec"
	"strings"
)

// IPMIPowerAction 合法的电源动作
var ipmiPowerActions = map[string]string{
	"status": "status", // 查询状态
	"on":     "on",     // 开机
	"off":    "off",    // 强制关机
	"reset":  "reset",  // 硬重启
	"soft":   "soft",   // 软关机
}

// ipmiPower 通过 ipmitool 控制服务器电源（需运行环境安装 ipmitool）
func ipmiPower(ip, user, password, action string) (string, error) {
	act := ipmiPowerActions[action]
	if act == "" {
		// 未知动作直接透传，交由 ipmitool 校验
		act = action
	}
	ip = strings.TrimSpace(ip)
	user = strings.TrimSpace(user)
	cmd := exec.Command("ipmitool", "-I", "lanplus", "-H", ip, "-U", user, "-P", password, "power", act)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
