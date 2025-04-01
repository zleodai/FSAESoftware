<?php
$serverName = "localhost"; 
$uid = "TelemetryDBUser";   
$pwd = "123";  
$databaseName = "TelemetryDB"; 
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
    $conn = new mysqli('localhost', 'TelemetryDBUser', '123', 'TelemetryDB', '3306');

    if ($conn->connect_error) {
        error_log('MySQL Connect Error (' . $conn->connect_errno . ') '
                . $conn->connect_error);
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

    $result = $conn->query($query);

    // printf("Select returned %d rows.\n", $result->num_rows);

    if ($result->num_rows > 0) {
        switch ($queryType) {
            case 2:
                foreach ($result as $packet) {
                    foreach ($result as $packet) {
                        $SessionID = $packet["SessionID"];
                        $LapID = $packet["LapID"];
                        $LapTime = $packet["LapTime"];
                        $DriverName = $packet["DriverName"];
                        $TrackName = $packet["TrackName"];
                        $TrackConfiguration = $packet["TrackConfiguration"];
                        $CarName = $packet["CarName"];
                        echo $SessionID.",".$LapID.",".$LapTime.",".$DriverName.",".$TrackName.",".$TrackConfiguration.",".$CarName."\n";
                    }
                }
                break;
            case 3:
                foreach ($result as $packet) { 
                    $PacketID = $packet["PacketID"];
                    $SpeedMPH = $packet["SpeedMPH"];
                    $Gas = $packet["Gas"];
                    $Brake = $packet["Brake"];
                    $Steer = $packet["Steer"];
                    $Clutch = $packet["Clutch"];
                    $Gear = $packet["Gear"];
                    $RPM = $packet["RPM"];
                    $TurboBoost = $packet["TurboBoost"];
                    $LocalAngularVelocityX = $packet["LocalAngularVelocityX"];
                    $LocalAngularVelocityY = $packet["LocalAngularVelocityY"];
                    $LocalAngularVelocityZ = $packet["LocalAngularVelocityZ"];
                    $VelocityX = $packet["VelocityX"];
                    $VelocityY = $packet["VelocityY"];
                    $VelocityZ = $packet["VelocityZ"];
                    $WorldPositionX = $packet["WorldPositionX"];
                    $WorldPositionY = $packet["WorldPositionY"];
                    $WorldPositionZ = $packet["WorldPositionZ"];
                    $Aero_DragCoeffcient = $packet["Aero_DragCoeffcient"];
                    $Aero_LiftCoefficientFront = $packet["Aero_LiftCoefficientFront"];
                    $Aero_LiftCoefficientRear = $packet["Aero_LiftCoefficientRear"];
                    echo $PacketID.",".$SpeedMPH.",".$Gas.",".$Brake.",".$Steer.",".$Clutch.",".$Gear.",".$RPM.",".$TurboBoost.",".$LocalAngularVelocityX.",".$LocalAngularVelocityY.",".$LocalAngularVelocityZ.",".$VelocityX.",".$VelocityY.",".$VelocityZ.",".$WorldPositionX.",".$WorldPositionY.",".$WorldPositionZ.",".$Aero_DragCoeffcient.",".$Aero_LiftCoefficientFront.",".$Aero_LiftCoefficientRear."\n";
                }
                break;
            case 4:
                foreach ($result as $packet) {
                    $PacketID = $packet["PacketID"];
                    $FL_CamberRad = $packet["FL_CamberRad"];
                    $FR_CamberRad = $packet["FR_CamberRad"];
                    $RL_CamberRad = $packet["RL_CamberRad"];
                    $RR_CamberRad = $packet["RR_CamberRad"];
                    $FL_SlipAngle = $packet["FL_SlipAngle"];
                    $FR_SlipAngle = $packet["FR_SlipAngle"];
                    $RL_SlipAngle = $packet["RL_SlipAngle"];
                    $RR_SlipAngle = $packet["RR_SlipAngle"];
                    $FL_SlipRatio = $packet["FL_SlipRatio"];
                    $FR_SlipRatio = $packet["FR_SlipRatio"];
                    $RL_SlipRatio = $packet["RL_SlipRatio"];
                    $RR_SlipRatio = $packet["RR_SlipRatio"];
                    $FL_SelfAligningTorque = $packet["FL_SelfAligningTorque"];
                    $FR_SelfAligningTorque = $packet["FR_SelfAligningTorque"];
                    $RL_SelfAligningTorque = $packet["RL_SelfAligningTorque"];
                    $RR_SelfAligningTorque = $packet["RR_SelfAligningTorque"];
                    $FL_Load = $packet["FL_Load"];
                    $FR_Load = $packet["FR_Load"];
                    $RL_Load = $packet["RL_Load"];
                    $RR_Load = $packet["RR_Load"];
                    $FL_TyreSlip = $packet["FL_TyreSlip"];
                    $FR_TyreSlip = $packet["FR_TyreSlip"];
                    $RL_TyreSlip = $packet["RL_TyreSlip"];
                    $RR_TyreSlip = $packet["RR_TyreSlip"];
                    $FL_ThermalState = $packet["FL_ThermalState"];
                    $FR_ThermalState = $packet["FR_ThermalState"];
                    $RL_ThermalState = $packet["RL_ThermalState"];
                    $RR_ThermalState = $packet["RR_ThermalState"];
                    $FL_DynamicPressure = $packet["FL_DynamicPressure"];
                    $FR_DynamicPressure = $packet["FR_DynamicPressure"];
                    $RL_DynamicPressure = $packet["RL_DynamicPressure"];
                    $RR_DynamicPressure = $packet["RR_DynamicPressure"];
                    $FL_TyreDirtyLevel = $packet["FL_TyreDirtyLevel"];
                    $FR_TyreDirtyLevel = $packet["FR_TyreDirtyLevel"];
                    $RL_TyreDirtyLevel = $packet["RL_TyreDirtyLevel"];
                    $RR_TyreDirtyLevel = $packet["RR_TyreDirtyLevel"];
                    echo $PacketID.",".$FL_CamberRad.",".$FR_CamberRad.",".$RL_CamberRad.",".$RR_CamberRad.",".$FL_SlipAngle.",".$FR_SlipAngle.",".$RL_SlipAngle.",".$RR_SlipAngle.",".$FL_SlipRatio.",".$FR_SlipRatio.",".$RL_SlipRatio.",".$RR_SlipRatio.",".$FL_SelfAligningTorque.",".$FR_SelfAligningTorque.",".$RL_SelfAligningTorque.",".$RR_SelfAligningTorque.",".$FL_Load.",".$FR_Load.",".$RL_Load.",".$RR_Load.",".$FL_TyreSlip.",".$FR_TyreSlip.",".$RL_TyreSlip.",".$RR_TyreSlip.",".$FL_ThermalState.",".$FR_ThermalState.",".$RL_ThermalState.",".$RR_ThermalState.",".$FL_DynamicPressure.",".$FR_DynamicPressure.",".$RL_DynamicPressure.",".$RR_DynamicPressure.",".$FL_TyreDirtyLevel.",".$FR_TyreDirtyLevel.",".$RL_TyreDirtyLevel.",".$RR_TyreDirtyLevel."\n";
                }
                break;
            default:
                foreach ($result as $packet) {
                    $PacketID = $packet["PacketID"];
                    $SessionID = $packet["SessionID"];
                    $LapID = $packet["LapID"];
                    $PacketDateTime = $packet["PacketDatetime"];
                    echo $PacketID.",".$SessionID.",".$LapID.",".$PacketDateTime."\n";
                }
                break;
        }
    } else {
        echo "Got 0 rows";
    }
    
    $conn->close();
} else {
    echo "WrongKey";
}
?>