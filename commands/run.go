package commands

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"monot/config"

	"github.com/urfave/cli/v2"
)

type process struct {
	command string
	path    string
}

// RunCmd runs a command across all registered services
func RunCmd(ctx *cli.Context) error {
	processes := make(chan *process, 0)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	// Listen for quit signals
	go func() {
		sig := <-sigs
		fmt.Println("=======")
		fmt.Println(sig)
		done <- true
	}()

	// Open the cache file
	f, err := os.Open("./.monot/cache")
	if err != nil {
		log.Panic(err)
	}

	task := ctx.Args().First()

	// Decode the cache file into a config struct again
	decoder := gob.NewDecoder(f)
	var conf config.Config
	if err := decoder.Decode(&conf); err != nil {
		log.Panic(err)
	}

	// Iterate over all registered services, and find the associated
	// task, which matches the requested task. I.e. if we run a task
	// named `run-local`, it will scan each service for a task called
	// `run-local` and will find the commands for each of those.
	for path, service := range conf.RegisteredServices {
		t := service.Tasks[task]
		for _, command := range t.Commands {
			cmd := command
			p := path
			go func(cmd, p string) {
				processes <- &process{command: cmd, path: p}
			}(cmd, p)
		}
	}

	// Listens for processes, which is a channel. We do this so we can
	// run multiple services/commands at the same time
	for {
		select {
		case proc := <-processes:
			cmdParts := strings.Split(proc.command, " ")

			first := cmdParts[0]
			rest := cmdParts[1:]

			cmd := exec.Command(first, rest...)
			cmd.Dir = proc.path

			// This prints a bunch of semi-garbled output from each service ran currently...
			// TODO: improve the output to do something similar to `docker compose`.
			// [info][service-a] - some log output
			// [error][service-b] - an error happened in service b
			//
			stderr, _ := cmd.StderrPipe()
			stdout, _ := cmd.StdoutPipe()
			_ = cmd.Start()

			// TODO: there's a lot of duplication in the next two go routines...
			// I'm sure we could neaten this up!
			go func() {
				scanner := bufio.NewScanner(stderr)
				scanner.Split(bufio.ScanWords)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
				}
				_ = cmd.Wait()
			}()

			go func() {
				scanner := bufio.NewScanner(stdout)
				scanner.Split(bufio.ScanWords)
				for scanner.Scan() {
					m := scanner.Text()
					fmt.Println(m)
				}
				_ = cmd.Wait()
			}()
		case <-done:
			log.Println("finished, exiting")
			return nil
		}
	}

	return nil
}
