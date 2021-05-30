#!/bin/bash

function create_binaries() {
    version=$1
    plataforms=("windows" "darwin" "linux")

    for goos in "${plataforms[@]}"; do
        GOOS=${goos} GOARCH=amd64 go build -ldflags="-X 'main.Version=${version}'" -o changelogger-${goos} .
    done
    mv changelogger-windows changelogger-windows.exe
}

create_binaries $1
