username =  nghia
password = root
uri = 127.0.0.1
port = 3306
databaseName = repeatword_dev

repeatDir = ${dir ${shell which repeatword}}

up:
	mysql -u ${username} -proot < vob.sql
seed:
	mysql -u ${username} -proot < data.sql
functions:
	mysql -u ${username} -proot < functions.sql
clean: 
	mysql -u ${username} -proot < clean.sql
init: 
	mysql -u ${username} -h ${uri} -P ${port} -p${password}  < schema.sql
sqltest:
	mysql -u ${username} -proot < script.sql
uf: clean functions

release:
	go install -v -ldflags='-s -w'
	