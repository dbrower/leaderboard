package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

// An Environment contains all the data needed for our state
type Environment struct {
	m        sync.RWMutex // protects everything below
	Teams    []*Team
	History  []string
	Filename string
}

// A Team represents a single team.
// The Scores have a funny meaning...a negative score means "not entered",
// whereas a zero score means "received 0 points that round"
type Team struct {
	Number int // Team ID
	Name   string
	Scores [10]int
	Total  int
	Group  string // top, middle, or bottom
	Rank   int    // rank of this team
}

func (e *Environment) findTeam(n int) *Team {
	for i := range e.Teams {
		if e.Teams[i].Number == n {
			return e.Teams[i]
		}
	}
	return nil
}

type byTotal []*Team

// The comparison in Less is reversed since we want the LARGEST total first
func (d byTotal) Len() int           { return len(d) }
func (d byTotal) Less(i, j int) bool { return d[i].Total > d[j].Total }
func (d byTotal) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func (e *Environment) RankTeams() {
	for i := range e.Teams {
		e.Teams[i].updatetotal()
	}
	sort.Stable(byTotal(e.Teams))

	level := 1
	n := len(e.Teams)
	for i := range e.Teams {
		var rank = level
		if i > 0 && e.Teams[i-1].Total == e.Teams[i].Total {
			rank = e.Teams[i-1].Rank
		}
		e.Teams[i].Rank = rank
		level++
		if rank <= n/3 {
			e.Teams[i].Group = "top"
		} else if rank <= 2*n/3 {
			e.Teams[i].Group = "middle"
		} else {
			e.Teams[i].Group = "bottom"
		}
	}
}

func NewTeam() *Team {
	var t Team
	for i := range t.Scores {
		t.Scores[i] = -1
	}
	return &t
}

func (t *Team) updatetotal() {
	var n int
	for _, x := range t.Scores {
		if x > 0 {
			n += x
		}
	}
	t.Total = n
}

func (e *Environment) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	for i := range e.Teams {
		err := writeTeam(f, e.Teams[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func writeTeam(w io.Writer, team *Team) error {
	_, err := fmt.Fprintf(w, "%d\t%s", team.Number, team.Name)
	if err != nil {
		return err
	}

	for _, score := range team.Scores {
		_, err = fmt.Fprintf(w, "\t%d", score)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(w, "\n")
	return err
}

func (e *Environment) Load(filename string) error {
	var result []*Team
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	lines := strings.Split(string(data), "\n")
	for i := range lines {
		t := readTeam(lines[i])
		if t != nil {
			result = append(result, t)
		}
	}
	e.Teams = result
	e.RankTeams()
	return nil
}

func readTeam(line string) *Team {
	t := NewTeam()
	fields := strings.Split(line, "\t")
	for i := range fields {
		switch i {
		case 0:
			var err error
			t.Number, err = strconv.Atoi(fields[i])
			if err != nil {
				return nil
			}
		case 1:
			t.Name = fields[i]
		case 2, 3, 4, 5, 6, 7, 8, 9, 10, 11:
			var err error
			n, err := strconv.Atoi(fields[i])
			if err == nil {
				t.Scores[i-2] = n
			}
		}
	}
	return t
}

var (
	env          Environment
	viewT        = template.Must(template.ParseFiles("admin_template.html"))
	leaderboardT = template.Must(template.ParseFiles("leaderboard_template.html"))
)

func main() {
	env.Filename = "team_data.tsv"
	env.Load(env.Filename)

	publicmux := mux.NewRouter()
	publicmux.HandleFunc("/", renderLeaderboardPage).Methods("GET")
	publicmux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("."))))
	go func() {
		err := http.ListenAndServe(":80", &logHandler{name: "public", h: publicmux})
		log.Println(err)
	}()

	privateMux := mux.NewRouter()
	privateMux.HandleFunc("/", renderAdminPage).Methods("GET")
	privateMux.HandleFunc("/", updateAdminPage).Methods("POST")
	http.ListenAndServe("localhost:8000", &logHandler{name: "admin", h: privateMux})
}

type logHandler struct {
	name string
	h    http.Handler
}

func (lh *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(lh.name, r.Method, r.URL)
	lh.h.ServeHTTP(w, r)
}

func renderAdminPage(w http.ResponseWriter, r *http.Request) {
	env.m.RLock()
	defer env.m.RUnlock()
	viewT.Execute(w, env)
}

func updateAdminPage(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("cmd")
	env.m.Lock()
	defer env.m.Unlock()
	if cmd != "" {
		env.History = append(env.History, cmd)
		cmdDispatch(cmd, &env)
	}
	viewT.Execute(w, env)
}

func renderLeaderboardPage(w http.ResponseWriter, r *http.Request) {
	env.m.RLock()
	defer env.m.RUnlock()
	leaderboardT.Execute(w, env)
}
