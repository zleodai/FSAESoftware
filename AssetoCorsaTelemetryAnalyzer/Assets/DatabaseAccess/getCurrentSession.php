<?php
$env = parse_ini_file('.env');
$user = $env["user"];
$password = $env["password"];
$host = $env["host"];
$port = $env["port"];
$dbname = $env["dbname"];

$privateKey = "945";

$publicKey = (int)$_GET["publicKey"];

$calculatedKey = (int)hash('sha256', $privateKey);

if ($publicKey == $calculatedKey) {
    $conn = pg_connect(sprintf("user=%s password=%s host=%s port=%s dbname=%s", $user, $password, $host, $port, $dbname));
    if (!$conn) {
        error_log('Connection to Postgresql database failed');
        exit;
    }

    $query = "SELECT MAX(SessionID) FROM PacketInfo"; 

    $result = pg_query($conn, $query);
    if (!$result) {
        error_log('Query Failed');
        exit;
    }

    echo (int) pg_fetch_row($result)[0] + 1;
} else {
    echo "WrongKey";
}
?>