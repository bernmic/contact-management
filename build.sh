#!/bin/sh
cd frontend
npm install
ng build --prod
cd ..
rm -rf static
mkdir static
cp frontend/dist/ContactManager/* static/
go get ./...
go build .
