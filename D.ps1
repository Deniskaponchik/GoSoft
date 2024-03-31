# Version:      0.0
# STATUS:       Не протестировано
# Цель:         
# реализация:   
# проблемы:     
# Планы:        Протестировать
# Last Update:  
#

#.\D.ps1 -task Web
#Set-ExecutionPolicy Unrestricted

#Внешние входные параметры для скрипта
[CmdletBinding()]
Param (
    [Parameter (Position=1)] #[Parameter (Mandatory=$true, Position=1)]
    #[alias("ARG","ArgumentName")]
    #[ValidatePattern("[0-9][0-9][0-9][0-9]")]
    #[ValidateLength(1,3)]
    [string]
    $task,

    [Parameter (Position=2)] #[Parameter (Mandatory=$true, Position=2)]
    [string]
    $hostname
)
Write-Host "Task : "$task
Write-Host "Hostname : "$hostname
#$Env:GISUP_PATH
$ScriptVersion = "Unifi_v3.29-PROD"
$TimeZone = "+5"
[Environment]::NewLine 

#Выбор варианта работы скрипта
if (!$task){
    do {
        Write-Host "Choose task for this script "
        [Environment]::NewLine

        
        Write-Host "Run"
        Write-Host "Build"
        Write-Host "Web"
        Write-Host "Swagger"
        Write-Host "GRPC"
        [Environment]::NewLine

        $Task = Read-Host "Input task "
        $task
    } while (
        ($task -ne "Run") -and ($task -ne "Build") -and ($task -ne "Web") -and ($task -ne "Swagger")
    )

    [Environment]::NewLine
    PAUSE
}

<#if (!$hostname){
    if ($language -eq "RUS"){
        Write-Host "Укажи имя компьютера" -ForegroundColor RED
        Write-Host "Имя должно начинаться с VCSXX-ConfRoomName" -ForegroundColor RED
        Write-Host "XX           - двухбуквенный код региона из AD" -ForegroundColor RED
        Write-Host "ConfRoomName - имя переговорной комнаты" -ForegroundColor RED
        Write-Host "Длина всего имени компьютера не больше 15 символов" -ForegroundColor RED
        Write-Host "Пример:  VCSIR-SELENGA" -ForegroundColor RED
    }else{
        Write-Host "Input computer name" -ForegroundColor RED
        Write-Host "name should match with mask:" -ForegroundColor RED
        Write-Host "VCSXX-ConfRoomName" -ForegroundColor RED
        Write-Host "where" -ForegroundColor RED
        Write-Host "XX           - 2 letter code from AD" -ForegroundColor RED
        Write-Host "ConfRoomName - Name of conf room" -ForegroundColor RED
        Write-Host "Name length must contains max 15 letters" -ForegroundColor RED
        Write-Host "Example:  VCSIR-SELENGA" -ForegroundColor RED
    }
    [Environment]::NewLine

    do {
        $hostname = Read-Host "Hostname "
    } while (
        ($hostname -eq '') -or ($hostname -notmatch "[v][c][s][A-z][A-z][-].\B" -and $hostname.Length -le 15)
    )
    $hostname

    [Environment]::NewLine
    PAUSE
}
#>

function Run {
     cmd.exe /c "D:/Clouds/GitHub/GoSoft/bin/$ScriptVersion -mode PROD -time $TimeZone"
    #cmd.exe /c "$Env:GISUP_PATH/bin/$ScriptVersion -mode PROD -time $Env:GISUP_TIMEZONE"
    #"$Env:GISUP_HTTP_URL/bin/"
}

function Build {

}

function Web {
    #$ScriptVersion = "Unifi_v3.29-PROD"
    #cmd.exe /c "/B bin/$ScriptVersion -mode WEB -time $TimeZone"
     cmd.exe /c "D:/Clouds/GitHub/GoSoft/bin/$ScriptVersion -mode WEB -time $TimeZone"
    #cmd.exe /c "$Env:GISUP_PATH/bin/$ScriptVersion -mode WEB -time $Env:GISUP_TIMEZONE"

    #Start-Process cmd.exe -ArgumentList "/B bin/$ScriptVersion -mode WEB -time $TimeZone" -NoNewWindow
    #start bin/$ScriptVersion -ArgumentList "-mode WEB -time $TimeZone"
    #go run ./bin/$ScriptVersion -mode WEB -time $TimeZone
}


