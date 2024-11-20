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

type TelemetryPacket struct {
	id                           int64
	date_entry                   time.Time
	time_step                    time.Time
	tire_temps                   [4096]float64
	tire_pressures               [4]float64
	velocity                     float64
	location                     [2]float64
	accelerator_input            float64
	brake_input                  float64
	steering_angle               float64
	gyro_pitch                   float64
	gyro_yaw                     float64
	gyro_roll                    float64
	x_acceleration               float64
	y_acceleration               float64
	z_acceleration               float64
	total_power_draw             float64
	active_suspension_power_draw [4]float64
	motor_power_draw             [2]float64
	battery_voltage              float64
	traction_loss                float64
	abs_throttle_limiting        [2]float64
	limited_slip_usage           [2]float64
}

func TempTelemtryPacket() TelemetryPacket {
	var packet TelemetryPacket = TelemetryPacket{
		tire_temps:                   [4096]float64{},
		tire_pressures:               [4]float64{rand.Float64()*5 + 25, rand.Float64()*5 + 25, rand.Float64()*5 + 25, rand.Float64()*5 + 25},
		velocity:                     rand.Float64() * 100,
		location:                     [2]float64{rand.Float64(), rand.Float64()},
		accelerator_input:            rand.Float64(),
		brake_input:                  rand.Float64(),
		steering_angle:               rand.Float64(),
		gyro_pitch:                   rand.Float64(),
		gyro_yaw:                     rand.Float64(),
		gyro_roll:                    rand.Float64(),
		x_acceleration:               rand.Float64(),
		y_acceleration:               rand.Float64(),
		z_acceleration:               rand.Float64(),
		total_power_draw:             rand.Float64(),
		active_suspension_power_draw: [4]float64{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()},
		motor_power_draw:             [2]float64{rand.Float64()},
		battery_voltage:              rand.Float64(),
		traction_loss:                rand.Float64(),
		abs_throttle_limiting:        [2]float64{rand.Float64(), rand.Float64()},
		limited_slip_usage:           [2]float64{rand.Float64(), rand.Float64()},
	}

	for i := 0; i < 4096; i++ {
		packet.tire_temps[i] = 1
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
			newLargestID = int64(math.Max(float64(newLargestID), float64(packet.id)))
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
		latestID = int64(math.Max(float64(latestID), float64(packet.id)))
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
			Elements: make([]pgtype.Float8, len(packet.tire_temps)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.tire_temps)), LowerBound: 1}},
		}
		for i, v := range packet.tire_temps {
			tireTemps.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var tirePressures = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.tire_pressures)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.tire_pressures)), LowerBound: 1}},
		}
		for i, v := range packet.tire_pressures {
			tirePressures.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var location = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.location)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.location)), LowerBound: 1}},
		}
		for i, v := range packet.location {
			location.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var activeSuspensionPowerDraw = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.active_suspension_power_draw)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.active_suspension_power_draw)), LowerBound: 1}},
		}
		for i, v := range packet.active_suspension_power_draw {
			activeSuspensionPowerDraw.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var motorPowerDraw = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.motor_power_draw)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.motor_power_draw)), LowerBound: 1}},
		}
		for i, v := range packet.motor_power_draw {
			motorPowerDraw.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var absThrottleLimiting = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.abs_throttle_limiting)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.abs_throttle_limiting)), LowerBound: 1}},
		}
		for i, v := range packet.abs_throttle_limiting {
			absThrottleLimiting.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var limitedSplitUsage = pgtype.Array[pgtype.Float8]{
			Elements: make([]pgtype.Float8, len(packet.limited_slip_usage)),
			Dims:     []pgtype.ArrayDimension{{Length: int32(len(packet.limited_slip_usage)), LowerBound: 1}},
		}
		for i, v := range packet.limited_slip_usage {
			limitedSplitUsage.Elements[i] = pgtype.Float8{Float64: v, Valid: true}
		}

		var id = getLargestId(dbpool)
		_, err := dbpool.Exec(context.Background(), "insert into telemetry (id, date_entry, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_loss, abs_throttle_limiting, limited_slip_usage) values ($1, CURRENT_DATE, NOW(), $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)", id, tireTemps, tirePressures, packet.velocity, location, packet.accelerator_input, packet.brake_input, packet.steering_angle, packet.gyro_pitch, packet.gyro_yaw, packet.gyro_roll, packet.x_acceleration, packet.y_acceleration, packet.z_acceleration, packet.total_power_draw, activeSuspensionPowerDraw, motorPowerDraw, packet.battery_voltage, packet.traction_loss, absThrottleLimiting, limitedSplitUsage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Insertion failed: %v\n", err)
			os.Exit(1)
		}
	}
}

