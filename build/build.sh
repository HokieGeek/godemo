#!/bin/sh

cd ../cmd/nx
GOOS=linux GOARCH=amd64 go build -o nx-nix .
GOOS=darwin GOARCH=amd64 go build -o nx-mac .
GOOS=windows GOARCH=amd64 go build -o nx-win.exe .
tar -cf nexus-cli.tar nx-nix nx-mac nx-win.exe
rm nx-nix nx-mac nx-win.exe
