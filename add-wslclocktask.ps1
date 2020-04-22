$taskCommand = "powershell.exe -noprofile -windowstyle Hidden -Command '" + $PSScriptRoot + "\update-clock.ps1'"

schtasks /Create /TN wsl-clock /TR $taskCommand /SC ONEVENT /EC System /MO "*[System[Provider[@Name='Microsoft-Windows-Kernel-Power'] and (EventID=107 or EventID=507)]]" /F