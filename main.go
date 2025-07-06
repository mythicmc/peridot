package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
)

func main() {
	log.SetFlags(log.Lshortfile)

	name := filepath.Base(os.Args[0])
	args := os.Args[1:]
	help := false

	if len(os.Args) < 2 {
		log.Fatalln("no arguments provided")
	} else if slices.ContainsFunc(os.Args, func(arg string) bool {
		return arg == "--version" || arg == "-v"
	}) || os.Args[1] == "version" {
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
			fmt.Println("Show current state of the Minecraft servers.")
			fmt.Println("If a server is specified, show its state only.")
		} else {
			// FIXME
		}
		return
	case "apply":
		if help {
			fmt.Println("Usage: " + name + " apply [server]")
			fmt.Println("")
			fmt.Println("Apply current config to the Minecraft servers.")
			fmt.Println("If a server is specified, apply the config to that server only.")
		} else {
			// FIXME
		}
		return
	case "apply-live":
		if help {
			fmt.Println("Usage: " + name + " apply-live [server]")
			fmt.Println("")
			fmt.Println("Apply current config to the Minecraft servers without restart.")
			fmt.Println("If a server is specified, apply the config to that server only.")
		} else {
			// FIXME
			return
		}
	default:
		if !help {
			log.Fatalf("unknown command: %s\n", os.Args[1])
		}
	}

	if help {
		fmt.Println("Usage: " + name + " (command) [options]")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  help                      Show this help message")
		fmt.Println("  version                   Show version information")
		fmt.Println("  status, state [server]    Show current state of the Minecraft servers")
		fmt.Println("  apply [server]            Apply current config to the Minecraft servers")
		fmt.Println("  apply-live [server]       Apply current config to the Minecraft servers without restart")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  --version, -v            Show version information")
		fmt.Println("  --help, -h               Show this help message")
		return
	}
}
