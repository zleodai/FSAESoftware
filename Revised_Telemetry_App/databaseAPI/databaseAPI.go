package databaseAPI

import (
	"context"
	"fmt"
	"math"
	"os"
	"time"

	"math/rand/v2"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MiniTelemetryPacket struct {
	Id			int64
	LapId 		int64
}

type TelemetryPacket struct {
	Id                           int64
	LapId						 int64
	Date_entry                   time.Time
	Time_step                    time.Time
	Tire_temps                   [4]float64
	Tire_pressures               [4]float64
	Velocity                     float64
	Location                     [2]float64
	Accelerator_input            float64
	Brake_input                  float64
	Steering_angle               float64
	Gyro_pitch                   float64
	Gyro_yaw                     float64
	Gyro_roll                    float64
	X_acceleration               float64
	Y_acceleration               float64
	Z_acceleration               float64
	Total_power_draw             float64
	Active_suspension_power_draw [4]float64
	Motor_power_draw             [2]float64
	Battery_voltage              float64
	Traction_loss                float64
	Abs_throttle_limiting        [2]float64
	Limited_slip_usage           [2]float64
}

func TempTelemetryPacket() TelemetryPacket {
	var packet TelemetryPacket = TelemetryPacket{
		LapId:						  rand.Int64N(5),
		Tire_temps:                   [4]float64{rand.Float64()*80 + 15, rand.Float64()*80 + 15, rand.Float64()*80 + 15, rand.Float64()*80 + 15},
		Tire_pressures:               [4]float64{rand.Float64()*5 + 25, rand.Float64()*5 + 25, rand.Float64()*5 + 25, rand.Float64()*5 + 25},
		Velocity:                     rand.Float64() * 100,
		Location:                     [2]float64{rand.Float64(), rand.Float64()},
		Accelerator_input:            rand.Float64(),
		Brake_input:                  rand.Float64(),
		Steering_angle:               rand.Float64(),
		Gyro_pitch:                   rand.Float64(),
		Gyro_yaw:                     rand.Float64(),
		Gyro_roll:                    rand.Float64(),
		X_acceleration:               rand.Float64(),
		Y_acceleration:               rand.Float64(),
		Z_acceleration:               rand.Float64(),
		Total_power_draw:             rand.Float64(),
		Active_suspension_power_draw: [4]float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()},
		Motor_power_draw:             [2]float64{rand.Float64()},
		Battery_voltage:              rand.Float64(),
		Traction_loss:                rand.Float64(),
		Abs_throttle_limiting:        [2]float64{rand.Float64(), rand.Float64()},
		Limited_slip_usage:           [2]float64{rand.Float64(), rand.Float64()},
	}
	return packet
}

var LargestID int64 = -404
var LatestID int64 = -404

func getLargestId(dbpool *pgxpool.Pool) int64 {
	if LargestID == -404 {
		var newLargestID int64 = 0
		var existingData = QueryFromPool(dbpool)
		for _, packet := range *existingData {
			newLargestID = int64(math.Max(float64(newLargestID), float64(packet.Id)))
		}

		LargestID = newLargestID
	}
	LargestID += 1
	return LargestID
}

func getLatestId(dbpool *pgxpool.Pool) int64 {
	var latestID int64 = 0
	var existingData = QueryFromPool(dbpool)
	for _, packet := range *existingData {
		latestID = int64(math.Max(float64(latestID), float64(packet.Id)))
	}

	return int64(latestID)
}

func NewConnection() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), "postgresql://FSAE_DB_User@localhost:5432/telemetrydb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	LatestID = getLatestId(dbpool)

	return dbpool
}

func CloseConnection(dbpool *pgxpool.Pool) {
	defer dbpool.Close()
}

