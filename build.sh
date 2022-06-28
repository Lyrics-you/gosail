#!/bin/bash

go build .
cd gocy && go build .
cd ../gobars && go build .
cd ..

mkdir sailcybars
mv ./gosail ./sailcybars
mv ./gocy/gocy ./sailcybars
mv ./gobars/gobars ./sailcybars
