### orm支持的数据库类型

```
mysql | sqlite | oracle | pgsql | TiDB
```

### 支持多数据源 

```
mysqlkeymultiple=default::QRQM,slave::WEBAPP
mysqlkey = QRQM
QRQM.Host = 127.0.0.1
QRQM.Port = 3306
QRQM.DbName = nicetuan_recommend_qrqm
QRQM.UserName = root
QRQM.Password = 123456
QRQM.DriverName = mysql

slave = WEBAPP
WEBAPP.Host = 127.0.0.1
WEBAPP.Port = 3306
WEBAPP.DbName = nicetuan_weapp_db_new
WEBAPP.UserName = root
WEBAPP.Password = 123456
WEBAPP.DriverName = mysql
```