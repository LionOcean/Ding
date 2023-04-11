#! /bin/bash

echo -e "Start running the script..."
cd ../

echo -e "Start building the app for macos platform..."
wails build  -platform darwin/universal

echo -e "End running the script!"
