// Package autocomplete provides bash autocomplete for golang programs.
package autocomplete

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var acFlag = flag.Bool("c", false, "autocomplete parameters")

func init() {
	prog := os.Args[0]
	if prog == "" {
		return
	}
	prog = filepath.Base(prog)
	script := fmt.Sprintf(`complete -C "%s -c" %s`, prog, prog)
	file := fmt.Sprintf("%s/.config/bash_completion/%s", os.Getenv("HOME"), prog)
	addBashStatement(file, script)
}

func addBashStatement(file, script string) {
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		panic(err)
	}

	actext, _ := ioutil.ReadFile(file)
	if string(actext) != script {
		if err := ioutil.WriteFile(file, []byte(script), 0644); err != nil {
			panic(err)
		}
		fmt.Println("autocomplete installed. Run source ~/.bashrc now")
	}

	what := `for i in ~/.config/bash_completion/*; do source $i ; done`
	bashrc := os.Getenv("HOME") + "/.bashrc"
	buf, _ := ioutil.ReadFile(bashrc)
	lines := strings.Split(string(buf), "\n")
	found := false
	for _, ln := range lines {
		if ln == what {
			found = true
			break
		}
	}

	if !found {
		f, err := os.OpenFile(bashrc, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(f, "\n%s\n", what)
		f.Close()
		fmt.Println("common autocomplete installed. Run source ~/.bashrc now")
	}
}

// Handler is a function type for Handler function.
type Handler func() []string

var handlers = make(map[string]Handler)

// Handle adds custom handler for the parameter. Handler should return list of possible completions.
func Handle(name string, h Handler) {
	_, ok := handlers[name]
	if ok {
		panic("handler is already assigned: " + name)
	}
	handlers[name] = h
}

// HandleArgs prints autocompletion variants if -c flag is set and then program exits.
func HandleArgs() {
	if *acFlag {
		printCompletions()
		os.Exit(0)
		return
	}
}

func printCompletions() {
	lastarg := flag.Arg(flag.NArg() - 1)
	lastpar := flag.Arg(flag.NArg() - 2)

	for name, handler := range handlers {
		if lastarg != "-"+name {
			continue
		}
		list := handler()
		for _, s := range list {
			if strings.HasPrefix(s, lastpar) {
				fmt.Println(s)
			}
		}
		return
	}

	flag.VisitAll(func(f *flag.Flag) {
		if f.Name == "c" {
			return
		}
		if strings.HasPrefix("-"+f.Name, flag.Arg(1)) {
			fmt.Println("-" + f.Name)
		}
	})
}
