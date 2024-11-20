$newPath = ";$(realpath ./FSAESoftwareBinaryImports/pgsql/bin)"
$oldPath = [Environment]::GetEnvironmentVariable('PATH', 'User');
$newPath = $newPath.Replace("/", "\")
$newPath = $newpath.Replace("\c\", "C:\")
setx Path "$oldPath$newPath"