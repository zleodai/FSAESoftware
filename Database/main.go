package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type telemetryPacket struct {
	// date_entry string
	// time_step string
	// tire_temps [4][1][1]pgtype.Float8
	// tire_pressures [4]pgtype.Float8
	// velocity pgtype.Numeric
	// location [2]pgtype.Numeric
	// accelerator_input pgtype.Numeric               
    // brake_input pgtype.Numeric                     
    // steering_angle pgtype.Numeric                  
    // gyro_pitch pgtype.Numeric                      
    // gyro_yaw pgtype.Numeric                        
    // gyro_roll pgtype.Numeric                       
    x_acceleration float64                
    y_acceleration float64                  
    z_acceleration float64                  
    // total_power_draw pgtype.Float8                 
	// active_suspension_power_draw [4]pgtype.Float8
	// motor_power_draw [2]pgtype.Float8
	// battery_voltage pgtype.Float8                  
	// traction_loss pgtype.Float8
	// abs_throttle_limiting [2]pgtype.Float8
	// limited_slip_usage [2]pgtype.Float8
	// log string                            
}

func main() {
	dbpool, err := pgxpool.New(context.Background(), "postgresql://FSAE_DB_User@localhost:5432/telemetrydb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// for i := 0; i < 10; i++ {
	// 	_, err = dbpool.Exec(context.Background(), "insert into telemetry(date_entry, time_step, x_acceleration, y_acceleration, z_acceleration) values (CURRENT_DATE, NOW(), 0.0, 0.0, 0.0)")
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Insertion failed: %v\n", err)
	// 		os.Exit(1)
	// 	}
	// }

	for i := 0; i < 10; i++ {
		_, err = dbpool.Exec(context.Background(), "insert into telemetry(x_acceleration, y_acceleration, z_acceleration) values (0.0, 0.0, 0.0)")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Insertion failed: %v\n", err)
			os.Exit(1)
		}
	}

	rows, _ := dbpool.Query(context.Background(), "select (x_acceleration, y_acceleration, z_acceleration) from telemetry")
	packets, err := pgx.CollectRows(rows, pgx.RowToStructByName[telemetryPacket])
	if err != nil {
		fmt.Fprintf(os.Stderr, "CollectRows error: %v\n", err)
		os.Exit(1)
	}

	for _, packet := range packets {
		fmt.Printf("\nAcceleration:[%f, %f, %f]", packet.x_acceleration, packet.y_acceleration, packet.z_acceleration)
	}

	// for rows.Next() {
	// 	var date pgtype.Date
	// 	oErr := rows.Scan(&date)
	// 	if oErr != nil {
	// 		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", oErr)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println(date)
	// }	
}