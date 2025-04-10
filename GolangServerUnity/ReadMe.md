# Golang Telemetry Server for Racecar Data

This Go server is designed to manage and serve telemetry data from a racecar. It uses a SQLite database to store the data and provides several HTTP endpoints to interact with the data, including querying, inserting, and deleting data.

## Overview

The server is built using Go and the `gorilla/mux` router for handling HTTP requests. It stores racecar telemetry data in a SQLite database file (`database.db`). The database schema is defined in `helpers/database.go` (startLine: 7, endLine: 92) and includes tables for:

- **PacketInfo**: Basic information about each data packet (PacketID, SessionID, LapID, Packet Datetime).
- **LapInfo**: Information about each lap (SessionID, LapID, LapTime, Driver, Track, Car).
- **TelemetryInfo**: Real-time telemetry data (Speed, Gas, Brake, Steer, RPM, etc.).
- **TireInfo**: Tire-related telemetry data (Camber, Slip Angle, Slip Ratio, Temperature, Pressure, etc.).

The server provides endpoints to:

- **Query data**: Retrieve telemetry data based on packet or lap ranges.
- **Insert data**: Insert data from CSV files or single JSON rows.
- **Delete data**: Delete entire sessions or individual laps.
- **Retrieve session and lap lists**: Get summaries of available sessions and laps.
- **Get fastest lap**: Find the fastest lap for a given session.
- **Compare laps**: Retrieve telemetry data for multiple laps for comparison.

## Endpoints

The following endpoints are available on the server. All endpoints return JSON responses unless otherwise specified.

### Data Querying

- **`GET /sqliteQuery`**:

  - **Purpose**: General-purpose endpoint to query data from the SQLite database based on PacketID or LapID ranges.
  - **Handler Function**: `handlers.SqliteQuery` in `handlers/requests.go` (startLine: 20, endLine: 127)
  - **Query Parameters**:
    - `table`: The name of the table to query (`PacketInfo`, `LapInfo`, `TelemetryInfo`, `TireInfo`).
    - `start`: The starting PacketID or LapID for the query range.
    - `end`: The ending PacketID or LapID for the query range.
  - **Example**:
    ```bash
    curl "http://localhost:8080/sqliteQuery?table=TelemetryInfo&start=100&end=200"
    ```

- **`GET /sessions`**:

  - **Purpose**: Retrieve a list of all recorded sessions, optionally filtered.
  - **Handler Function**: `handlers.GetSessionList` in `handlers/requests.go` (startLine: 373, endLine: 461)
  - **Query Parameters (Optional)**:
    - `driver`: Filter by driver name.
    - `track`: Filter by track name.
    - `car`: Filter by car name.
    - `date_start`: Filter sessions starting from this date (YYYY-MM-DD format, adjust based on your PacketDatetime format).
    - `date_end`: Filter sessions ending by this date (YYYY-MM-DD format, adjust based on your PacketDatetime format).
  - **Example**:
    ```bash
    curl "http://localhost:8080/sessions?driver=Lewis&track=Monza"
    ```

- **`GET /sessions/{session_id}/laps`**:

  - **Purpose**: Retrieve a list of laps for a specific session.
  - **Handler Function**: `handlers.GetLapList` in `handlers/requests.go` (startLine: 463, endLine: 507)
  - **Path Parameter**:
    - `session_id`: The ID of the session.
  - **Example**:
    ```bash
    curl "http://localhost:8080/sessions/123/laps"
    ```

- **`GET /sessions/{session_id}/fastest_lap`**:

  - **Purpose**: Retrieve the fastest lap for a specific session.
  - **Handler Function**: `handlers.GetFastestLap` in `handlers/requests.go` (startLine: 509, endLine: 579)
  - **Path Parameter**:
    - `session_id`: The ID of the session.
  - **Example**:
    ```bash
    curl "http://localhost:8080/sessions/123/fastest_lap"
    ```

- **`GET /laps/compare`**:
  - **Purpose**: Retrieve telemetry data for multiple laps for comparison.
  - **Handler Function**: `handlers.GetLapComparisonData` in `handlers/requests.go` (startLine: 581, endLine: 703)
  - **Query Parameter**:
    - `lap_ids`: A comma-separated list of LapIDs to compare.
  - **Example**:
    ```bash
    curl "http://localhost:8080/laps/compare?lap_ids=1,2,3"
    ```

### Data Insertion

- **`POST /csvInsert`**:

  - **Purpose**: Insert data from a CSV file into the database. The CSV headers should match the column names in the database tables. Data is inserted into relevant tables based on the headers present in the CSV.
  - **Handler Function**: `handlers.CsvInsert` in `handlers/requests.go` (startLine: 139, endLine: 242)
  - **Request Body**: CSV data in the request body.
  - **Example**:
    ```bash
    curl -X POST -H "Content-Type: text/csv" --data-binary "@data.csv" http://localhost:8080/csvInsert
    ```
    _(Assuming you have a `data.csv` file in the same directory)_

