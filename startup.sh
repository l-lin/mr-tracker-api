#!/bin/bash

export PORT=3000
export GOOGLE_CLIENT_ID=""
export GOOGLE_CLIENT_SECRET=""
export GOOGLE_REDIRECT_URI="http://localhost:3000/oauth2callback"
export DATABASE_URL="postgres://postgres@localhost:5432/mr?sslmode=disable"
export ADMIN_ID=""

if [ -z ${GOOGLE_CLIENT_ID} ]; then
    echo "[x] You need to provide the variable GOOGLE_CLIENT_ID from Google developer console: https://console.developers.google.com/project"
    exit -1
fi

if [ -z ${GOOGLE_CLIENT_SECRET} ]; then
    echo "[x] You need to provide the variable GOOGLE_CLIENT_SECRET from Google developer console: https://console.developers.google.com/project"
    exit -1
fi

gin run