function DeleteOldFiles {
    #(Get-Content Input.json) -replace '"(\d+),(\d{1,})"', '$1.$2' `
    #-replace 'second regex', 'second replacement' | 
    #Out-File output.json
    #"HostMetadata=Region=Нижний Новгород:UserLogin=roman.novotorov:RoomName=Бежин Луг:IsVcs=true:VcsType=Lenovo"

    #(Get-Content -Path "D:\Test.txt" -ErrorVariable ErrGetZabbixConfig -ErrorAction Stop) `
    #-replace 'TEST', 'VCS' | 
    #Out-File "D:\Test.txt"

    #$TestTxt = Get-Content -Path "D:\Test.txt" -ErrorVariable ErrGetZabbixConfig -ErrorAction Stop
    #$TestTxt.replace('TEST', 'VCS').replace('tst', 'txt') | 
    #Out-File "D:\Test.txt"

    <#
    (Get-Content -Path $PathZabbixConf -ErrorVariable ErrGetZabbixConfig -ErrorAction Stop) `
    -replace '#Server=WillChangeFromScript', "Server=$ZabbixServer"`
    -replace '#ServerActive=WillChangeFromScript', "ServerActive=$ZabbixServer"`
    -replace '#Hostname=WillChangeFromScript', "Hostname=$PcNewName"`
    -replace '#HostMetadata=WillChangeFromScript', "HostMetadata=Region=$Region:UserLogin=$UserLogin:RoomName=$RoomName:IsVcs=true:VcsType=$VcsType" | 
    Out-File $PathZabbixConf #zabbix_agentd.conf
    #>

    <#
    (Get-Content -Path $PathZabbixConf -Encoding UTF8 -ErrorVariable ErrGetZabbixConfig -ErrorAction Stop).replace('#Server=WillChangeFromScript', "Server=$ZabbixServer").replace('#ServerActive=WillChangeFromScript', "ServerActive=$ZabbixServer").replace('#Hostname=WillChangeFromScript', "Hostname=$PcNewName").replace('#HostMetadata=WillChangeFromScript', "HostMetadata=Region=$Region;UserLogin=$UserLogin;RoomName=$RoomName;IsVcs=true;VcsType=$VcsType;") | 
    Out-File $PathZabbixConf -Encoding UTF8 #ASCII
    #>

    #https://stackoverflow.com/questions/5596982/using-powershell-to-write-a-file-in-utf-8-without-the-bom
    $MyRawString = Get-Content -Raw $PathZabbixConf
    $MyRawStringReplace = $MyRawString.
    replace('#Server=WillChangeFromScript', "Server=$ZabbixServer").
    replace('#ServerActive=WillChangeFromScript', "ServerActive=$ZabbixServer").
    replace('#Hostname=WillChangeFromScript', "Hostname=$hostname").
    replace('#HostMetadata=WillChangeFromScript', "HostMetadata=Region=$Region;UserLogin=$UserLogin;RoomName=$RoomName;IsVcs=true;VcsType=$VcsType;")
    $Utf8NoBomEncoding = New-Object System.Text.UTF8Encoding $False
    [System.IO.File]::WriteAllLines($PathZabbixConf, $MyRawStringReplace, $Utf8NoBomEncoding)


    [Environment]::NewLine
    if ($language -eq "RUS"){
        Write-Host "Изменения в файл конфигурации Zabbix-агента внесены" -ForegroundColor GREEN
    }else{
        Write-Host "The Zabbix Agent conf.file was changed successfully" -ForegroundColor GREEN
    }
   
}

function CreateTaskScheduler{
    #https://blog.netwrix.com/2018/07/03/how-to-automate-powershell-scripts-with-task-scheduler/
    $Trigger= New-ScheduledTaskTrigger -At 10:00am –Daily
    $User= "NT AUTHORITYSYSTEM" # Specify the account to run the script
    # Specify what program to run and with its parameters
    $Action= New-ScheduledTaskAction -Execute "PowerShell.exe" -Argument "C:PSStartupScript.ps1" 
    # Specify the name of the task
    Register-ScheduledTask -TaskName "MonitorGroupMembership" -Trigger $Trigger -User $User -Action $Action -RunLevel Highest –Force 
}


switch($task){
    "Run"{Run}
    "Build"{Build}  # -ipadd $ip
    "Web"{Web}
    "Swagger"{Swagger}
    "GRPC"{GRPC}
Default {"EMPTY"}
}
[Environment]::NewLine