- **`POST /addRow`**:
  - **Purpose**: Insert a single row of data into the database across relevant tables using a JSON request body. The JSON keys should correspond to column names.
  - **Handler Function**: `handlers.AddRow` in `handlers/requests.go` (startLine: 249, endLine: 371)
  - **Request Body**: JSON data in the request body.
  - **Example**:
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"PacketID": "1234", "SessionID": "5678", "LapID": "1", "PacketDatetime": "2024-08-03 12:00:00", "SpeedMPH": "100.5"}' http://localhost:8080/addRow
    ```

### Data Deletion

- **`DELETE /sessions/{session_id}`**:

  - **Purpose**: Delete a specific session and all associated data (laps, telemetry, packets).
  - **Handler Function**: `handlers.DeleteSession` in `handlers/requests.go` (startLine: 705, endLine: 770)
  - **Path Parameter**:
    - `session_id`: The ID of the session to delete.
  - **Example**:
    ```bash
    curl -X DELETE http://localhost:8080/sessions/123
    ```

- **`DELETE /laps/{lap_id}`**:
  - **Purpose**: Delete a specific lap and its associated telemetry and packet data.
  - **Handler Function**: `handlers.DeleteLap` in `handlers/requests.go` (startLine: 772, endLine: 837)
  - **Path Parameter**:
    - `lap_id`: The ID of the lap to delete.
  - **Example**:
    ```bash
    curl -X DELETE http://localhost:8080/laps/123
    ```

## Database Schema

The database is structured into four main tables, each designed to store specific aspects of racecar telemetry data. Below is a breakdown of each table and its columns:

### PacketInfo Table

This table stores basic information about each data packet received.

- **PacketID**: `INTEGER` - Unique identifier for each packet.
- **SessionID**: `INTEGER` - Identifier for the session to which the packet belongs.
- **LapID**: `INTEGER` - Identifier for the lap to which the packet belongs.
- **PacketDatetime**: `TEXT` - Timestamp indicating when the packet data was recorded.

### LapInfo Table

This table contains information about each lap within a session.

- **SessionID**: `INTEGER` - Identifier for the session.
- **LapID**: `INTEGER` - Identifier for the lap within the session.
- **LapTime**: `INTEGER` - Duration of the lap, likely in milliseconds or a similar unit.
- **DriverName**: `TEXT` - Name of the driver for the lap.
- **TrackName**: `TEXT` - Name of the track where the lap was driven.
- **TrackConfiguration**: `TEXT` - Specific configuration of the track, if applicable.
- **CarName**: `TEXT` - Name of the car used for the lap.

### TelemetryInfo Table

This table stores real-time telemetry data from the racecar.

- **PacketID**: `INTEGER` - Foreign key referencing `PacketInfo.PacketID`.
- **SessionID**: `INTEGER` - Foreign key referencing `LapInfo.SessionID`.
- **LapID**: `INTEGER` - Foreign key referencing `LapInfo.LapID`.
- **SpeedMPH**: `REAL` - Speed of the car in miles per hour.
- **Gas**: `REAL` - Throttle input, typically a value between 0 and 1.
- **Brake**: `REAL` - Brake input, typically a value between 0 and 1.
- **Steer**: `REAL` - Steering input, typically a value between -1 and 1.
- **Clutch**: `REAL` - Clutch input.
- **Gear**: `INTEGER` - Currently engaged gear.
- **RPM**: `REAL` - Engine revolutions per minute.
- **TurboBoost**: `REAL` - Turbo boost pressure.
- **LocalAngularVelocityX**: `REAL` - Angular velocity around the X-axis in local space.
- **LocalAngularVelocityY**: `REAL` - Angular velocity around the Y-axis in local space.
- **LocalAngularVelocityZ**: `REAL` - Angular velocity around the Z-axis in local space.
- **VelocityX**: `REAL` - Velocity in the X direction.
- **VelocityY**: `REAL` - Velocity in the Y direction.
- **VelocityZ**: `REAL` - Velocity in the Z direction.
- **WorldPositionX**: `REAL` - Car's position in world space X coordinate.
- **WorldPositionY**: `REAL` - Car's position in world space Y coordinate.
- **WorldPositionZ**: `REAL` - Car's position in world space Z coordinate.
- **Aero_DragCoeffcient**: `REAL` - Aerodynamic drag coefficient.
- **Aero_LiftCoefficientFront**: `REAL` - Front aerodynamic lift coefficient.
- **Aero_LiftCoefficientRear**: `REAL` - Rear aerodynamic lift coefficient.

### TireInfo Table

This table stores tire-related telemetry data.

- **PacketID**: `INTEGER` - Foreign key referencing `PacketInfo.PacketID`.
- **SessionID**: `INTEGER` - Foreign key referencing `LapInfo.SessionID`.
- **LapID**: `INTEGER` - Foreign key referencing `LapInfo.LapID`.
- **FL_CamberRad**: `REAL` - Front Left tire camber in radians.
- **FR_CamberRad**: `REAL` - Front Right tire camber in radians.
- **RL_CamberRad**: `REAL` - Rear Left tire camber in radians.
- **RR_CamberRad**: `REAL` - Rear Right tire camber in radians.
- **FL_SlipAngle**: `REAL` - Front Left tire slip angle.
- **FR_SlipAngle**: `REAL` - Front Right tire slip angle.
- **RL_SlipAngle**: `REAL` - Rear Left tire slip angle.
- **RR_SlipAngle**: `REAL` - Rear Right tire slip angle.
- **FL_SlipRatio**: `REAL` - Front Left tire slip ratio.
- **FR_SlipRatio**: `REAL` - Front Right tire slip ratio.
- **RL_SlipRatio**: `REAL` - Rear Left tire slip ratio.
- **RR_SlipRatio**: `REAL` - Rear Right tire slip ratio.
- **FL_SelfAligningTorque**: `REAL` - Front Left tire self-aligning torque.
- **FR_SelfAligningTorque**: `REAL` - Front Right tire self-aligning torque.
- **RL_SelfAligningTorque**: `REAL` - Rear Left tire self-aligning torque.
- **RR_SelfAligningTorque**: `REAL` - Rear Right tire self-aligning torque.
- **FL_Load**: `REAL` - Front Left tire load.
- **FR_Load**: `REAL` - Front Right tire load.
- **RL_Load**: `REAL` - Rear Left tire load.
- **RR_Load**: `REAL` - Rear Right tire load.
- **FL_TyreSlip**: `REAL` - Front Left tire slip.
- **FR_TyreSlip**: `REAL` - Front Right tire slip.
- **RL_TyreSlip**: `REAL` - Rear Left tire slip.
- **RR_TyreSlip**: `REAL` - Rear Right tire slip.
- **FL_ThermalState**: `REAL` - Front Left tire thermal state.
- **FR_ThermalState**: `REAL` - Front Right tire thermal state.
- **RL_ThermalState**: `REAL` - Rear Left tire thermal state.
- **RR_ThermalState**: `REAL` - Rear Right tire thermal state.
- **FL_DynamicPressure**: `REAL` - Front Left tire dynamic pressure.
- **FR_DynamicPressure**: `REAL` - Front Right tire dynamic pressure.
- **RL_DynamicPressure**: `REAL` - Rear Left tire dynamic pressure.
- **RR_DynamicPressure**: `REAL` - Rear Right tire dynamic pressure.
- **FL_TyreDirtyLevel**: `REAL` - Front Left tire dirt level.
- **FR_TyreDirtyLevel**: `REAL` - Front Right tire dirt level.
- **RL_TyreDirtyLevel**: `REAL` - Rear Left tire dirt level.
- **RR_TyreDirtyLevel**: `REAL` - Rear Right tire dirt level.

## Database Setup

The server uses a SQLite database file named `database.db` in the same directory as the server executable. The database schema is automatically created when the server starts if it doesn't already exist. The schema is defined in the `SetupDatabaseSchema` function in `helpers/database.go` (startLine: 7, endLine: 92).

## Running the Server

1.  **Install Go**: Ensure you have Go installed on your system.
2.  **Navigate to the project directory**: Open a terminal and navigate to the directory containing the server code (`cmd/main.go`, `handlers/requests.go`, `helpers/`).
3.  **Run the server**: Execute the following command in the terminal:
    ```bash
    go run cmd/main.go
    ```
    The server will start and listen on `http://localhost:8080`.

