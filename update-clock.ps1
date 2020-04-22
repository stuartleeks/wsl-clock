function Log($message) {
    $output = (Get-Date).ToUniversalTime().ToString("u") + "`t$message"
    Add-Content -Path "~/.wsl-clock.log" -Value $output
}

Log "********************************"
Log "*** Update WSL clock starting..."

$runningDistroCount = wsl --list --running --quiet | 
        Where-Object {$_ -ne ""} | 
        Measure-Object | 
        Select-Object -ExpandProperty Count

if ($runningDistroCount -eq 0){
    Log "No Distros - quitting"
    exit 1
}

$originalDate=wsl sh -c "date -Iseconds"

Log "Performing reset..."
$result = wsl -u root sh -c "hwclock -s" 2>&1
$success = $?
if (-not $success){
    Log "reset failed:"
    Log $result
    exit 2
} else {
    $newDate=wsl bash -c "date -Iseconds"
    Log "clock reset"
    Log "OriginalDate:$originalDate"
    Log "NewDate:     $newDate"
    exit 0
}