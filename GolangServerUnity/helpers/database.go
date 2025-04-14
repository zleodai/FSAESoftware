package helpers

import (
	"database/sql"
)

func SetupDatabaseSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS PacketInfo (
			PacketID INTEGER PRIMARY KEY,
			SessionID INTEGER,
			LapID INTEGER,
			PacketDatetime TEXT
		);

		CREATE TABLE IF NOT EXISTS LapInfo (
			SessionID INTEGER,
			LapID INTEGER,
			LapTime INTEGER,
			DriverName TEXT,
			TrackName TEXT,
			TrackConfiguration TEXT,
			CarName TEXT
		);

		CREATE TABLE IF NOT EXISTS TelemetryInfo (
			PacketID INTEGER PRIMARY KEY,
			SessionID INTEGER,
			LapID INTEGER,
			SpeedMPH REAL,
			Gas REAL,
			Brake REAL,
			Steer REAL,
			Clutch REAL,
			Gear INTEGER,
			RPM REAL,
			TurboBoost REAL,
			LocalAngularVelocityX REAL,
			LocalAngularVelocityY REAL,
			LocalAngularVelocityZ REAL,
			VelocityX REAL,
			VelocityY REAL,
			VelocityZ REAL,
			WorldPositionX REAL,
			WorldPositionY REAL,
			WorldPositionZ REAL,
			Aero_DragCoeffcient REAL,
			Aero_LiftCoefficientFront REAL,
			Aero_LiftCoefficientRear REAL
		);

		CREATE TABLE IF NOT EXISTS TireInfo (
			PacketID INTEGER PRIMARY KEY,
			SessionID INTEGER,
			LapID INTEGER,
			FL_CamberRad REAL,
			FR_CamberRad REAL,
			RL_CamberRad REAL,
			RR_CamberRad REAL,
			FL_SlipAngle REAL,
			FR_SlipAngle REAL,
			RL_SlipAngle REAL,
			RR_SlipAngle REAL,
			FL_SlipRatio REAL,
			FR_SlipRatio REAL,
			RL_SlipRatio REAL,
			RR_SlipRatio REAL,
			FL_SelfAligningTorque REAL,
			FR_SelfAligningTorque REAL,
			RL_SelfAligningTorque REAL,
			RR_SelfAligningTorque REAL,
			FL_Load REAL,
			FR_Load REAL,
			RL_Load REAL,
			RR_Load REAL,
			FL_TyreSlip REAL,
			FR_TyreSlip REAL,
			RL_TyreSlip REAL,
			RR_TyreSlip REAL,
			FL_ThermalState REAL,
			FR_ThermalState REAL,
			RL_ThermalState REAL,
			RR_ThermalState REAL,
			FL_DynamicPressure REAL,
			FR_DynamicPressure REAL,
			RL_DynamicPressure REAL,
			RR_DynamicPressure REAL,
			FL_TyreDirtyLevel REAL,
			FR_TyreDirtyLevel REAL,
			RL_TyreDirtyLevel REAL,
			RR_TyreDirtyLevel REAL
		);
	`)
	return err
}

// ConvertSqlValue converts sql.Rows value types to Go types for JSON encoding.
// Specifically handles []byte to string conversion.
func ConvertSqlValue(val interface{}) interface{} {
	switch v := val.(type) {
	case []byte:
		return string(v)
	default:
		return v
	}
}
