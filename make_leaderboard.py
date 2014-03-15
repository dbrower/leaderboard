#!/usr/bin/env python

import csv

def int_maybe(x):
    try:
        return int(x)
    except:
        return 0

class Cycle:
    def __init__(self):
        self.is_odd = False
    def cycle(self):
        self.is_odd = not self.is_odd
        return self.is_odd

class Team:
    def __init__(self, name, tablenr, scores):
        self.name = name
        self.tablenr = tablenr
        self.scores = scores
        self.total = sum(map(int_maybe, scores))
        self.rank = None
    def __repr__(self):
        return "#{0.rank} {0.name} (@{0.tablenr}) {0.scores} = {0.total}".format(self)

def compare_teams(x, y):
    # primary sort is total score, decreasing
    r = y.total - x.total
    if r == 0:
        # secondary sort is table number, increasing
        r = x.tablenr - y.tablenr
        # secondary sort is team name alphabetical
        #r = cmp(x.name, y.name)
    return r

def assign_ranks(teams):
    rank = 1
    previous_total = None
    current_index = 0
    for team in teams:
        if previous_total != team.total:
            rank = current_index + 1
            previous_total = team.total
        team.rank = rank
        current_index += 1

# read data file
teams = []
with open('team_data.txt') as f:
    # csv.reader seems to only support ascii characters
    reader = csv.reader(f, delimiter=',')
    for row in reader:
        teams.append(Team(row[1], int(row[0]), row[2:]))

teams.sort(compare_teams)
# only assign ranks if points have been scored
if teams[0].total != 0:
    assign_ranks(teams)

# now generate html file
from wheezy.template.engine import Engine
from wheezy.template.ext.core import CoreExtension
from wheezy.template.loader import FileLoader

searchpath = ['.']
engine = Engine(
    loader=FileLoader(searchpath),
    extensions=[CoreExtension()]
)
template = engine.get_template('leaderboard_template.html')

# cycle for the table zebra striping
c = Cycle()
with open('leaderboard.html', 'w') as f:
    f.write(template.render({"teams" : teams, "cycle" : c.cycle}))
