#!/bin/bash

./make_leaderboard.py
# for system apache
#-target=/Library/WebServer/Documents/
# for brew nginx
target=/usr/local/var/www
files="leaderboard.css sjcpl_logo_bw.png"
sudo cp leaderboard.html $target/index.html
sudo cp $files $target
