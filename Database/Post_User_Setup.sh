mkdir -p /usr/local/pgsql/data
chown postgres /usr/local/pgsql/data

su - postgres
/usr/local/pgsql/bin/initdb -D /usr/local/pgsql/data
/usr/local/pgsql/bin/pg_ctl -D /usr/local/pgsql/data -l logfile start

PATH=/usr/local/pgsql/bin:$PATH
export PATH

createdb telemetrydb
psql telemetrydb

#for datatypes in postgresql https://www.postgresql.org/docs/current/datatype.html

CREATE TABLE telemetry {
    time_step       time,
    tire_temps      float ARRAY[4],
    velocity        float,
    acceleration    float,
    location        point
};

\q