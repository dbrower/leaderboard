# SJCPL Leaderboard

This code runs the leaderboard for trivia night. Once the system is set up
there are no runtime dependencies on the internet. The code opens a web server
on port 80 to display the leaderboard, and runs a second web server on port
8000 to provide an admin console. The admin console allows for adding teams,
adding scores, and modifying scores. All changes to the scores are saved to
the file `team_data.tsv`.

The displayed score page has javascript to progressively revealed the scores.
At first no scores will show. Pressing `space` will reveal the bottom
third of teams. Press `space` again to show the bottom two thirds, and pressing
`space` a third time will display the rest.

This code originally was a python script for the 2013 trivia night. It was
rewritten in Go for the 2017 trivia night.

# Details

The program is in Go; compile it by running `go build`. Then run the program
by running `sudo ./leaderboard`. The `sudo` is required since root privildges
are needed for the program to listen on port 80, the usual HTTP server port.

View the admin page by visiting `localhost:8000` in a browser window.

You can view the score page by visiting the web page on the computer the program
is running on. In my case, my machine name is `chance` (run `hostname` at the
terminal). The published page can then be accessed using the url
`http://chance.local`.

The teams are grouped into the bottom three, the top three, and everything in
the middle. When the page is opened, no teams are displayed. By hitting the
spacebar, the teams are revealed. First the bottom group is shown, then the
middle, and finally, on the third press of the spacebar, the top teams are
revealed. Hitting shift-spacebar will hide the groups in the reverse order.

The program keeps the current scores saved into this file. In case of an emergency
the file could be edited by hand and reloaded with the `load team_data.txt`
command in the admin console.

The save file consists of a sequence of lines, each line having the format

    <table number> \t <team name> \t <score round 1> \t <score round 2> \t ...

where `\t` is a tab character.
A sample line is

    6\tDeliberators\t8\t10\t4

This line describes the team "Deliberators" at table 6. They have scored 8 in
round 1, 10 points in round 2, and 4 in round 3. Currently, the csv file must
be in the current directory and named `team_data.txt`.

The file `leaderboard_template.html` is the template for the score display page.
The file `admin_template.html` is the template for the admin page.

Enjoy.
