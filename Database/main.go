package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type telemetryPacket struct {
	date_entry time.Time
	time_step time.Time
	tire_temps [4][64][64]float64
	tire_pressures [4]float64
	velocity float64
	location [2]float64
	accelerator_input float64               
    brake_input float64                     
    steering_angle float64                  
    gyro_pitch float64                      
    gyro_yaw float64                        
    gyro_roll float64                       
    x_acceleration float64                
    y_acceleration float64                  
    z_acceleration float64                  
    total_power_draw float64                 
	active_suspension_power_draw [4]float64
	motor_power_draw [2]float64
	battery_voltage float64                  
	traction_loss float64
	abs_throttle_limiting [2]float64
	limited_slip_usage [2]float64
	log string                            
}

func NewConnection() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), "postgresql://FSAE_DB_User@localhost:5432/telemetrydb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	return dbpool
}

func InsertIntoPool(dbpool *pgxpool.Pool, data []telemetryPacket) {
	for _, packet := range data {
		_, err := dbpool.Exec(context.Background(), "insert into telemetry (date_entry, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_lost, abs_throttle_limiting, limited_slip_usage, log) values (CURRENT_DATE, NOW(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)", pgtype.Array([2]float64{1, 1}), )
		if err != nil {
			fmt.Fprintf(os.Stderr, "Insertion failed: %v\n", err)
			os.Exit(1)
		}
	}
}

func QueryFromPool(dbpool *pgxpool.Pool) *[]telemetryPacket {
	rows, err := dbpool.Query(context.Background(), "select date_entry, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_lost, abs_throttle_limiting, limited_slip_usage, log from telemetry")
	if err != nil {
		fmt.Fprintf(os.Stderr, "CollectRows error: %v\n", err)
		os.Exit(1)
	}

	var data []telemetryPacket

	for rows.Next() {
		var packet telemetryPacket
		err = rows.Scan(&packet.date_entry, &packet.time_step, &packet.tire_temps, &packet.tire_pressures, &packet.velocity, &packet.location, &packet.accelerator_input, &packet.brake_input, &packet.steering_angle, &packet.gyro_pitch, &packet.gyro_yaw, &packet.gyro_roll, &packet.x_acceleration, &packet.y_acceleration, &packet.z_acceleration, &packet.total_power_draw, &packet.active_suspension_power_draw, &packet.motor_power_draw, &packet.battery_voltage, &packet.traction_loss, &packet.abs_throttle_limiting, &packet.limited_slip_usage, &packet.log)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nGot Packet at TimeStep", packet.date_entry.GoString())
		data = append(data, packet)
	}

	return &data
}

func main() {
	var dbpool = NewConnection()

	var testData []telemetryPacket
	InsertIntoPool(dbpool, testData)

	var data = QueryFromPool(dbpool)
}