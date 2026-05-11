$mysql = Get-ChildItem -Path "C:\" -Recurse -Filter "mysql.exe" -ErrorAction SilentlyContinue -Force | Select-Object -First 3
$mysql | ForEach-Object { Write-Host $_.FullName }
if (-not $mysql) { Write-Host "mysql.exe not found on C: drive" }