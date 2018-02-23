package main

import (
	"strconv"
	"strings"
)

type command struct {
	name    string
	handler func(string, *Environment)
}

var cmdTable = []command{
	{"addteam", addteam},
	{"score", alterScore},
	{"add", addScore},
	{"save", saveEnv},
	{"load", loadEnv},
}

func takeTwo(s string) (a, b string) {
	pieces := strings.SplitN(s, " ", 2)
	if len(pieces) == 0 {
		return
	}
	a = strings.TrimSpace(pieces[0])
	if len(pieces) > 1 {
		b = strings.TrimSpace(pieces[1])
	}
	return
}

func takeThree(s string) (a, b, c string) {
	var bc string
	a, bc = takeTwo(s)
	b, c = takeTwo(bc)
	return
}

func cmdDispatch(s string, e *Environment) {
	cmd, arg := takeTwo(s)
	for i := range cmdTable {
		if cmdTable[i].name == cmd {
			cmdTable[i].handler(arg, e)
			return
		}
	}
}

func addteam(arg string, e *Environment) {
	number, name := takeTwo(arg)
	n, err := strconv.Atoi(number)
	if err != nil {
		return
	}
	t := e.findTeam(n)
	if t != nil {
		// alter team
		t.Name = name
		return
	}
	// new team
	t = NewTeam()
	t.Number = n
	t.Name = name
	e.Teams = append(e.Teams, t)
	e.RankTeams()
	e.Save(e.Filename)
}

func alterScore(arg string, e *Environment) {
	// set <team> <round> <score>
	number, round, score := takeThree(arg)
	genAlterScore(number, round, score, e)
	e.RankTeams()
	e.Save(e.Filename)
}

func genAlterScore(number, round, score string, e *Environment) {
	n, err := strconv.Atoi(number)
	if err != nil {
		return
	}
	r, err := strconv.Atoi(round)
	if err != nil {
		return
	}
	// round is 1 based. special case -1, though
	if (r < 1 || r > 10) && r != -1 {
		return
	}
	scr, err := strconv.Atoi(score)
	if err != nil {
		return
	}
	// clamp negative scores to -1
	if scr < 0 {
		scr = -1
	}
	t := e.findTeam(n)
	if t == nil {
		return
	}
	if r > 0 {
		t.Scores[r-1] = scr
		return
	}
	// update first empty score for team
	for i := range t.Scores {
		if t.Scores[i] < 0 {
			t.Scores[i] = scr
			break
		}
	}
}

func addScore(arg string, e *Environment) {
	// add <n> <score> [ "/" <n> <score> ... ]

	first := true
	for arg != "" {
		if !first {
			var slash string
			slash, arg = takeTwo(arg)
			if slash != "/" {
				break
			}
		}
		first = false
		var number, score string
		number, score, arg = takeThree(arg)
		genAlterScore(number, "-1", score, e)
	}
	e.RankTeams()
	e.Save(e.Filename)
}

func loadEnv(arg string, e *Environment) {
	e.Load(arg)
}

func saveEnv(arg string, e *Environment) {
	e.Save(arg)
}
