<?php
$serverName = "localhost"; 
$uid = "TelemetryDBUser";   
$pwd = "123";  
$databaseName = "TelemetryDB"; 
$privateKey = "945";

$publicKey = (int)$_GET["publicKey"];

$calculatedKey = (int)hash('sha256', $privateKey);

if ($publicKey == $calculatedKey) {
    $conn = new mysqli('localhost', 'TelemetryDBUser', '123', 'TelemetryDB', '3306');

    if ($conn->connect_error) {
        error_log('MySQL Connect Error (' . $conn->connect_errno . ') '
                . $conn->connect_error);
    }

    $query = "SELECT MAX(SessionID) FROM PacketInfo"; 

    $result = $conn->query($query);

    echo (int) $result->fetch_array()[0] + 1;
    
    $conn->close();
} else {
    echo "WrongKey";
}
?>