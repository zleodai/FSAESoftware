<?php
$env = parse_ini_file('.env');
$user = $env["user"];
$password = $env["password"];
$host = $env["host"];
$port = $env["port"];
$dbname = $env["dbname"];

$privateKey = "945";

//Query Types
//1 = query PacketInfo from PackedID to PacketID
//2 = query LapInfo from SessionID and LapID
//3 = query TelemetryInfo from PacketID to PacketID
//4 = query TireInfo from PacketID to PacketID

$publicKey = (int)$_GET["publicKey"];
$queryType = (int)$_GET["queryType"];
$x = (int)$_GET["x"];
$y = (int)$_GET["y"];

$calculatedKey = (int)hash('sha256', $queryType * $privateKey);

if ($publicKey == $calculatedKey) {
    $conn = pg_connect(sprintf("user=%s password=%s host=%s port=%s dbname=%s", $user, $password, $host, $port, $dbname));
    if (!$conn) {
        error_log('Connection to Postgresql database failed');
        exit;
    }

    $query; 

    switch ($queryType) {
        case 2:
            $query = sprintf("SELECT SessionID, LapID, LapTime, DriverName, TrackName, TrackConfiguration, CarName FROM LapInfo WHERE SessionID = %d AND LapID = %d;", $x, $y); 
            break;
        case 3:
            $query = sprintf("SELECT PacketID, SpeedMPH, Gas, Brake, Steer, Clutch, Gear, RPM, TurboBoost, LocalAngularVelocityX, LocalAngularVelocityY, LocalAngularVelocityZ, VelocityX, VelocityY, VelocityZ, WorldPositionX, WorldPositionY, WorldPositionZ, Aero_DragCoeffcient, Aero_LiftCoefficientFront, Aero_LiftCoefficientRear FROM TelemetryInfo WHERE PacketID BETWEEN %d AND %d;", $x, $y); 
            break;
        case 4:
            $query = sprintf("SELECT PacketID, FL_CamberRad, FR_CamberRad, RL_CamberRad, RR_CamberRad, FL_SlipAngle, FR_SlipAngle, RL_SlipAngle, RR_SlipAngle, FL_SlipRatio, FR_SlipRatio, RL_SlipRatio, RR_SlipRatio, FL_SelfAligningTorque, FR_SelfAligningTorque, RL_SelfAligningTorque, RR_SelfAligningTorque, FL_Load, FR_Load, RL_Load, RR_Load, FL_TyreSlip, FR_TyreSlip, RL_TyreSlip, RR_TyreSlip, FL_ThermalState, FR_ThermalState, RL_ThermalState, RR_ThermalState, FL_DynamicPressure, FR_DynamicPressure, RL_DynamicPressure, RR_DynamicPressure, FL_TyreDirtyLevel, FR_TyreDirtyLevel, RL_TyreDirtyLevel, RR_TyreDirtyLevel FROM TireInfo WHERE PacketID BETWEEN %d AND %d;", $x, $y); 
            break;
        default:
            $query = sprintf("SELECT PacketID, SessionID, LapID, PacketDatetime FROM PacketInfo WHERE PacketID BETWEEN %d AND %d", $x, $y); 
            break;
    }
    // echo sprintf("Querying from packets %d - %d ...\n", $x, $y);

    $result = pg_query($conn, $query);
    if (!$result) {
        error_log('Query Failed');
        exit;
    }

    // printf("Select returned %d rows.\n", $result->num_rows);

    while ($packet = pg_fetch_row($result)) {
        switch ($queryType) {
            case 2:
                echo $packet;
                $SessionID = $packet[0];
                $LapID = $packet[1];
                $LapTime = $packet[2];
                $DriverName = $packet[3];
                $TrackName = $packet[4];
                $TrackConfiguration = $packet[5];
                $CarName = $packet[6];
                echo $SessionID.",".$LapID.",".$LapTime.",".$DriverName.",".$TrackName.",".$TrackConfiguration.",".$CarName."\n";
                break;
            case 3:
                $PacketID = $packet[0];
                $SpeedMPH = $packet[1];
                $Gas = $packet[2];
                $Brake = $packet[3];
                $Steer = $packet[4];
                $Clutch = $packet[5];
                $Gear = $packet[6];
                $RPM = $packet[7];
                $TurboBoost = $packet[8];
                $LocalAngularVelocityX = $packet[9];
                $LocalAngularVelocityY = $packet[10];
                $LocalAngularVelocityZ = $packet[11];
                $VelocityX = $packet[12];
                $VelocityY = $packet[13];
                $VelocityZ = $packet[14];
                $WorldPositionX = $packet[15];
                $WorldPositionY = $packet[16];
                $WorldPositionZ = $packet[17];
                $Aero_DragCoeffcient = $packet[18];
                $Aero_LiftCoefficientFront = $packet[19];
                $Aero_LiftCoefficientRear = $packet[20];
                echo $PacketID.",".$SpeedMPH.",".$Gas.",".$Brake.",".$Steer.",".$Clutch.",".$Gear.",".$RPM.",".$TurboBoost.",".$LocalAngularVelocityX.",".$LocalAngularVelocityY.",".$LocalAngularVelocityZ.",".$VelocityX.",".$VelocityY.",".$VelocityZ.",".$WorldPositionX.",".$WorldPositionY.",".$WorldPositionZ.",".$Aero_DragCoeffcient.",".$Aero_LiftCoefficientFront.",".$Aero_LiftCoefficientRear."\n";
                break;
            case 4:
                $PacketID = $packet[0];
                $FL_CamberRad = $packet[1];
                $FR_CamberRad = $packet[2];
                $RL_CamberRad = $packet[3];
                $RR_CamberRad = $packet[4];
                $FL_SlipAngle = $packet[5];
                $FR_SlipAngle = $packet[6];
                $RL_SlipAngle = $packet[7];
                $RR_SlipAngle = $packet[8];
                $FL_SlipRatio = $packet[9];
                $FR_SlipRatio = $packet[10];
                $RL_SlipRatio = $packet[11];
                $RR_SlipRatio = $packet[12];
                $FL_SelfAligningTorque = $packet[13];
                $FR_SelfAligningTorque = $packet[14];
                $RL_SelfAligningTorque = $packet[15];
                $RR_SelfAligningTorque = $packet[16];
                $FL_Load = $packet[17];
                $FR_Load = $packet[18];
                $RL_Load = $packet[19];
                $RR_Load = $packet[20];
                $FL_TyreSlip = $packet[21];
                $FR_TyreSlip = $packet[22];
                $RL_TyreSlip = $packet[23];
                $RR_TyreSlip = $packet[24];
                $FL_ThermalState = $packet[25];
                $FR_ThermalState = $packet[26];
                $RL_ThermalState = $packet[27];
                $RR_ThermalState = $packet[28];
                $FL_DynamicPressure = $packet[29];
                $FR_DynamicPressure = $packet[30];
                $RL_DynamicPressure = $packet[31];
                $RR_DynamicPressure = $packet[32];
                $FL_TyreDirtyLevel = $packet[33];
                $FR_TyreDirtyLevel = $packet[34];
                $RL_TyreDirtyLevel = $packet[35];
                $RR_TyreDirtyLevel = $packet[36];
                echo $PacketID.",".$FL_CamberRad.",".$FR_CamberRad.",".$RL_CamberRad.",".$RR_CamberRad.",".$FL_SlipAngle.",".$FR_SlipAngle.",".$RL_SlipAngle.",".$RR_SlipAngle.",".$FL_SlipRatio.",".$FR_SlipRatio.",".$RL_SlipRatio.",".$RR_SlipRatio.",".$FL_SelfAligningTorque.",".$FR_SelfAligningTorque.",".$RL_SelfAligningTorque.",".$RR_SelfAligningTorque.",".$FL_Load.",".$FR_Load.",".$RL_Load.",".$RR_Load.",".$FL_TyreSlip.",".$FR_TyreSlip.",".$RL_TyreSlip.",".$RR_TyreSlip.",".$FL_ThermalState.",".$FR_ThermalState.",".$RL_ThermalState.",".$RR_ThermalState.",".$FL_DynamicPressure.",".$FR_DynamicPressure.",".$RL_DynamicPressure.",".$RR_DynamicPressure.",".$FL_TyreDirtyLevel.",".$FR_TyreDirtyLevel.",".$RL_TyreDirtyLevel.",".$RR_TyreDirtyLevel."\n";
                break;
            default:
                $PacketID = $packet[0];
                $SessionID = $packet[1];
                $LapID = $packet[2];
                $PacketDateTime = $packet[3];
                echo $PacketID.",".$SessionID.",".$LapID.",".$PacketDateTime."\n";
                break;
        }
    }
} else {
    echo "WrongKey";
}
?>