func InsertIntoPool(dbpool *pgxpool.Pool, data []TelemetryPacket) {
	for _, packet := range data {
		var tireTemps = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.Tire_temps)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.Tire_temps)), LowerBound: 1}},
		}
		for i, v := range packet.Tire_temps {
			tireTemps.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var tirePressures = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.Tire_pressures)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.Tire_pressures)), LowerBound: 1}},
		}
		for i, v := range packet.Tire_pressures {
			tirePressures.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var location = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.Location)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.Location)), LowerBound: 1}},
		}
		for i, v := range packet.Location {
			location.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var activeSuspensionPowerDraw = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.Active_suspension_power_draw)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.Active_suspension_power_draw)), LowerBound: 1}},
		}
		for i, v := range packet.Active_suspension_power_draw {
			activeSuspensionPowerDraw.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var motorPowerDraw = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.Motor_power_draw)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.Motor_power_draw)), LowerBound: 1}},
		}
		for i, v := range packet.Motor_power_draw {
			motorPowerDraw.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var absThrottleLimiting = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.Abs_throttle_limiting)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.Abs_throttle_limiting)), LowerBound: 1}},
		}
		for i, v := range packet.Abs_throttle_limiting {
			absThrottleLimiting.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var limitedSplitUsage = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.Limited_slip_usage)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.Limited_slip_usage)), LowerBound: 1}},
		}
		for i, v := range packet.Limited_slip_usage {
			limitedSplitUsage.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var id = getLargestId(dbpool)
		_, err := dbpool.Exec(context.Background(), "insert into telemetry (id, lap_id, date_entry, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_loss, abs_throttle_limiting, limited_slip_usage) values ($1, $2, CURRENT_DATE, NOW(), $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)", id, packet.LapId, tireTemps, tirePressures, packet.Velocity, location, packet.Accelerator_input, packet.Brake_input, packet.Steering_angle, packet.Gyro_pitch, packet.Gyro_yaw, packet.Gyro_roll, packet.X_acceleration, packet.Y_acceleration, packet.Z_acceleration, packet.Total_power_draw, activeSuspensionPowerDraw, motorPowerDraw, packet.Battery_voltage, packet.Traction_loss, absThrottleLimiting, limitedSplitUsage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Insertion failed: %v\n", err)
			os.Exit(1)
		}
	}
}

