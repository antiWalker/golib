package global

import "flag"

var NacosDir = flag.String("nacosDir", "./conf", "nacos配置文件目录,默认是./conf 目录，")
var NacosFileName = flag.String("filename", "nacos", "nacos配置文件名称")
var LogsDir = flag.String("logDir", "./logs", "日志路径")
