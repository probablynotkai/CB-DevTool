@echo off

del ".\dist\cb.exe"

mkdir ".\dist"
go build -o dist/cb.exe

echo f | xcopy /f /y ".\cb.json" ".\dist\cb.json"