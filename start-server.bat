@echo off
echo =================================
echo  Wallpaper API Server
echo =================================
echo.
echo Starting Go server...
echo Server will be available at: http://localhost:8080
echo.
echo Place your images in:
echo   - images/nature/
echo   - images/culture/ 
echo   - images/digital/
echo.
echo Press Ctrl+C to stop the server
echo.

REM Initialize Go modules if not already done
if not exist go.sum (
    echo Initializing Go modules...
    go mod tidy
)

REM Start the server
go run main.go

pause