package config

import (
	"regexp"
)

// IpReg 验证ip正则表达式
const IpReg = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"

// localhost正则表达式
const localhostReg = "^(localhost|127\\.0\\.0\\.1)$"

// PortReg 验证端口正则表达式
const PortReg = "^([0-9]|[1-9]\\d|[1-9]\\d{2}|[1-9]\\d{3}|[1-5]\\d{4}|6[0-4]\\d{3}|65[0-4]\\d{2}|655[0-2]\\d|6553[0-5])$"

// PositiveReg 验证正数正则表达式
const PositiveReg = "^[1-9]\\d*$"

func VerifyReg(reg, value string) bool {
	res, _ := regexp.Match(reg, []byte(value))
	return res
}
