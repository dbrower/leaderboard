#!/bin/bash

./make_leaderboard.py
target=/Library/WebServer/Documents/
files="leaderboard.css sjcpl_logo.jpg"
sudo cp leaderboard.html $target/index.html
sudo cp $files $target
