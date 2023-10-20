## 数据库自动取样工具DataMiner使用说明

中文文档 | [English Document](https://github.com/wjlab/DataMiner/blob/master/README.md)

### 1.功能概览

- 支持对数据库全部数据进行取样，并支持指定条数取样
- 支持对数据库指定单一数据库进行取样，并支持指定条数取样
- 支持对数据库指定数据表进行取样，并支持指定条数取样
- 支持对数据库中的关键敏感内容进行捕获，目前支持邮箱、身份证、手机号、密码
- 支持对数据库内容进行自定义正则表达式进行匹配
- 支持socks5代理连接远程数据库
- 支持对数据库中的数据量进行统计
- 支持批量连接数据库进行信息收集
- 数据输出为HTML和CSV
- 目前支持Mysql、Mssql、Oracle、Mongodb和Postgre数据库

### 2.功能命令说明

- 命令参数说明

  ```
  命令：
  Sampledata,缩减命令: SD              //数据库全部取样功能
  SampleSingleDatabase,缩减命令：SSD   //指定单一数据库取样功能
  Overview,缩减命令: OV                //数据库数据量统计功能
  SearchSensitiveData,缩减命令: SS     //数据库敏感数据捕获功能
  SingleTable,缩减命令: ST             //数据库单表取样功能
  参数：
  -T  databaseType                    //数据库类型(必选参数，目前支持 mysql、mssql、oracle、mongodb)
  -da 127.0.0.1:3306                  //数据库地址(必选参数，除非使用-f参数文件输入数据)
  -du name                            //数据库用户名(必选参数，除非使用-f参数文件输入数据)
  -dp passwd                          //数据库密码(必选参数，除非使用-f参数文件输入数据)
  -dn databaseName                    //指定数据库名称(单一数据库取样功能必选参数)
  -ds databaseSchema                  //指定数据库schema(可选参数，postgre数据库单表取样时可指定schema)
  -pa 127.0.0.1:8080                  //代理地址(可选参数)
  -pu name                            //代理用户名(可选参数)
  -pp passwd                          //代理密码(可选参数)
  -n  1                               //指定取样数据条数，默认为3(可选参数)
  -t 1                                //数据库敏感数据捕获功能使用线程数量，默认为5(可选参数)
  -p 自定义正则表达式                  //数据库敏感数据捕获功能自定义正则匹配参数(可选参数)
  -tf filepath                        //使用TNS方式登录Oracle数据库(仅对于Oracle数据库)
  -WA                                 //使用Windows本地认证方式登录(仅针对于mssql数据库)
  -f data.txt                         //批量数据库信息导入文件，文本中一条数据库信息占用一行
                                      文本格式：schema://user:password@host:port 
                                      如：mysql://root:123321@127.0.0.1:3306
                                          mssql://sa:123321@127.0.0.1:1433
                                          oracle://system:123321@127.0.0.1:1521
                                          postgre://postgres:123321@127.0.0.1:5432
                                          mongo://admin:123321@127.0.0.1:27017
                                          mongo://admin:123321@127.0.0.1:27017?admin
                                          mongo://:@127.0.0.1:27017
                                          上述后两条分别为mongodb数据库 指定admin数据库登录模式与无用户密码登录模式
  ```

- 全部数据库取样功能

  ```
  //指定mysql数据库，连接数据库，每个表中内容取样条数为2
  DataMiner SD -T mysql -da 127.0.0.1:3306 -du name -dp passwd -n 2
  //指定mssql数据库，使用socks代理连接数据库，每个表中内容取样条数为2
  DataMiner SD -T mssql -da 127.0.0.1:1433 -du name -dp passwd -pa 127.0.0.1:8080 -pu name -pp passwd -n 2
  //使用文件批量导入数据库连接信息进行连接，每个表中内容取样条数为2
  DataMiner SD -f data.txt  -n 2
  //使用文件批量导入数据库连接信息并使用socks代理进行连接，每个表中内容取样条数为2
  DataMiner SD -f data.txt -pa 127.0.0.1:8080 -pu name -pp passwd -n 2
  //Oracle数据库TNS方式登录使用全部数据库取样功能
  DataMiner SD -T oracle -du name -dp passwd -tf tnsnames.ora
  //MSSQL数据库本地Windows认证登录使用全部数据库取样功能
  DataMiner SD -T mssql -WA
  //Mongodb数据库无用户密码登录模式使用全部数据库取样功能
  DataMiner SD -T mongo -da 127.0.0.1:27017
  //Mongodb数据库指定admin数据库登录模式使用全部数据库取样功能
  DataMiner SD -T mongo -da 127.0.0.1:27017?admin -du name -dp password
  ```

- 指定单一数据库取样功能

  ```
  //指定postgre数据库,连接数据库，对名为'test'的数据库下所有表进行取样，取样条数为2
  DataMiner ST -T postgre -da 127.0.0.1:5432 -du name -dp passwd -dn test -n 2
  //指定postgre数据库,使用socks代理连接数据库，对名为'test'的数据库下所有表进行取样，取样条数为2
  DataMiner ST -T postgre -da 127.0.0.1:5432 -du name -dp passwd -pa 127.0.0.1:8080 -pu name -pp passwd -n 2 -dn test
  ```
  
- 数据量统计概览功能

  ```
  //指定oracle数据库，连接数据库，使用数据量统计命令
  DataMiner OV -T oracle -da 127.0.0.1:1521 -du name -dp passwd
  //指定mysql数据库,使用socks代理连接数据库，使用数据量统计命令
  DataMiner OV -T mysql -da 127.0.0.1:3306 -du name -dp passwd -pa 127.0.0.1:8080 -pu name -pp passwd
  //使用文件批量导入数据库连接信息进行连接，使用数据量统计命令
  DataMiner OV -f data.txt
  //使用文件批量导入数据库连接信息并使用socks代理进行连接，使用数据量统计命令
  DataMiner OV -f data.txt -pa 127.0.0.1:8080 -pu name -pp passwd
  //Oracle数据库TNS方式登录使用数据量统计概览功能
  DataMiner OV -T oracle -du name -dp passwd -tf tnsnames.ora
  //MSSQL数据库本地Windows认证登录使用数据量统计概览功能
  DataMiner OV -T mssql -WA
  //Mongodb数据库无用户密码登录模式使用数据量统计概览功能
  DataMiner OV -T mongo -da 127.0.0.1:27017
  //Mongodb数据库指定admin数据库登录模式使用数据量统计概览功能
  DataMiner OV -T mongo -da 127.0.0.1:27017?admin -du name -dp password
  ```

- 关键敏感信息捕获功能

  ```
  //指定mssql数据库，连接数据库，每个表中内容取样条数为2,并指定使用6个线程
  DataMiner SS -T mssql -da 127.0.0.1:1433 -du name -dp passwd -n 2 -t 6
  //指定mysql数据库,使用socks代理连接数据库，每个表中内容取样条数为2，并指定使用6个线程
  DataMiner SS -T mysql -da 127.0.0.1:3306 -du name -dp passwd -pa 127.0.0.1:8080 -pu name -pp passwd -n 2 -t 6
  //使用文件批量导入数据库连接信息进行连接，每个表中内容取样条数为2,并指定使用6个线程
  DataMiner SS -f data.txt  -n 2 -t 6
  //使用文件批量导入数据库连接信息并使用socks代理进行连接，每个表中内容取样条数为2,并指定使用6个线程
  DataMiner SS -f data.txt -pa 127.0.0.1:8080 -pu name -pp passwd -n 2 -t 6
  //指定mysql数据库,连接数据库，每个表中内容取样条数为2,指定使用6个线程，并使用自定义正则匹配用户名
  DataMiner SS -T mysql -da 127.0.0.1:3306 -du name -dp passwd -n 2 -t 6 -p ^[\x{4e00}-\x{9fa5}]{2,4}$
  //Oracle数据库TNS方式登录使用关键敏感信息捕获功能
  DataMiner SS -T oracle -du name -dp passwd -tf tnsnames.ora
  //MSSQL数据库本地Windows认证登录使用关键敏感信息捕获功能
  DataMiner SS -T mssql -WA
  //Mongodb数据库无用户密码登录模式使用关键敏感信息捕获功能
  DataMiner SS -T mongo -da 127.0.0.1:27017
  //Mongodb数据库指定admin数据库登录模式使用关键敏感信息捕获功能
  DataMiner SS -T mongo -da 127.0.0.1:27017?admin -du name -dp password
  ```


- 指定数据库单表取样功能

  ```
  //指定mysql数据库,连接数据库，指定test数据库中users表，取样条数为2
  DataMiner ST -T mysql -da 127.0.0.1:3306 -du name -dp passwd -n 2 -dt test.users
  //指定mysql数据库,使用socks代理连接数据库，指定test数据库中users表，取样条数为2
  DataMiner ST -T mysql -da 127.0.0.1:3306 -du name -dp passwd -pa 127.0.0.1:8080 -pu name -pp passwd -n 2 -dt test.users
  //指定postgre数据库,指定test数据库'other'schema下的users表，取样条数为2 (若不指定ds参数，schema默认为public)
  DataMiner ST -T postgre -da 127.0.0.1:5432 -du name -dp passwd -n 2 -dt test.users -ds other
  ```

- 取样Sample模块HTML结果输出样例

  ![](https://github.com/wjlab/DataMiner/blob/master/image/HtmlOutput.png)

- 取样Sample模块CSV结果输出样例

  ![](https://github.com/wjlab/DataMiner/blob/master/image/CsvOutput.png)

- 数据量统计模块HTML结果输出样例

  ![](https://github.com/wjlab/DataMiner/blob/master/image/Overview.png)

- 敏感数据捕获模块 CSV结果输出样例

  ![](https://github.com/wjlab/DataMiner/blob/master/image/Secret.png)

