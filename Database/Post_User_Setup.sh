rm -r /usr/local/pgsql/data

mkdir -p /usr/local/pgsql/data
chown postgres /usr/local/pgsql/data

su - postgres -c "  /usr/local/pgsql/bin/initdb -D /usr/local/pgsql/data; 
                    /usr/local/pgsql/bin/pg_ctl -D /usr/local/var/postgres stop;
                    /usr/local/pgsql/bin/pg_ctl -D /usr/local/pgsql/data -l logfile start; 
                    /usr/local/pgsql/bin/createdb telemetrydb;
                    /usr/local/pgsql/bin/psql telemetrydb -c '  CREATE TABLE telemetry ( 
                                                                    time_step time, 
                                                                    tire_temps float ARRAY[4], 
                                                                    velocity float, 
                                                                    acceleration float, 
                                                                    location point
                                                                );
                                                                INSERT INTO telemetry (velocity, acceleration)
                                                                    VALUES (22.0, 0.0);
                                                                SELECT * FROM telemetry;'"

#Scuffed as hell but works sometimes :)   (works when postgres server has not been started already but if its first time initalizing it shouldnt be a problem)