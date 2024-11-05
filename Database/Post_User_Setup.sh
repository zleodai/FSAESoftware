rm -r /usr/local/pgsql/data

mkdir -p /usr/local/pgsql/data
chown postgres /usr/local/pgsql/data

su - postgres -c "  /usr/local/pgsql/bin/initdb -D /usr/local/pgsql/data;"

cp -a Create_Database.sql createDB.sql 
mv createDB.sql /usr/local/pgsql/data

su - postgres -c "  /usr/local/pgsql/bin/pg_ctl -D /usr/local/pgsql/data -l logfile start; 
                    /usr/local/pgsql/bin/createdb telemetrydb;
                    /usr/local/pgsql/bin/psql -a -f /usr/local/pgsql/data/createDB.sql;
                    /usr/local/pgsql/bin/pg_ctl stop -D /usr/local/pgsql/data"


#Scuffed as hell but works sometimes :)   (works when postgres server has not been started already but if its first time initalizing it shouldnt be a problem)