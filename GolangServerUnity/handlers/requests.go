package handlers

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	"github.com/zleodai/FSAESoftware/GolangServerUnity/helpers"
)

// SqliteQuery handles HTTP requests to query data from an SQLite database.
// It expects 'start', 'end', and 'table' query parameters to define the data range and table to query.
// It returns the query results as a JSON response.
func SqliteQuery(w http.ResponseWriter, r *http.Request) {

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	table := r.URL.Query().Get("table")

	if start == "" || end == "" || table == "" {
		http.Error(w, "Missing start, end, or table", 512)
		return
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = helpers.SetupDatabaseSchema(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var query string
	var rows *sql.Rows

	switch table {
	case "PacketInfo":
		query = "SELECT * FROM PacketInfo WHERE PacketID >= ? AND PacketID <= ?"
		rows, err = db.Query(query, start, end)
	case "LapInfo":
		query = "SELECT * FROM LapInfo WHERE LapID >= ? AND LapID <= ?" // Assuming LapID is the relevant range for LapInfo
		rows, err = db.Query(query, start, end)
	case "TelemetryInfo":
		query = "SELECT * FROM TelemetryInfo WHERE PacketID >= ? AND PacketID <= ?"
		rows, err = db.Query(query, start, end)
	case "TireInfo":
		query = "SELECT * FROM TireInfo WHERE PacketID >= ? AND PacketID <= ?"
		rows, err = db.Query(query, start, end)
	default:
		http.Error(w, "Invalid table name", 512)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		err = rows.Scan(scanArgs...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			if b, ok := val.([]byte); ok {
				var intVal int
				if table == "LapInfo" && (col == "LapTime") { //special case for LapTime in LapInfo table
					intVal, err = strconv.Atoi(string(b))
					if err != nil {
						http.Error(w, "Error converting LapTime to integer", http.StatusInternalServerError)
						return
					}
					row[col] = intVal
				} else {
					row[col] = string(b)
				}

			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(results) == 0 {
		http.Error(w, "No records found in the specified range", 513)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

var tableSchemas = map[string][]string{
	"PacketInfo":    {"PacketID", "SessionID", "LapID", "PacketDatetime"},
	"LapInfo":       {"SessionID", "LapID", "LapTime", "DriverName", "TrackName", "TrackConfiguration", "CarName"},
	"TelemetryInfo": {"PacketID", "SessionID", "LapID", "SpeedMPH", "Gas", "Brake", "Steer", "Clutch", "Gear", "RPM", "TurboBoost", "LocalAngularVelocityX", "LocalAngularVelocityY", "LocalAngularVelocityZ", "VelocityX", "VelocityY", "VelocityZ", "WorldPositionX", "WorldPositionY", "WorldPositionZ", "Aero_DragCoeffcient", "Aero_LiftCoefficientFront", "Aero_LiftCoefficientRear"},
	"TireInfo":      {"PacketID", "SessionID", "LapID", "FL_CamberRad", "FR_CamberRad", "RL_CamberRad", "RR_CamberRad", "FL_SlipAngle", "FR_SlipAngle", "RL_SlipAngle", "RR_SlipAngle", "FL_SlipRatio", "FR_SlipRatio", "RL_SlipRatio", "RR_SlipRatio", "FL_SelfAligningTorque", "FR_SelfAligningTorque", "RL_SelfAligningTorque", "RR_SelfAligningTorque", "FL_Load", "FR_Load", "RL_Load", "RR_Load", "FL_TyreSlip", "FR_TyreSlip", "RL_TyreSlip", "RR_TyreSlip", "FL_ThermalState", "FR_ThermalState", "RL_ThermalState", "RR_ThermalState", "FL_DynamicPressure", "FR_DynamicPressure", "RL_DynamicPressure", "RR_DynamicPressure", "FL_TyreDirtyLevel", "FR_TyreDirtyLevel", "RL_TyreDirtyLevel", "RR_TyreDirtyLevel"},
}

// CsvInsert handles HTTP POST requests to insert CSV data into an SQLite database.
// It reads CSV data from the request body and inserts it into the appropriate tables based on the headers.
// It expects the CSV data to have a header row and at least one data row.
func CsvInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = helpers.SetupDatabaseSchema(db)
	if err != nil {
		http.Error(w, "Failed to setup database schema", http.StatusInternalServerError)
		return
	}

	csvReader := csv.NewReader(r.Body) // Read from request body
	defer r.Body.Close()               // Close request body after reading
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		http.Error(w, "Failed to read CSV data from request body", http.StatusBadRequest)
		return
	}

	if len(records) <= 1 {
		http.Error(w, "CSV data must have header and at least one data row", http.StatusBadRequest)
		return
	}

	headerRow := records[0]

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for i, record := range records[1:] { // Skip header row
		if len(record) != len(headerRow) {
			http.Error(w, fmt.Sprintf("Row %d column count does not match header", i+2), http.StatusBadRequest)
			return
		}

		for tableName, schema := range tableSchemas {
			tableHeaders := make([]string, 0)
			tableValues := make([]string, 0)

			for colIdx, header := range headerRow {
				for _, schemaCol := range schema {
					if strings.EqualFold(header, schemaCol) { // Case-insensitive comparison
						value := record[colIdx]
						if value != "" { // Only include columns with non-empty values
							tableHeaders = append(tableHeaders, header)
							tableValues = append(tableValues, value)
						}
						break // Move to the next header after finding a schema match
					}
				}
			}

			if len(tableHeaders) > 0 { // Only insert if there are columns for this table in the current row
				valuePlaceholders := make([]string, len(tableHeaders))
				for i := range tableHeaders {
					valuePlaceholders[i] = "?"
				}
				insertQuery := "INSERT INTO " + tableName + " (" + strings.Join(tableHeaders, ", ") + ") VALUES (" + strings.Join(valuePlaceholders, ", ") + ")"

				stmt, err := tx.Prepare(insertQuery)
				if err != nil {
					http.Error(w, fmt.Sprintf("Failed to prepare insert statement for table '%s': %s", tableName, err.Error()), http.StatusInternalServerError)
					return
				}
				defer stmt.Close()

				var args []interface{}
				for _, val := range tableValues {
					args = append(args, val)
				}

				_, err = stmt.Exec(args...)
				if err != nil {
					http.Error(w, fmt.Sprintf("Failed to execute insert statement for table '%s' in row %d: %s", tableName, i+2, err.Error()), http.StatusInternalServerError)
					return
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("CSV data successfully inserted into relevant tables"))
}

// AddRow handles HTTP POST requests to insert a single row of data into multiple SQLite database tables.
// It expects a JSON request body where the keys correspond to column names across different tables.
// The function parses the JSON, identifies the relevant table for each column based on predefined schemas,
// and inserts the data into the corresponding tables within a single transaction.
// It responds with a success message or an error if insertion fails.
func AddRow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var rowData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&rowData)
	if err != nil {
		http.Error(w, "Failed to decode JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = helpers.SetupDatabaseSchema(db)
	if err != nil {
		http.Error(w, "Failed to setup database schema", http.StatusInternalServerError)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	for tableName, schema := range tableSchemas {
		tableHeaders := make([]string, 0)
		tableValues := make([]interface{}, 0) // Use interface{} to accommodate various data types

		for header, value := range rowData {
			for _, schemaCol := range schema {
				if strings.EqualFold(header, schemaCol) {
					tableHeaders = append(tableHeaders, header)
					tableValues = append(tableValues, value)
					break
				}
			}
		}

		if len(tableHeaders) > 0 {
			valuePlaceholders := make([]string, len(tableHeaders))
			for i := range tableHeaders {
				valuePlaceholders[i] = "?"
			}
			insertQuery := "INSERT INTO " + tableName + " (" + strings.Join(tableHeaders, ", ") + ") VALUES (" + strings.Join(valuePlaceholders, ", ") + ")"

			stmt, err := tx.Prepare(insertQuery)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to prepare insert statement for table '%s': %s", tableName, err.Error()), http.StatusInternalServerError)
				return
			}
			defer stmt.Close()

			_, err = stmt.Exec(tableValues...)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to execute insert statement for table '%s': %s", tableName, err.Error()), http.StatusInternalServerError)
				return
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Row data successfully inserted into relevant tables"))
}

// GetSessionList handles GET requests to retrieve a list of sessions, optionally filtered by driver, track, car, and date range.
func GetSessionList(w http.ResponseWriter, r *http.Request) {
	driver := r.URL.Query().Get("driver")
	track := r.URL.Query().Get("track")
	car := r.URL.Query().Get("car")
	dateStart := r.URL.Query().Get("date_start")
	dateEnd := r.URL.Query().Get("date_end")

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = helpers.SetupDatabaseSchema(db)
	if err != nil {
		http.Error(w, "Failed to setup database schema", http.StatusInternalServerError)
		return
	}

	query := `SELECT DISTINCT SessionID, DriverName, TrackName, TrackConfiguration, CarName FROM LapInfo WHERE 1=1` // 1=1 to easily append filters
	var filters []string
	var args []interface{}

	if driver != "" {
		filters = append(filters, "DriverName = ?")
		args = append(args, driver)
	}
	if track != "" {
		filters = append(filters, "TrackName = ?")
		args = append(args, track)
	}
	if car != "" {
		filters = append(filters, "CarName = ?")
		args = append(args, car)
	}
	if dateStart != "" && dateEnd != "" { // Assuming you have a date/time column, adjust accordingly
		filters = append(filters, "PacketDatetime BETWEEN ? AND ?") //Example, adjust column name
		args = append(args, dateStart, dateEnd)                     //Example, adjust column name
	}

	if len(filters) > 0 {
		query += " AND " + strings.Join(filters, " AND ")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query sessions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sessions []map[string]interface{}
	columns, _ := rows.Columns() // Assuming column names are consistent
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err = rows.Scan(scanArgs...)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning session row: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		session := make(map[string]interface{})
		for i, col := range columns {
			session[col] = helpers.ConvertSqlValue(values[i]) // Assuming you have a helper to convert sql values
		}
		sessions = append(sessions, session)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// GetLapList handles GET requests to retrieve a list of laps for a specific session.
func GetLapList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)             // gorilla mux
	sessionID := vars["session_id"] // gorilla mux

	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = helpers.SetupDatabaseSchema(db)
	if err != nil {
		http.Error(w, "Failed to setup database schema", http.StatusInternalServerError)
		return
	}

	query := "SELECT LapID, LapTime FROM LapInfo WHERE SessionID = ?" // Select relevant lap info
	rows, err := db.Query(query, sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query laps for session %s: %s", sessionID, err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var laps []map[string]interface{}
	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err = rows.Scan(scanArgs...)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning lap row: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		lap := make(map[string]interface{})
		for i, col := range columns {
			lap[col] = helpers.ConvertSqlValue(values[i])
		}
		laps = append(laps, lap)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(laps)
}

// GetFastestLap handles GET requests to find and return the fastest lap for a given session.
func GetFastestLap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)             // gorilla mux
	sessionID := vars["session_id"] // gorilla mux

	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = helpers.SetupDatabaseSchema(db)
	if err != nil {
		http.Error(w, "Failed to setup database schema", http.StatusInternalServerError)
		return
	}

	query := `SELECT LapID, LapTime FROM LapInfo WHERE SessionID = ? ORDER BY LapTime ASC LIMIT 1`
	row := db.QueryRow(query, sessionID)

	var lap map[string]interface{} = make(map[string]interface{})
	var lapID int
	var lapTime int // Assuming LapTime is int, adjust if needed

	err = row.Scan(&lapID, &lapTime)
	if err == sql.ErrNoRows {
		http.Error(w, fmt.Sprintf("No laps found for session %s", sessionID), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("Failed to query fastest lap: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	lap["LapID"] = lapID
	lap["LapTime"] = lapTime

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lap)
}

// GetLapComparisonData handles GET requests to retrieve telemetry data for multiple laps for comparison.
func GetLapComparisonData(w http.ResponseWriter, r *http.Request) {
	lapIDsStr := r.URL.Query().Get("lap_ids")
	if lapIDsStr == "" {
		http.Error(w, "Lap IDs are required", http.StatusBadRequest)
		return
	}

	lapIDStrs := strings.Split(lapIDsStr, ",")
	var lapIDs []int
	for _, idStr := range lapIDStrs {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid lap ID: %s", idStr), http.StatusBadRequest)
			return
		}
		lapIDs = append(lapIDs, id)
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = helpers.SetupDatabaseSchema(db)
	if err != nil {
		http.Error(w, "Failed to setup database schema", http.StatusInternalServerError)
		return
	}

	comparisonData := make(map[string][]map[string]interface{}) // Map of lapID -> telemetry data

	for _, lapID := range lapIDs {
		var lapTelemetryData []map[string]interface{}

		// Example: Query TelemetryInfo for each lap. Adjust tables as needed.
		query := "SELECT * FROM TelemetryInfo WHERE LapID = ?"
		rows, err := db.Query(query, lapID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to query telemetry data for lap %d: %s", lapID, err.Error()), http.StatusInternalServerError)
			return // Or continue and just skip this lap if you want to be more lenient
		}
		defer rows.Close()

		columns, _ := rows.Columns()
		for rows.Next() {
			values := make([]interface{}, len(columns))
			scanArgs := make([]interface{}, len(columns))
			for i := range values {
				scanArgs[i] = &values[i]
			}
			err = rows.Scan(scanArgs...)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning telemetry row for lap %d: %s", lapID, err.Error()), http.StatusInternalServerError)
				return // Or continue and just skip this row
			}

			rowMap := make(map[string]interface{})
			for i, col := range columns {
				rowMap[col] = helpers.ConvertSqlValue(values[i])
			}
			lapTelemetryData = append(lapTelemetryData, rowMap)
		}
		comparisonData[strconv.Itoa(lapID)] = lapTelemetryData // Store telemetry data in the map, key is lapID string
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comparisonData)
}

// DeleteSession handles DELETE requests to remove a session and all associated data.
func DeleteSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)             // gorilla mux
	sessionID := vars["session_id"] // gorilla mux

	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete from PacketInfo
	_, err = tx.Exec("DELETE FROM PacketInfo WHERE SessionID = ?", sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete PacketInfo for session %s: %s", sessionID, err.Error()), http.StatusInternalServerError)
		return
	}

	// Delete from TelemetryInfo
	_, err = tx.Exec("DELETE FROM TelemetryInfo WHERE SessionID = ?", sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete TelemetryInfo for session %s: %s", sessionID, err.Error()), http.StatusInternalServerError)
		return
	}

	// Delete from TireInfo
	_, err = tx.Exec("DELETE FROM TireInfo WHERE SessionID = ?", sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete TireInfo for session %s: %s", sessionID, err.Error()), http.StatusInternalServerError)
		return
	}

	// Delete from LapInfo (delete session info last in case of FK constraints)
	_, err = tx.Exec("DELETE FROM LapInfo WHERE SessionID = ?", sessionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete LapInfo for session %s: %s", sessionID, err.Error()), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to commit transaction for session deletion: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Session %s and associated data deleted successfully", sessionID)))
}

// DeleteLap handles DELETE requests to remove a specific lap and its associated telemetry data.
func DeleteLap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)        // gorilla mux
	lapIDStr := vars["lap_id"] // gorilla mux
	lapID, err := strconv.Atoi(lapIDStr)
	if err != nil {
		http.Error(w, "Invalid Lap ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		http.Error(w, "Failed to open database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete from PacketInfo
	_, err = tx.Exec("DELETE FROM PacketInfo WHERE LapID = ?", lapID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete PacketInfo for lap %d: %s", lapID, err.Error()), http.StatusInternalServerError)
		return
	}

	// Delete from TelemetryInfo
	_, err = tx.Exec("DELETE FROM TelemetryInfo WHERE LapID = ?", lapID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete TelemetryInfo for lap %d: %s", lapID, err.Error()), http.StatusInternalServerError)
		return
	}

	// Delete from TireInfo
	_, err = tx.Exec("DELETE FROM TireInfo WHERE LapID = ?", lapID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete TireInfo for lap %d: %s", lapID, err.Error()), http.StatusInternalServerError)
		return
	}

	// Delete from LapInfo (delete lap info last)
	_, err = tx.Exec("DELETE FROM LapInfo WHERE LapID = ?", lapID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete LapInfo for lap %d: %s", lapID, err.Error()), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to commit transaction for lap deletion: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Lap %d and associated data deleted successfully", lapID)))
}