## Example Usage

**1. Insert CSV data:**

Create a CSV file named `data.csv` with headers matching your database schema (e.g., `PacketID,SessionID,LapID,PacketDatetime,SpeedMPH`).

```csv
PacketID,SessionID,LapID,PacketDatetime,SpeedMPH
1,100,1,2024-08-03 10:00:00,50.2
2,100,1,2024-08-03 10:00:01,51.5
3,100,1,2024-08-03 10:00:02,52.8
```

Run the CSV insert command:

```bash
curl -X POST -H "Content-Type: text/csv" --data-binary "@data.csv" http://localhost:8080/csvInsert
```

**2. Add a single row using JSON:**

```bash
curl -X POST -H "Content-Type: application/json" -d '{"PacketID": "4", "SessionID": "100", "LapID": "1", "PacketDatetime": "2024-08-03 10:00:03", "SpeedMPH": "53.1"}' http://localhost:8080/addRow
```

**3. Query telemetry data:**

```bash
curl "http://localhost:8080/sqliteQuery?table=TelemetryInfo&start=1&end=4"
```

**4. Get laps for session 100:**

```bash
curl "http://localhost:8080/sessions/100/laps"
```

**5. Delete lap 1:**

```bash
curl -X DELETE http://localhost:8080/laps/1
```
