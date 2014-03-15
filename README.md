# SJCPL Leaderboard

This code was used to maintain the leaderboard for the 2013 trivia night. The
overall workflow was

 1. Enter team data (name, table number, scores) into a csv file
 2. Tabulate and rank the teams; output the results as an html file
 3. Publish the html file on a web server running on the same laptop
 4. Display the html file using the target computer

Once this system is set up, there are no runtime dependencies on the internet.
However, there is a dependency on having the local wifi network working. The
only way around that is to generate the html file on the same computer as the
one which is displaying the leaderboard to the projector.

# Details

The csv file consists of a sequence of lines, each line having the format

    <table number>,<team name>, <score round 1>, <score round 2>, ...

A sample line is

    6,Deliberators,   8, 10, 4

This line describes the team "Deliberators" at table 6. They have scored 8 in
round 1, 10 points in round 2, and 4 in round 3. Currently, the csv file must
be in the current directory and named `team_data.txt`.

Execute the `update.sh` script to translate the csv file into the ranked html
document and copy it into the local web server's root. The generated files are
copied to the directory `/Library/WebServer/Documents/`, which is the default
root for the Apache install in OSX. (Enable "Web Sharing" in the Sharing
control panel). In my case, my machine name is `chance` (run `hostname` at the
terminal). The published page can then be accessed using the url
`http://chance.local`.

The `update.sh` script uses the `make_leaderboard.py` utility to generate the
html page. It requires the [wheezy][] template engine to be installed, for
better or worse. This could be easily swapped out for others. The file
`leaderboard_template.html` is the template for the html.

 [wheezy]: https://pypi.python.org/pypi/wheezy.template

Enjoy.


Don Brower  
don.brower@gmail.com
