#! /bin/bash

rm -r ./postgresServer
mkdir -p ./postgresServer
./FSAESoftwareBinaryImports/pgsql/bin/initdb -D ./postgresServer

./FSAESoftwareBinaryImports/pgsql/bin/pg_ctl -D ./postgresServer -l ./postgresServer/logfile.txt start 
./FSAESoftwareBinaryImports/pgsql/bin/createuser FSAE_DB_User -s
./FSAESoftwareBinaryImports/pgsql/bin/createdb telemetrydb
./FSAESoftwareBinaryImports/pgsql/bin/psql -h localhost -p 5432 -d telemetrydb -U FSAE_DB_User -a -f ./SQLCreateDB.sql
./FSAESoftwareBinaryImports/pgsql/bin/pg_ctl stop -D ./postgresServer