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
    pause
    goto :EXIT_WITH_ERROR
)

where npm >nul 2>nul
if %errorlevel% neq 0 (
    echo %RED%Error: 'npm' command not found. Please install Node.js and npm first.%NC%
    pause
    goto :EXIT_WITH_ERROR
)

:: Check project directory
if not exist "edunexus\" (
    if not exist "wails.json" (
        echo %RED%Error: Please run this script in the project root or edunexus directory.%NC%
        pause
        goto :EXIT_WITH_ERROR
    )
) else (
    cd edunexus
    if !errorlevel! neq 0 (
        echo %RED%Error: Failed to change directory to 'edunexus'.%NC%
        pause
        goto :EXIT_WITH_ERROR
    )
)

echo.
echo %GREEN%[1/2] Building Frontend...%NC%
cd frontend
if %errorlevel% neq 0 (
    echo %RED%Error: Failed to change directory to 'frontend'.%NC%
    pause
    goto :EXIT_WITH_ERROR
)

call npm install
if %errorlevel% neq 0 (
    echo %RED%Error: 'npm install' failed.%NC%
    pause
    goto :EXIT_WITH_ERROR
)
call npm run build
if %errorlevel% neq 0 (
    echo %RED%Error: 'npm run build' failed.%NC%
    pause
    goto :EXIT_WITH_ERROR
)
cd ..
if %errorlevel% neq 0 (
    echo %RED%Error: Failed to change directory back from 'frontend'.%NC%
    pause
    goto :EXIT_WITH_ERROR
)

echo.
echo %GREEN%[2/2] Starting EduNexus Client...%NC%

:: Check if Wails is installed
where wails >nul 2>nul
if %errorlevel% equ 0 (
    echo %YELLOW%Wails CLI found. Starting with 'wails dev' (Press Ctrl+C to stop)%NC%
    call wails dev
) else (
    echo %YELLOW%Wails CLI not found. Starting with 'go run -tags dev .' (Press Ctrl+C to stop)%NC%
    call go run -tags dev .
)

set "EXIT_CODE=%errorlevel%"

if %EXIT_CODE% neq 0 (
    if %EXIT_CODE% neq 130 (
        if %EXIT_CODE% neq 3221225786 (
            :: 3221225786 is 0xC000013A (STATUS_CONTROL_C_EXIT) in Windows
            echo.
            echo %RED%Error: The application failed to start or exited with an error.^ %NC%
        ) else (
            echo.
            echo %GREEN%Application closed.^ %NC%
            set "EXIT_CODE=0"
        )
    ) else (
        echo.
        echo %GREEN%Application closed.^ %NC%
        set "EXIT_CODE=0"
    )
) else (
    echo.
    echo %GREEN%Application closed.^ %NC%
)
pause

:EXIT_WITH_ERROR_CODE
cmd /c exit %EXIT_CODE%
goto :EOF

:EXIT_WITH_ERROR
cmd /c exit 1
goto :EOF

endlocal
