package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/cache"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"github.com/antiWalker/golib/common"
	"github.com/antiWalker/golib/global"
	"os"
	"strings"
	"time"
)

type NacosSetting struct {
	Ip          string
	Port        uint64
	NamespaceId string
	DataId      string
	Group       string
	LogDir      string
	CacheDir    string
	RotateTime  string
	LogLevel    string
	TimeoutMs   uint64
	MaxAge      int64
}

func ParseProperties(v *viper.Viper) NacosSetting {
	var nacos NacosSetting
	if err := v.Unmarshal(&nacos); err != nil {
		common.ErrorLogger.Infof("err:%v ", err)
	}
	return nacos
}

var ConfigClient config_client.IConfigClient
var Setting NacosSetting

func init() {
	for _, arg := range os.Args {
		if strings.Contains(arg, "nacosDir") {
			dir := strings.Split(arg, "=")[1]
			*global.NacosDir = dir
			continue
		}
		if strings.Contains(arg, "filename") {
			filename := strings.Split(arg, "=")[1]
			*global.NacosFileName = filename
			continue
		}
		if strings.Contains(arg, "logDir") {
			logDir := strings.Split(arg, "=")[1]
			*global.LogsDir = logDir
			continue
		}
	}
	v := viper.New()
	v.SetConfigName(*global.NacosFileName)
	v.SetConfigType("properties")
	v.AddConfigPath(*global.NacosDir)
	err := v.ReadInConfig()
	if err != nil {
		common.ErrorLogger.Infof("Fatal error config file: %v ", err)
	}
	Setting = ParseProperties(v)
	sc := []constant.ServerConfig{
		{
			IpAddr: Setting.Ip,
			Port:   Setting.Port,
		},
	}
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(Setting.NamespaceId),
		constant.WithTimeoutMs(Setting.TimeoutMs),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(*global.LogsDir),
		constant.WithCacheDir(*global.NacosDir),
		constant.WithRotateTime(Setting.RotateTime),
		constant.WithMaxAge(Setting.MaxAge),
		constant.WithLogLevel(Setting.LogLevel),
	)

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		common.ErrorLogger.Infof("Fatal error config file: %v  ", err)
	}
	ConfigClient = client

	common.InfoLogger.Infof(" nacos init end :%v", time.Now().String())
}

func GetConfigProxy() (string, error) {
	content, err := ConfigClient.GetConfig(vo.ConfigParam{
		DataId: Setting.DataId,
		Group:  Setting.Group,
	})
	return content, err
}

func GetConfigs() ([]string, error) {
	content, err := GetConfigsFromLocal()
	if err != nil {
		content, err = GetConfigProxy()
		return strings.Split(content, "\n"), err
	}
	return strings.Split(content, "\n"), err
}

func GetConfigsFromLocal() (string, error) {
	cacheKey := util.GetConfigCacheKey(Setting.DataId, Setting.Group, Setting.NamespaceId)
	return cache.ReadConfigFromFile(cacheKey, Setting.CacheDir+string(os.PathSeparator)+"config")
}

// GetConfigMap 获取配置文件，转成map格式
func GetConfigMap() map[string]string {
	var result = make(map[string]string)
	configs, err := GetConfigs()
	if err != nil {
		common.ErrorLogger.Infof("GetConfigs error : %v ", err)
	}
	for _, config := range configs {
		if config == "" {
			continue
		}
		indexOf := strings.Index(config, "=")
		length := len(config)
		result[strings.TrimSpace(config[0:indexOf])] = strings.TrimSpace(config[indexOf+1 : length])

	}
	return result
}