func QueryFromPool(dbpool *pgxpool.Pool) *[]TelemetryPacket {
	rows, err := dbpool.Query(context.Background(), "select id, lap_id, date_entry, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_loss, abs_throttle_limiting, limited_slip_usage from telemetry")
	if err != nil {
		fmt.Fprintf(os.Stderr, "CollectRows error: %v\n", err)
		os.Exit(1)
	}

	var data []TelemetryPacket

	for rows.Next() {
		var packet TelemetryPacket
		err = rows.Scan(&packet.Id, &packet.LapId, &packet.Date_entry, &packet.Time_step, &packet.Tire_temps, &packet.Tire_pressures, &packet.Velocity, &packet.Location, &packet.Accelerator_input, &packet.Brake_input, &packet.Steering_angle, &packet.Gyro_pitch, &packet.Gyro_yaw, &packet.Gyro_roll, &packet.X_acceleration, &packet.Y_acceleration, &packet.Z_acceleration, &packet.Total_power_draw, &packet.Active_suspension_power_draw, &packet.Motor_power_draw, &packet.Battery_voltage, &packet.Traction_loss, &packet.Abs_throttle_limiting, &packet.Limited_slip_usage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		data = append(data, packet)
	}

	return &data
}

func QueryLatestFromPool(dbpool *pgxpool.Pool) *[]TelemetryPacket {
	if LatestID == -404 {
		LatestID = getLatestId(dbpool)
	}
	rows, err := dbpool.Query(context.Background(), "select id, lap_id, date_entry, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_loss, abs_throttle_limiting, limited_slip_usage from telemetry where id > $1", LatestID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CollectRows error: %v\n", err)
		os.Exit(1)
	}
	var data []TelemetryPacket

	for rows.Next() {
		var packet TelemetryPacket
		err = rows.Scan(&packet.Id, &packet.LapId, &packet.Date_entry, &packet.Time_step, &packet.Tire_temps, &packet.Tire_pressures, &packet.Velocity, &packet.Location, &packet.Accelerator_input, &packet.Brake_input, &packet.Steering_angle, &packet.Gyro_pitch, &packet.Gyro_yaw, &packet.Gyro_roll, &packet.X_acceleration, &packet.Y_acceleration, &packet.Z_acceleration, &packet.Total_power_draw, &packet.Active_suspension_power_draw, &packet.Motor_power_draw, &packet.Battery_voltage, &packet.Traction_loss, &packet.Abs_throttle_limiting, &packet.Limited_slip_usage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		data = append(data, packet)
		LatestID = int64(math.Max(float64(packet.Id), float64(LatestID)))
	}

	return &data
}

func QueryLapFromPool(dbpool *pgxpool.Pool, lapId int64) *[]TelemetryPacket {
	rows, err := dbpool.Query(context.Background(), "select id, lap_id, date_entry, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_loss, abs_throttle_limiting, limited_slip_usage from telemetry where lap_id = $1", lapId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CollectRows error: %v\n", err)
		os.Exit(1)
	}
	var data []TelemetryPacket

	for rows.Next() {
		var packet TelemetryPacket
		err = rows.Scan(&packet.Id, &packet.LapId, &packet.Date_entry, &packet.Time_step, &packet.Tire_temps, &packet.Tire_pressures, &packet.Velocity, &packet.Location, &packet.Accelerator_input, &packet.Brake_input, &packet.Steering_angle, &packet.Gyro_pitch, &packet.Gyro_yaw, &packet.Gyro_roll, &packet.X_acceleration, &packet.Y_acceleration, &packet.Z_acceleration, &packet.Total_power_draw, &packet.Active_suspension_power_draw, &packet.Motor_power_draw, &packet.Battery_voltage, &packet.Traction_loss, &packet.Abs_throttle_limiting, &packet.Limited_slip_usage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		data = append(data, packet)
		LatestID = int64(math.Max(float64(packet.Id), float64(LatestID)))
	}

	return &data
}

func QueryBetweenIdsFromPool(dbpool *pgxpool.Pool, startId int64, endId int64) *[]TelemetryPacket {
	rows, err := dbpool.Query(context.Background(), "select id, lap_id, date_entry, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_loss, abs_throttle_limiting, limited_slip_usage from telemetry where id between $1 and $2", startId, endId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CollectRows error: %v\n", err)
		os.Exit(1)
	}
	var data []TelemetryPacket

	for rows.Next() {
		var packet TelemetryPacket
		err = rows.Scan(&packet.Id, &packet.LapId, &packet.Date_entry, &packet.Time_step, &packet.Tire_temps, &packet.Tire_pressures, &packet.Velocity, &packet.Location, &packet.Accelerator_input, &packet.Brake_input, &packet.Steering_angle, &packet.Gyro_pitch, &packet.Gyro_yaw, &packet.Gyro_roll, &packet.X_acceleration, &packet.Y_acceleration, &packet.Z_acceleration, &packet.Total_power_draw, &packet.Active_suspension_power_draw, &packet.Motor_power_draw, &packet.Battery_voltage, &packet.Traction_loss, &packet.Abs_throttle_limiting, &packet.Limited_slip_usage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		data = append(data, packet)
		LatestID = int64(math.Max(float64(packet.Id), float64(LatestID)))
	}

	return &data
}

func QueryMiniPacketsFromPool(dbpool *pgxpool.Pool) *[]MiniTelemetryPacket {
	rows, err := dbpool.Query(context.Background(), "select id, lap_id from telemetry")
	if err != nil {
		fmt.Fprintf(os.Stderr, "CollectRows error: %v\n", err)
		os.Exit(1)
	}

	var data []MiniTelemetryPacket

	for rows.Next() {
		var packet MiniTelemetryPacket
		err = rows.Scan(&packet.Id, &packet.LapId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		data = append(data, packet)
	}

	return &data
}

//Demonstration of code
func main() {
	var dbpool = NewConnection()

	var testData []TelemetryPacket
	for i := 0; i < 100; i++ {
		testData = append(testData, TempTelemetryPacket())
	}
	InsertIntoPool(dbpool, testData)

	var data = QueryFromPool(dbpool)

	println("\nGot Data:")
	for _, packet := range *data {
		fmt.Printf("%d, ", packet.Id)
	}

	var latestData = QueryLatestFromPool(dbpool)

	println("\nGot Latest Data:")
	for _, packet := range *latestData {
		fmt.Printf("%d, ", packet.Id)
	}

	var lapData = QueryLapFromPool(dbpool, 0)

	println("\nGot Lap 0 Data:")
	for _, packet := range *lapData {
		fmt.Printf("%d, ", packet.Id)
	}

	CloseConnection(dbpool)
}
