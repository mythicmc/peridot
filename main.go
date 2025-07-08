package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/mythicmc/peridot/config"
	"github.com/mythicmc/peridot/repos"
)

func main() {
	repos, err := repos.LoadRepositories()
	if err != nil {
		println(err.Error())
	}
	_, err = config.LoadConfigs(repos)
	if err != nil {
		println(err.Error())
	}

	name := filepath.Base(os.Args[0])
	args := os.Args[1:]
	help := false

	if len(os.Args) < 2 {
		log.SetFlags(log.Lshortfile)
		log.Fatalln("no arguments provided")
	} else if slices.ContainsFunc(os.Args, func(arg string) bool {
		return arg == "--version" || arg == "-v"
	}) || os.Args[1] == "version" {
		log.SetFlags(log.Lshortfile)
		log.Println("version unknown")
		return
	} else if slices.ContainsFunc(os.Args, func(arg string) bool {
		return arg == "--help" || arg == "-h"
	}) || os.Args[1] == "help" {
		help = true
		if os.Args[1] == "help" && len(args) > 1 {
			args = args[1:]
		}
	}

	switch args[0] {
	case "status", "state":
		if help {
			fmt.Println("Usage: " + name + " status [server]")
			fmt.Println("    OR " + name + " state [server]")
			fmt.Println("")
			fmt.Println("Show current state of configured Minecraft servers.")
			fmt.Println("This command compares the current state of the server files on disk with the")
			fmt.Println("desired state defined in Peridot's configuration.")
			fmt.Println("If a server is specified, it shows only its state.")
		} else {
			// FIXME
		}
		return
	case "apply":
		if help {
			fmt.Println("Usage: " + name + " apply [server]")
			fmt.Println("")
			fmt.Println("Apply current config to the Minecraft server files on disk.")
			fmt.Println("This command will restart the server(s) if possible.")
			fmt.Println("If a server is specified, apply the config to that server only.")
		} else {
			// FIXME
		}
		return
	case "apply-live":
		if help {
			fmt.Println("Usage: " + name + " apply-live [server]")
			fmt.Println("")
			fmt.Println("Apply current config to the Minecraft server files on disk, without restarts.")
			fmt.Println("If a server is specified, apply the config to that server only.")
		} else {
			// FIXME
			return
		}
	default:
		if !help {
			log.SetFlags(log.Lshortfile)
			log.Fatalf("unknown command: %s\n", os.Args[1])
		}
	}

	if help {
		fmt.Println("Usage: " + name + " (command) [options]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  help                      Show this help message")
		fmt.Println("  version                   Show version information")
		fmt.Println("  status, state [server]    Show current state of configured Minecraft servers")
		fmt.Println("  apply [server]            Apply current config to Minecraft server files")
		fmt.Println("  apply-live [server]       Apply current config to Minecraft server files,")
		fmt.Println("                            without restarts")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  --version, -v            Show version information")
		fmt.Println("  --help, -h               Show this help message")
		return
	}
}
