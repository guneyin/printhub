package utils

type (
	ConfigParam string
)

const (
	ConfigParamDisk   ConfigParam = "disk"
	ConfigParamTenant ConfigParam = "tenant"

	ConfigParamOAuth ConfigParam = "oauth"
	ConfigParamToken ConfigParam = "token"
)

func (c ConfigParam) String() string {
	return string(c)
}
