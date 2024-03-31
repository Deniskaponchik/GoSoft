# Version:      0.0
# STATUS:       рабочее
# Цель:         Сделать аналог MakeFile на powershell
# реализация:   
# проблемы:     
# Планы:        Протестировать
# Last Update:  

#.\Make.ps1 -device Unifi -task Run
#.\Make.ps1 -device Poly -task Run
#Set-ExecutionPolicy Unrestricted

#Внешние входные параметры для скрипта
[CmdletBinding()]
Param (
    [Parameter (Position=1)] #[Parameter (Mandatory=$true, Position=2)]
    [string]
    $device,

    [Parameter (Position=2)] #[Parameter (Mandatory=$true, Position=1)]
    #[alias("ARG","ArgumentName")]
    #[ValidatePattern("[0-9][0-9][0-9][0-9]")]
    #[ValidateLength(1,3)]
    [string]
    $task    
)
[Environment]::NewLine
Write-Host "Device : "$device
Write-Host "Task   : "$task

$VersionUnifi = "Unifi_v3.29-PROD"
$VersionPoly  = "Poly_v3.1-PROD"

$Env:GISUP_PATH = "D:/Clouds/GitHub/GoSoft" #Comment after reboot
$Env:GISUP_TIMEZONE = "+5"                  #Comment after reboot
#$TimeZone = "+5"
[Environment]::NewLine

#Выбор устройства
if (!$device){
    do {
        Write-Host "Choose device for this script "
        [Environment]::NewLine
        
        $DeviceArray = @("Unifi", "Poly", "Eltex")
        $DeviceArray | ForEach-Object {Write-Host $_}
        #Write-Host "Run" Write-Host "Build" Write-Host "Web" Write-Host "Swagger" Write-Host "GRPC" Write-Host "DeleteOldFiles" Write-Host "CreateTaskScheduler"
        [Environment]::NewLine

        $device = Read-Host "Input device "
        $device
    } until (
        #($task -ne "Run") -and ($task -ne "Build") -and ($task -ne "Web") -and ($task -ne "Swagger")
        $DeviceArray -match $device
    )

    [Environment]::NewLine
    PAUSE
}
#Выбор варианта работы скрипта
if (!$task){
    do {
        Write-Host "Choose task for this script "
        [Environment]::NewLine
        
        $TaskArray = @("Run", "Build", "Web", "Swagger", "GRPC", "DeleteOldLogs", "CreateTaskScheduler")
        $TaskArray | ForEach-Object {Write-Host $_}
        #Write-Host "Run" Write-Host "Build" Write-Host "Web" Write-Host "Swagger" Write-Host "GRPC" Write-Host "DeleteOldFiles" Write-Host "CreateTaskScheduler"
        [Environment]::NewLine

        $Task = Read-Host "Input task "
        $task
    } until (
        #($task -ne "Run") -and ($task -ne "Build") -and ($task -ne "Web") -and ($task -ne "Swagger")
        $TaskArray -match $task
    )

    [Environment]::NewLine
    PAUSE
}

function DeleteOldLogs {
    if ($device -eq "Unifi"){
        #$Env:GISUP_PATH = "D:/Clouds/GitHub/GoSoft"
        try{
            $OldLogsFiles = Get-ChildItem "$Env:GISUP_PATH/logs" -Force -ErrorAction Stop -ErrorVariable ErrGetLogs
            #1 variant
            $OldLogsFiles | Where-Object lastwritetime -lt (Get-Date).AddDays(-60) | ForEach-Object{
                #Pause
                Remove-Item -Path "$Env:GISUP_PATH/logs/$_" -verbose #-Confirm
            }
            <#2 variant
            foreach ($olf in $OldLogsFiles) {
            Write-Host $olf
            #PAUSE
            if ($olf.lastwritetime -lt (Get-Date).AddDays(-15)){   #-and $NoDelUsers -notcontains $UF.name){
                PAUSE
                try{
                    Remove-Item $olf -force -recurse -verbose #-ErrorVariable ErrDelUF2
                }catch{
                    
                }
                [Environment]::NewLine
            }   
            }
            #>
            Write-Host "Очистка файлов старых логов завершилась" -ForegroundColor Green
        }
        catch{
            Write-Host $ErrGetLogs -ForegroundColor Red
            Write-Host "Не удалось получить файлы логов" -ForegroundColor Red
        }
    }
    if ($device -eq "Poly"){
        <#ПОКА ЛОГИ НЕ БУДУТ ЛЕЖАТЬ В ПАПКЕ ЛОГОВ - НЕ ВЫПОЛНЯТЬ!!!
        #$Env:GISUP_PATH = "D:/Clouds/GitHub/GoSoft"
        try{
            $OldLogsFiles = Get-ChildItem "$Env:GISUP_PATH/cmd/poly" -Force -ErrorAction Stop -ErrorVariable ErrGetLogs
            #1 variant
            $OldLogsFiles | Where-Object lastwritetime -lt (Get-Date).AddDays(-60) | ForEach-Object{
                #Pause
                Remove-Item -Path "$Env:GISUP_PATH/cmd/poly/$_" -verbose #-Confirm
            }
            Write-Host "Очистка файлов старых логов завершилась" -ForegroundColor Green
        }
        catch{
            Write-Host $ErrGetLogs -ForegroundColor Red
            Write-Host "Не удалось получить файлы логов" -ForegroundColor Red
        }
        #>
    }
    
    [Environment]::NewLine
}

