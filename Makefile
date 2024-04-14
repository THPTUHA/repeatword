up:
	mysql -u nghia -proot < vob.sql
seed:
	mysql -u nghia -proot < data.sql
functions:
	mysql -u nghia -proot < functions.sql
clean: 
	mysql -u nghia -proot < clean.sql
init: 
	mysql -u nghia -proot < schema.sql
sqltest:
	mysql -u nghia -proot < script.sql
uf: clean functions