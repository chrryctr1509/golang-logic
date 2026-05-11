$MYSQL = "C:\laragon\bin\mysql\mysql-8.0.40-winx64\bin\mysql.exe"
$DB = "wallet_db"
$PASS = "123456"

$output = & $MYSQL -u root -p$PASS $DB -e "SHOW CREATE TABLE users;" 2>$null
Write-Output $output
$output = & $MYSQL -u root -p$PASS $DB -e "SHOW CREATE TABLE transactions;" 2>$null
Write-Output $output
$output = & $MYSQL -u root -p$PASS $DB -e "SELECT user_id, first_name, last_name, phone_number, address, pin, balance, DATE_FORMAT(created_date, '%Y-%m-%d %H:%i:%s') as created_date FROM users;" 2>$null
Write-Output $output
$output = & $MYSQL -u root -p$PASS $DB -e "SELECT transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, DATE_FORMAT(created_date, '%Y-%m-%d %H:%i:%s') as created_date FROM transactions;" 2>$null
Write-Output $output