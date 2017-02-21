package main

import (
	"net/http"
	"goji.io"
	"goji.io/pat"
	"fmt"
	"encoding/json"
	"os/exec"
	"log"
	"bufio"
	"os"
)

func main() {
	mux := goji.NewMux()
	mux.Use(checkSecret)

	mux.HandleFunc(pat.Get("/tables"), tables)
	mux.HandleFunc(pat.Get("/query"), query)
	mux.HandleFunc(pat.Get("/table/:table"), table)

	bind := ""
	if os.Getenv("BIND") != "" {
		bind = os.Getenv("BIND")
	} else {
		bind = "localhost:8000"
	}

	http.ListenAndServe(bind, mux)
}

func tables(w http.ResponseWriter, r *http.Request) {
	lines := osqueryiDotTable()
	enc := json.NewEncoder(w)
	enc.Encode(lines)
}

func table(w http.ResponseWriter, r *http.Request) {
	sql := "select * from " + pat.Param(r, "table")
	result := osqueryiQuery(sql)
	fmt.Fprint(w, result)
}

func query(w http.ResponseWriter, r *http.Request) {
	result := osqueryiQuery(r.URL.Query().Get("q"))
	fmt.Fprint(w, result)
}

func checkSecret(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("secret") != os.Getenv("SECRET") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Go away!\n"))
			return
		}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func osqueryiQuery(sql string) string {

	bs, err := exec.Command("osqueryi", sql, "--json").Output()

	if err != nil {
		panic(err)
	}

	return string(bs)
}

func osqueryiDotTable() []string {

	cmd := exec.Command("osqueryi", ".table")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 2 {
			continue
		}
		lines = append(lines, line)
	}

	return lines

}