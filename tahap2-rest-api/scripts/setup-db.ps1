$mysql = Get-ChildItem -Path "C:\Program Files\MySQL" -Recurse -Filter "mysql.exe" -ErrorAction SilentlyContinue | Select-Object -First 1
if ($mysql) {
    Write-Host "Found MySQL at: $($mysql.FullName)"
    & $mysql.FullName -u root -e "CREATE DATABASE IF NOT EXISTS wallet_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    if ($LASTEXITCODE -eq 0) { Write-Host "Database created successfully" }
    else { Write-Host "MySQL command failed with exit code: $LASTEXITCODE" }
} else {
    Write-Host "MySQL not found in Program Files"
}