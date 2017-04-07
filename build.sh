#!/bin/bash
GOARCH=386 GOOS=windows CGO_ENABLED=0 go build -i
mv YourRansom.exe bin/YourRansom-win32.exe
GOARCH=386 GOOS=linux CGO_ENABLED=0 go build -i
mv YourRansom bin/YourRansom-linux32
GOARCH=386 GOOS=darwin CGO_ENABLED=0 go build -i
mv YourRansom bin/YourRansom-darwin32
