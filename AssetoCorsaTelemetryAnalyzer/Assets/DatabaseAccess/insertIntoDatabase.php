<?php
$serverName = "localhost"; 
$uid = "TelemetryDBUser";   
$pwd = "123";  
$databaseName = "TelemetryDB"; 
$privateKey = "945";

//Insert types
//1 = insert into PacketInfo, TelemetryInfo, and TireInfo
//2 = insert into LapInfo

$publicKey = (int)$_GET["publicKey"];
$insertType = (int)$_GET["insertType"];
$sessionID = (int)$_GET["SessionID"];
$lapID = (int)$_GET["LapID"];

$calculatedKey = (int)hash('sha256', $insertType * $privateKey * ($sessionID + $lapID));

if ($publicKey == $calculatedKey) {
    $conn = new mysqli('localhost', 'TelemetryDBUser', '123', 'TelemetryDB', '3306');

    if ($conn->connect_error) {
        error_log('MySQL Connect Error (' . $conn->connect_errno . ') '
                . $conn->connect_error);
    }

    $insertQueries;

    if ( $insertType == 1 ) {
        $packetID = $conn->query("SELECT MAX(PacketID) FROM PacketInfo")->fetch_array()[0];

        if (is_null($packetID)) {
            $packetID = 1;
        } else {
            $packetID += 1;
        }

        echo "Inserting ".$packetID." as packetID\n";

        $SpeedMPH = $_GET["SpeedMPH"];
        $Gas = $_GET["Gas"];
        $Brake = $_GET["Brake"];
        $Steer = $_GET["Steer"];
        $Clutch = $_GET["Clutch"];
        $Gear = $_GET["Gear"];
        $RPM = $_GET["RPM"];
        $TurboBoost = $_GET["TurboBoost"];
        $LocalAngularVelocityX = $_GET["LocalAngularVelocityX"];
        $LocalAngularVelocityY = $_GET["LocalAngularVelocityY"];
        $LocalAngularVelocityZ = $_GET["LocalAngularVelocityZ"];
        $VelocityX = $_GET["VelocityX"];
        $VelocityY = $_GET["VelocityY"];
        $VelocityZ = $_GET["VelocityZ"];
        $WorldPositionX = $_GET["WorldPositionX"];
        $WorldPositionY = $_GET["WorldPositionY"];
        $WorldPositionZ = $_GET["WorldPositionZ"];
        $Aero_DragCoeffcient = $_GET["Aero_DragCoeffcient"];
        $Aero_LiftCoefficientFront = $_GET["Aero_LiftCoefficientFront"];
        $Aero_LiftCoefficientRear = $_GET["Aero_LiftCoefficientRear"];
        $FL_CamberRad = $_GET["FL_CamberRad"];
        $FR_CamberRad = $_GET["FR_CamberRad"];
        $RL_CamberRad = $_GET["RL_CamberRad"];
        $RR_CamberRad = $_GET["RR_CamberRad"];
        $FL_SlipAngle = $_GET["FL_SlipAngle"];
        $FR_SlipAngle = $_GET["FR_SlipAngle"];
        $RL_SlipAngle = $_GET["RL_SlipAngle"];
        $RR_SlipAngle = $_GET["RR_SlipAngle"];
        $FL_SlipRatio = $_GET["FL_SlipRatio"];
        $FR_SlipRatio = $_GET["FR_SlipRatio"];
        $RL_SlipRatio = $_GET["RL_SlipRatio"];
        $RR_SlipRatio = $_GET["RR_SlipRatio"];
        $FL_SelfAligningTorque = $_GET["FL_SelfAligningTorque"];
        $FR_SelfAligningTorque = $_GET["FR_SelfAligningTorque"];
        $RL_SelfAligningTorque = $_GET["RL_SelfAligningTorque"];
        $RR_SelfAligningTorque = $_GET["RR_SelfAligningTorque"];
        $FL_Load = $_GET["FL_Load"];
        $FR_Load = $_GET["FR_Load"];
        $RL_Load = $_GET["RL_Load"];
        $RR_Load = $_GET["RR_Load"];
        $FL_TyreSlip = $_GET["FL_TyreSlip"];
        $FR_TyreSlip = $_GET["FR_TyreSlip"];
        $RL_TyreSlip = $_GET["RL_TyreSlip"];
        $RR_TyreSlip = $_GET["RR_TyreSlip"];
        $FL_ThermalState = $_GET["FL_ThermalState"];
        $FR_ThermalState = $_GET["FR_ThermalState"];
        $RL_ThermalState = $_GET["RL_ThermalState"];
        $RR_ThermalState = $_GET["RR_ThermalState"];
        $FL_DynamicPressure = $_GET["FL_DynamicPressure"];
        $FR_DynamicPressure = $_GET["FR_DynamicPressure"];
        $RL_DynamicPressure = $_GET["RL_DynamicPressure"];
        $RR_DynamicPressure = $_GET["RR_DynamicPressure"];
        $FL_TyreDirtyLevel = $_GET["FL_TyreDirtyLevel"];
        $FR_TyreDirtyLevel = $_GET["FR_TyreDirtyLevel"];
        $RL_TyreDirtyLevel = $_GET["RL_TyreDirtyLevel"];
        $RR_TyreDirtyLevel = $_GET["RR_TyreDirtyLevel"];

        $packetInfoInsert = sprintf(    "INSERT INTO PacketInfo (PacketID, SessionID, LapID, PacketDatetime) VALUES (%s, %s, %s, NOW());", 
                                        $packetID, $sessionID, $lapID);
        $telemetryInfoInsert = sprintf( "INSERT INTO TelemetryInfo (PacketID, SpeedMPH, Gas, Brake, Steer, Clutch, Gear, RPM, TurboBoost, LocalAngularVelocityX, LocalAngularVelocityY, LocalAngularVelocityZ, VelocityX, VelocityY, VelocityZ, WorldPositionX, WorldPositionY, WorldPositionZ, Aero_DragCoeffcient, Aero_LiftCoefficientFront, Aero_LiftCoefficientRear) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s);", 
                                        $packetID, $SpeedMPH, $Gas, $Brake, $Steer, $Clutch, $Gear, $RPM, $TurboBoost, $LocalAngularVelocityX, $LocalAngularVelocityY, $LocalAngularVelocityZ, $VelocityX, $VelocityY, $VelocityZ, $WorldPositionX, $WorldPositionY, $WorldPositionZ, $Aero_DragCoeffcient, $Aero_LiftCoefficientFront, $Aero_LiftCoefficientRear);
        $tireInfoInsert = sprintf(      "INSERT INTO TireInfo (PacketID, FL_CamberRad, FR_CamberRad, RL_CamberRad, RR_CamberRad, FL_SlipAngle, FR_SlipAngle, RL_SlipAngle, RR_SlipAngle, FL_SlipRatio, FR_SlipRatio, RL_SlipRatio, RR_SlipRatio, FL_SelfAligningTorque, FR_SelfAligningTorque, RL_SelfAligningTorque, RR_SelfAligningTorque, FL_Load, FR_Load, RL_Load, RR_Load, FL_TyreSlip, FR_TyreSlip, RL_TyreSlip, RR_TyreSlip, FL_ThermalState, FR_ThermalState, RL_ThermalState, RR_ThermalState, FL_DynamicPressure, FR_DynamicPressure, RL_DynamicPressure, RR_DynamicPressure, FL_TyreDirtyLevel, FR_TyreDirtyLevel, RL_TyreDirtyLevel, RR_TyreDirtyLevel) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s);", 
                                        $packetID, $FL_CamberRad, $FR_CamberRad, $RL_CamberRad, $RR_CamberRad, $FL_SlipAngle, $FR_SlipAngle, $RL_SlipAngle, $RR_SlipAngle, $FL_SlipRatio, $FR_SlipRatio, $RL_SlipRatio, $RR_SlipRatio, $FL_SelfAligningTorque, $FR_SelfAligningTorque, $RL_SelfAligningTorque, $RR_SelfAligningTorque, $FL_Load, $FR_Load, $RL_Load, $RR_Load, $FL_TyreSlip, $FR_TyreSlip, $RL_TyreSlip, $RR_TyreSlip, $FL_ThermalState, $FR_ThermalState, $RL_ThermalState, $RR_ThermalState, $FL_DynamicPressure, $FR_DynamicPressure, $RL_DynamicPressure, $RR_DynamicPressure, $FL_TyreDirtyLevel, $FR_TyreDirtyLevel, $RL_TyreDirtyLevel, $RR_TyreDirtyLevel);

        $insertQueries = array($packetInfoInsert, $telemetryInfoInsert, $tireInfoInsert);
    } elseif ( $insertType == 2 ) {
        $existingRequest = $conn->query(sprintf("SELECT SessionID FROM LapInfo WHERE SessionID=%s AND LapID=%s", $sessionID, $lapID))->fetch_array()[0];
        if (!is_null($existingRequest)) {
            $conn->query(sprintf("DELETE FROM LapInfo WHERE SessionID=%s AND LapID=%s", $sessionID, $lapID));
        }

        $LapTime = $_GET["LapTime"];
        $DriverName = $_GET["DriverName"];
        $TrackName = $_GET["TrackName"];
        $TrackConfiguration = $_GET["TrackConfiguration"];
        $CarName = $_GET["CarName"];

        $DriverName = str_replace(",","{Comma}", $DriverName);
        $TrackName = str_replace(",", "{Comma}", $TrackName);
        $TrackConfiguration = str_replace(",", "{Comma}", $TrackConfiguration);
        $CarName = str_replace(",", "{Comma}", $CarName);

        $lapInfoInsert = sprintf(   "INSERT INTO LapInfo (SessionID, LapID, LapTime, DriverName, TrackName, TrackConfiguration, CarName) VALUES (%s, %s, %s, '%s', '%s', '%s', '%s');",
                                    $sessionID, $lapID, $LapTime, $DriverName, $TrackName, $TrackConfiguration, $CarName);

        $insertQueries = array($lapInfoInsert);
    }

    foreach ($insertQueries as $insertQuery) {
        echo $insertQuery;
        $conn->query($insertQuery);
    }
    
    $conn->close();
} else {
    echo "WrongKey";
}
?>