function Swag {
    swag init -g internal/controller/http/fokInterface/routes.go --parseInternal true
}

function GRPC {
    rem D:\Programms\Coding\protoc-25.1-win64\bin\protoc.exe -I api/unifi/v1 api/unifi/v1/unifi.proto ^
    rem --go_out=pkg/grpc/unifi/v1 ^
    rem --go_opt=paths=source_relative ^
    rem --go-grpc_out=pkg/grpc/unifi/v1 ^
    rem --go-grpc_opt=paths=source_relative

    rem --go_out-plugins-grpc:pkg/grpc/server ^
    rem --go_out        сюда складывается результат protoc
    rem --go_opt        сгенерированные файлы будут использовать тот же пакет, что и proto-файлы. source_relative означает, что выходные файлы будут иметь тот же пакет, что и исходные .proto файлы
    rem --go-grpc_out   куда складывать go grpc код
    rem --go-grpc_opt   как создавать имена пакетов для gRPC
}

function Build {
    if ($device -eq "Unifi"){
        cd cmd/unifi
        #go build -o ..\..\bin\Unifi_v3.28-PROD -ldflags="-s -w"
         go build -o ..\..\bin\$VersionUnifi -ldflags="-s -w"
         cd ../..
    }
    if ($device -eq "Poly"){
        cd cmd/poly
        #go build -o ..\..\bin\Poly_v3.0-TEST -ldflags="-s -w"
         go build -o ..\..\bin\$VersionPoly -ldflags="-s -w"
         cd ../..
    }
}

function Run {
    if ($device -eq "Unifi"){
        DeleteOldLogs

        #cmd.exe /c "D:/Clouds/GitHub/GoSoft/bin/$VersionUnifi -mode PROD -time $TimeZone"
        #cmd.exe /c "$Env:GISUP_PATH/bin/$VersionUnifi -mode PROD -time $Env:GISUP_TIMEZONE"

        start "$Env:GISUP_PATH/bin/$VersionUnifi" -ArgumentList "-mode PROD -time $Env:GISUP_TIMEZONE" -NoNewWindow
    }
    if ($device -eq "Poly"){
        DeleteOldLogs

        #cmd.exe /c "D:/Clouds/GitHub/GoSoft/bin/$VersionPoly -mode PROD -time $TimeZone"
        #cmd.exe /c "$Env:GISUP_PATH/bin/$VersionPoly -mode PROD -time $Env:GISUP_TIMEZONE"

        start "$Env:GISUP_PATH/bin/$VersionPoly" -ArgumentList "-mode PROD -time $Env:GISUP_TIMEZONE" -NoNewWindow
    }
}

function Web {
    if ($device -eq "Unifi"){
        DeleteOldLogs

        #$Env:GISUP_PATH = "D:/Clouds/GitHub/GoSoft"
        #$VersionUnifi = "Unifi_v3.29-PROD"

        #cmd.exe /c "/B bin/$VersionUnifi -mode WEB -time $TimeZone"
        #cmd.exe /c "D:/Clouds/GitHub/GoSoft/bin/$VersionUnifi -mode WEB -time $TimeZone"
        #. "C:\Program Files\Go\bin\go.exe" D:/Clouds/GitHub/GoSoft/bin/$VersionUnifi -mode WEB -time $TimeZone
        #cmd.exe /c "$Env:GISUP_PATH/bin/$VersionUnifi -mode WEB -time $Env:GISUP_TIMEZONE"

        #Start-Process cmd.exe -ArgumentList "/B bin/$VersionUnifi -mode WEB -time $TimeZone" -NoNewWindow
        #Start-Process "D:/Clouds/GitHub/GoSoft/bin/Unifi_v3.29-PROD" -ArgumentList "-mode WEB -time +5" -NoNewWindow
        start "$Env:GISUP_PATH/bin/$VersionUnifi" -ArgumentList "-mode WEB -time $Env:GISUP_TIMEZONE" -NoNewWindow
    }
}


function CreateTaskScheduler{
    if ($device -eq "Unifi"){
        #https://blog.netwrix.com/2018/07/03/how-to-automate-powershell-scripts-with-task-scheduler/
        $Trigger= New-ScheduledTaskTrigger -At 10:00am –Daily
        $User= "NT AUTHORITYSYSTEM" # Specify the account to run the script
        # Specify what program to run and with its parameters
        $Action= New-ScheduledTaskAction -Execute "PowerShell.exe" -Argument "C:PSStartupScript.ps1" 
        # Specify the name of the task
        Register-ScheduledTask -TaskName "MonitorGroupMembership" -Trigger $Trigger -User $User -Action $Action -RunLevel Highest –Force
    }
    if ($device -eq "Poly"){

    }
}


switch($task){
    "Run"{Run}
    "Build"{Build}  # -ipadd $ip
    "Web"{Web}
    "Swagger"{Swagger}
    "GRPC"{GRPC}
    "DeleteOldFiles"{DeleteOldFiles}
    "CreateTaskScheduler"{CreateTaskScheduler}
Default {"Run"}
}
[Environment]::NewLine




