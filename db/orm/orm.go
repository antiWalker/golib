package orm

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/antiWalker/golib/common"
	"github.com/antiWalker/golib/nacos"
	"strconv"
	"strings"
	"time"
)

type dbConfig struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	Host            string `json:"host"`
	Port            int64  `json:"port"`
	DbName          string `json:"dbname"`
	DriverName      string `json:"driverName"` //mysql | sqlite | oracle | pgsql | TiDB
	SetMaxIdleConns int    `json:"set_max_idle_conns"`
	SetMaxOpenConns int    `json:"set_max_open_conns"`
}

type groups struct {
	Name string
	Key  string
}

func init() {
	configMap := nacos.GetConfigMap()

	mysqlKey := configMap["mysqlkey"]
	orm.RegisterDriver("mysql", orm.DRMySQL)

	var dbGroup map[string]groups
	dbGroup = make(map[string]groups)

	if len(mysqlKey) != 0 {
		// 默认设置
		dbGroup["default"] = groups{"default", mysqlKey}
	}
	// 支持多组
	addMultiple(configMap, dbGroup)

	/*获取多个数据库配置*/
	configs := getDbConfigs(dbGroup, configMap)

	//统一注册
	register(configs)

}

func addMultiple(configMap map[string]string, dbGroup map[string]groups) {
	mysqlKeyMultiple := configMap["mysqlkeymultiple"]
	if len(mysqlKeyMultiple) > 0 {
		multiple := strings.Split(mysqlKeyMultiple, ",")
		for _, atom := range multiple {
			group := strings.Split(atom, "::")
			if len(group) != 2 || len(group[0]) == 0 || len(group[1]) == 0 {
				panic("database.conf mysqlkeymultiple format error,for example:  mysqlkeymultiple = name1::key1,name2::key2")
			}

			dbGroup[group[0]] = groups{group[0], group[1]}
		}
	}
}

func getDbConfigs(dbGroup map[string]groups, configMap map[string]string) map[string]dbConfig {
	var configs map[string]dbConfig
	configs = make(map[string]dbConfig)
	for _, group := range dbGroup {
		var config dbConfig
		config.DbName = configMap[group.Key+".DbName"]
		config.UserName = configMap[group.Key+".UserName"]
		config.Host = configMap[group.Key+".Host"]
		port, _ := strconv.Atoi(configMap[group.Key+".Port"])
		config.Port = int64(port)
		config.Password = configMap[group.Key+".Password"]
		config.DriverName = configMap[group.Key+".DriverName"]
		SetMaxIdleConns, err := strconv.Atoi(configMap[group.Key+".SetMaxIdleConns"])
		if err != nil {
			common.ErrorLogger.Info("orm getDbConfigs SetMaxIdleConns Error:", err)
		}
		config.SetMaxIdleConns = SetMaxIdleConns
		SetMaxOpenConns, err := strconv.Atoi(configMap[group.Key+".SetMaxOpenConns"])
		if err != nil {
			common.ErrorLogger.Info("orm getDbConfigs SetMaxOpenConns Error:", err)
		}
		config.SetMaxOpenConns = SetMaxOpenConns
		configs[group.Name] = config
	}

	return configs
}

func register(configs map[string]dbConfig) {
	var (
		SetMaxIdleConns, SetMaxOpenConns = 50, 50
	)
	for name, config := range configs {
		if config.SetMaxIdleConns == 0 && config.SetMaxOpenConns == 0 {
			config.SetMaxIdleConns, config.SetMaxOpenConns = SetMaxIdleConns, SetMaxOpenConns
		} else if config.SetMaxIdleConns == 0 {
			config.SetMaxIdleConns = SetMaxIdleConns
		} else if config.SetMaxOpenConns == 0 {
			config.SetMaxOpenConns = SetMaxOpenConns
		}
		orm.RegisterDataBase(name, config.DriverName, config.UserName+":"+config.Password+
			"@tcp("+config.Host+":"+strconv.FormatInt(config.Port, 10)+")/"+config.DbName, config.SetMaxIdleConns, config.SetMaxOpenConns) // SetMaxIdleConns,SetMaxOpenConns
		//设置连接超时为25s，比目前测试、生产的30s小即可，防止连接池可用连接不可用
		db, _ := orm.GetDB(name)
		db.SetConnMaxLifetime(25 * time.Second)
	}
}
