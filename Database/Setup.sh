#! /bin/bash

printf "\nStarting Installation\n\n"

export GIT_TRACE_PACKET=1
export GIT_TRACE=1
export GIT_CURL_VERBOSE=1
git config --global http.postBuffer 157286400

git clone https://github.com/zleodai/FSAESoftwareBinaryImports --depth 1;
cd pgsql/
git fetch --unshallow 

mkdir -p ./postgresServer
./FSAESoftwareBinaryImports/pgsql/bin/initdb -D ./postgresServer

./FSAESoftwareBinaryImports/pgsql/bin/pg_ctl -D ./postgresServer -l ./postgresServer/logfile.txt start 
./FSAESoftwareBinaryImports/pgsql/bin/createuser FSAE_DB_User -s
./FSAESoftwareBinaryImports/pgsql/bin/createdb telemetrydb
./FSAESoftwareBinaryImports/pgsql/bin/psql -h localhost -p 5432 -d telemetrydb -U FSAE_DB_User -a -f ./SQLCreateDB.sql
./FSAESoftwareBinaryImports/pgsql/bin/pg_ctl stop -D ./postgresServer

powershell ./SetupEnvVar.ps1