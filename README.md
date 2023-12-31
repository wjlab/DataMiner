## DataMiner  Instructions

English Document | [中文文档](https://github.com/wjlab/DataMiner/blob/master/README-zh-CN.md)

### 1. Function Overview

- Supports sampling of all database data, and specifies the number of samples to be taken.
- Supports sampling of a specified database, and specifies the number of samples to be taken.
- Supports sampling of a specified database table, and specifies the number of samples to be taken.
- Supports for  the capture of key sensitive content in the database, currently supports mailboxes, ID cards, cell phone numbers, passwords.
- Supports for custom regular expression matching on database content.
- Supports for socks5 proxy connection to remote databases.
- Supports for statistics on the amount of data in the database.
- Supports for batch connection to databases for information collection.
- Supports for outputs data in HTML and CSV formats.
- Currently supports Mysql, Mssql, Oracle, Mongodb  and Postgre database.

### 2. Function Command Description

- Command Parameter Description

  ```
  Command：
  Sampledata,Abbreviated command: SD              //Full database sampling module
  SampleSingleDatabase,Abbreviated command:SSD    //Specified database sampling module
  Overview,Abbreviated command: OV                //Database data volume statistics module
  SearchSensitiveData,Abbreviated command: SS     //Database sensitive data capture module
  SingleTable,Abbreviated command: ST             //Database single table sampling module
  Parameters：
  -T  databaseType           //Database type (mandatory, currently supports mysql、mssql、oracle and mongodb)
  -da 127.0.0.1:3306        //Database address (mandatory, unless using -f parameter to specify file as input)
  -du name                  //Database username (mandatory, unless using -f parameter to specify file as input)
  -dp passwd                //Database password (mandatory, unless using -f parameter to specify file as input)
  -dn databaseName          //Database name (mandatory for single database sampling function)
  -ds databaseSchema        //Database schema(optional for sampling a single table to specify schema in postgre)
  -pa 127.0.0.1:8080        //Proxy address (optional)
  -pu name                  //Proxy username(optional)
  -pp passwd                //Proxy password(optional)
  -n  1                     //Specify the number of sampling data, default is 3 (optional)
  -t 1           //Number of threads used by the sensitive data capture module, default is 5 (optional)
  -p Customized Regular Expression  //Sensitive data capture module custom regular matching parameters (optional)
  -tf filepath                      //Login to the Oracle database using TNS file (for Oracle database only)
  -WA                               //Login using Windows local authentication (for Mssql database only)
  -f data.txt   //Batch database information import file, with each line of the text file containing one database record. Text format：schema://user:password@host:port 
  e.g.：mysql://root:123321@127.0.0.1:3306
        mssql://sa:123321@127.0.0.1:1433
        oracle://system:123321@127.0.0.1:1521
        postgre://postgres:123321@127.0.0.1:5432
        mongo://admin:123321@127.0.0.1:27017
        mongo://admin:123321@127.0.0.1:27017?admin
        mongo://:@127.0.0.1:27017
  The last two entries above are respectively for MongoDB database, specifying the admin database login mode and the mode without user and password login.
  ```

- Sample Data Extraction Module

  ```
  //Specify the mysql database, connect to the database, and sample 2 items from each table
  DataMiner SD -T mysql -da 127.0.0.1:3306 -du name -dp passwd -n 2
  
  //Specify mssql database，using socks5 proxy to connect to the database，and sample 2 items from each table
  DataMiner SD -T mssql -da 127.0.0.1:1433 -du name -dp passwd -pa 127.0.0.1:8080 -pu name -pp passwd -n 2
  
  //Use file to import the database connection information to connect, and sample 2 items from each table
  DataMiner SD -f data.txt  -n 2
  
  //Use file to import the database connection information to connect using socks5 proxy, and sample 2 items from each table
  DataMiner SD -f data.txt -pa 127.0.0.1:8080 -pu name -pp passwd -n 2
  
  //Oracle database TNS method login mode and use the full database sampling module
  DataMiner SD -T oracle -du name -dp passwd -tf tnsnames.ora
  
  //MSSQL database local Windows authentication login mode and use full database sampling module
  DataMiner SD -T mssql -WA
  
  //Mongodb database without user password login mode using full database sampling module
  DataMiner SD -T mongo -da 127.0.0.1:27017
  
  //Mongodb database specified 'admin' database login mode using full database sampling module
  DataMiner SD -T mongo -da 127.0.0.1:27017?admin -du name -dp password
  ```

- Sample a Single Database Module

  ```
  //Specify postgre database, meanwhile specify the 'test' database, and sample 2 items from all the tables in this database
  DataMiner SSD -T postgre -da 127.0.0.1:5432 -du name -dp passwd -dn test -n 2
  
  //Specify postgre database, connect to the database using socks5 proxy, meanwhile specify the 'test' database, and sample 2 items from all the tables in this database
  DataMiner SSD -T postgre -da 127.0.0.1:5432 -du name -dp passwd -pa 127.0.0.1:8080 -pu name -pp passwd -dn test -n 2
  ```

- Database Overview Module

  ```
  //Specify oracle database, connect to the database, and use database data volume statistics module
  DataMiner OV -T oracle -da 127.0.0.1:1521 -du name -dp passwd
  
  //Specify mysql database, using socks5 proxy to connect to the databse, and use database data volume statistics module
  DataMiner OV -T mysql -da 127.0.0.1:3306 -du name -dp passwd -pa 127.0.0.1:8080 -pu name -pp passwd
  
  //Use file to import the database connection information to connect, and use database data volume statistics module
  DataMiner OV -f data.txt
  
  //Use file to import the database connection information to connect using socks5 proxy, and use database data volume statistics module
  DataMiner OV -f data.txt -pa 127.0.0.1:8080 -pu name -pp passwd
  
  //Oracle database TNS method login mode and use database data volume statistics module
  DataMiner OV -T oracle -du name -dp passwd -tf tnsnames.ora
  
  //MSSQL database local Windows authentication login mode and use database data volume statistics module
  DataMiner OV -T mssql -WA
  
  //Mongodb database without user password login mode and use database data volume statistics module
  DataMiner OV -T mongo -da 127.0.0.1:27017
  
  //Mongodb database specified 'admin' database login mode and use database data volume statistics module
  DataMiner OV -T mongo -da 127.0.0.1:27017?admin -du name -dp password
  ```

- Key Sensitive Information Searching Module

  ```
  //Specify mssql database, connect to the database, sample 2 items from each table to find the key sensitive informaiton, and specify the use of 6 threads
  DataMiner SS -T mssql -da 127.0.0.1:1433 -du name -dp passwd -n 2 -t 6
  
  //Specify mysql database, connect to the database using socks5 proxy, sample 2 items from each table to find the key sensitive informaiton, and specify the use of 6 threads
  DataMiner SS -T mysql -da 127.0.0.1:3306 -du name -dp passwd -pa 127.0.0.1:8080 -pu name -pp passwd -n 2 -t 6
  
  //Use file to import the database connection information to connect, sample 2 items from each table to find the key sensitive informaiton, and specify the use of 6 threads
  DataMiner SS -f data.txt  -n 2 -t 6
  
  //Use file to import the database connection information to connect using socks5 proxy, sample 2 items from each table to find the key sensitive informaiton, and specify the use of 6 threads
  DataMiner SS -f data.txt -pa 127.0.0.1:8080 -pu name -pp passwd -n 2 -t 6
  
  //Specify mysql database, connect to the database, sample 2 items from each table to find the key sensitive informaiton using customized regular expression to match user name, and specify the use of 6 threads
  DataMiner SS -T mysql -da 127.0.0.1:3306 -du name -dp passwd -n 2 -t 6 -p ^[\x{4e00}-\x{9fa5}]{2,4}$
  
  //Oracle database TNS method login mode and use key sensitive information searching module
  DataMiner SS -T oracle -du name -dp passwd -tf tnsnames.ora
  
  //MSSQL database local Windows authentication login mode and use key sensitive information searching module
  DataMiner SS -T mssql -WA
  
  //Mongodb database without user password login mode and use key sensitive information searching module
  DataMiner SS -T mongo -da 127.0.0.1:27017
  
  //Mongodb database specified 'admin' database login mode and use key sensitive information searching module
  DataMiner SS -T mongo -da 127.0.0.1:27017?admin -du name -dp password
  ```


- Specify Single Table Sampling Function

  ```
  //Specify mysql database, connect to the database, specify the 'users' table in the 'test' database, and sample 2 items from this table
  DataMiner ST -T mysql -da 127.0.0.1:3306 -du name -dp passwd -n 2 -dt test.users
  
  //Specify mysql database, connect to the database using socks5 proxy, specify the 'users' table in the 'test' database, and sample 2 items from this table
  DataMiner ST -T mysql -da 127.0.0.1:3306 -du name -dp passwd -pa 127.0.0.1:8080 -pu name -pp passwd -n 2 -dt test.users
  
  //Specify postgre database, specify the 'users' table under the 'other' schema of the 'test' database, and set the number of sampled entries to 2 (if the 'ds' parameter is not specified, the schema defaults to 'public')
  DataMiner ST -T postgre -da 127.0.0.1:5432 -du name -dp passwd -n 2 -dt test.users -ds other
  ```

- Sampling Module HTML Output Example

  ![](https://github.com/wjlab/DataMiner/blob/master/image/HtmlOutput.png)

- Sampling Module CSV Output Example

  ![](https://github.com/wjlab/DataMiner/blob/master/image/CsvOutput.png)

- Database Overview Module HTML Output Example

  ![](https://github.com/wjlab/DataMiner/blob/master/image/Overview.png)

- Sensitive Data Capture Module  CSV Output Example

  ![](https://github.com/wjlab/DataMiner/blob/master/image/Secret.png)

