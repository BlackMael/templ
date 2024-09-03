@echo off

@echo "Bump Version"

go run ./get-version > .version

@rem build for host os (Windows)

@echo "Build Windows..."

go build -o templ.exe .\cmd\templ

@rem this is a quick example of building go app to target linux

@echo "Build Linux..."

set GOOS=linux
set GOARCH=amd64

go build -o templ .\cmd\templ

@rem this is a quick hack to copy the exe's to the bananas repo

@echo "Copy executables..."

xcopy /y templ ..\tradesomething-htmx-demo\src\
xcopy /y templ.exe ..\tradesomething-htmx-demo\src\
xcopy /y templ.exe ..\tradesomething-htmx-demo\alpha\