func QueryFromPool(dbpool *pgxpool.Pool) *[]TelemetryPacket {
	rows, err := dbpool.Query(context.Background(), "select id, date_entry, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_loss, abs_throttle_limiting, limited_slip_usage from telemetry")
	if err != nil {
		fmt.Fprintf(os.Stderr, "CollectRows error: %v\n", err)
		os.Exit(1)
	}

	var data []TelemetryPacket

	for rows.Next() {
		var packet TelemetryPacket
		err = rows.Scan(&packet.id, &packet.date_entry, &packet.time_step, &packet.tire_temps, &packet.tire_pressures, &packet.velocity, &packet.location, &packet.accelerator_input, &packet.brake_input, &packet.steering_angle, &packet.gyro_pitch, &packet.gyro_yaw, &packet.gyro_roll, &packet.x_acceleration, &packet.y_acceleration, &packet.z_acceleration, &packet.total_power_draw, &packet.active_suspension_power_draw, &packet.motor_power_draw, &packet.battery_voltage, &packet.traction_loss, &packet.abs_throttle_limiting, &packet.limited_slip_usage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		data = append(data, packet)
	}

	return &data
}

func QueryLatestFromPool(dbpool *pgxpool.Pool) *[]TelemetryPacket {
	rows, err := dbpool.Query(context.Background(), "select id, date_entry, time_step, tire_temps, tire_pressures, velocity, location, accelerator_input, brake_input, steering_angle, gyro_pitch, gyro_yaw, gyro_roll, x_acceleration, y_acceleration, z_acceleration, total_power_draw, active_suspension_power_draw, motor_power_draw, battery_voltage, traction_loss, abs_throttle_limiting, limited_slip_usage from telemetry where id > $1", LatestID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CollectRows error: %v\n", err)
		os.Exit(1)
	}
	if LatestID == -404 {
		LatestID = getLatestId(dbpool)
	}
	var data []TelemetryPacket

	for rows.Next() {
		var packet TelemetryPacket
		err = rows.Scan(&packet.id, &packet.date_entry, &packet.time_step, &packet.tire_temps, &packet.tire_pressures, &packet.velocity, &packet.location, &packet.accelerator_input, &packet.brake_input, &packet.steering_angle, &packet.gyro_pitch, &packet.gyro_yaw, &packet.gyro_roll, &packet.x_acceleration, &packet.y_acceleration, &packet.z_acceleration, &packet.total_power_draw, &packet.active_suspension_power_draw, &packet.motor_power_draw, &packet.battery_voltage, &packet.traction_loss, &packet.abs_throttle_limiting, &packet.limited_slip_usage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		data = append(data, packet)
		LatestID = int64(math.Max(float64(packet.id), float64(LatestID)))
	}

	return &data
}

func main() {
	var dbpool = NewConnection()

	var testData []TelemetryPacket
	for i := 0; i < 100; i++ {
		testData = append(testData, TempTelemtryPacket())
	}
	InsertIntoPool(dbpool, testData)

	var data = QueryFromPool(dbpool)

	println("\nGot Data:")
	for _, packet := range *data {
		fmt.Printf("%d, ", packet.id)
	}

	var latestData = QueryLatestFromPool(dbpool)

	println("\nGot Latest Data:")
	for _, packet := range *latestData {
		fmt.Printf("%d, ", packet.id)
	}

	CloseConnection(dbpool)
}
