@echo off
setlocal enabledelayedexpansion

:: Enable ANSI escape sequence support (Windows 10+)
for /F "delims=#" %%E in ('"prompt #$E# & for %%E in (1) do rem"') do set "ESC=%%E"

:: Set console colors
set "GREEN=%ESC%[0;32m"
set "BLUE=%ESC%[0;34m"
set "YELLOW=%ESC%[1;33m"
set "RED=%ESC%[0;31m"
set "NC=%ESC%[0m"

echo  %BLUE%================================================%NC%
echo         Project EduNexus - One-Click Start
echo  %BLUE%================================================%NC%

:: Check dependencies
echo %YELLOW%Checking dependencies...%NC%

where go >nul 2>nul
if %errorlevel% neq 0 (
    echo %RED%Error: 'go' command not found. Please install Go environment first.%NC%
    exit /b 1
)

where npm >nul 2>nul
if %errorlevel% neq 0 (
    echo %RED%Error: 'npm' command not found. Please install Node.js and npm first.%NC%
    exit /b 1
)

:: Check project directory
if not exist "edunexus\" (
    if not exist "wails.json" (
        echo %RED%Error: Please run this script in the project root or edunexus directory.%NC%
        exit /b 1
    )
) else (
    cd edunexus
)

echo.
echo %GREEN%[1/2] Building Frontend...%NC%
cd frontend
call npm install
call npm run build
cd ..

echo.
echo %GREEN%[2/2] Starting EduNexus Client...%NC%
echo %YELLOW%Starting with go run -tags dev . (Press Ctrl+C to stop)%NC%

:: Run wails project
go run -tags dev .

endlocal
