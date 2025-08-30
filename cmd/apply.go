package cmd

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mythicmc/peridot/utils"
)

func HandleApplyCommand(args []string) {
	_, configs, softwareUpdates, serverPropertiesUpdates, pluginUpdates := loadReposConfigUpdateState()
	previewUpdates(softwareUpdates, serverPropertiesUpdates, pluginUpdates)

	fmt.Print("Proceed to apply updates? [y/N] ")
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		fmt.Println("Aborting update.")
		return
	}

	affectedServers := make(map[string]struct{})
	for server := range softwareUpdates {
		affectedServers[server] = struct{}{}
	}
	for server := range serverPropertiesUpdates {
		affectedServers[server] = struct{}{}
	}
	for server := range pluginUpdates {
		affectedServers[server] = struct{}{}
	}

	// Send stop signals to all servers
	fmt.Println("Stopping affected servers via Octyne...")
	for server := range affectedServers {
		err := utils.OctyneTerminateServer(server)
		if err != nil {
			log.Fatalf("Error stopping server %s: %v\n", server, err)
		}
	}

	// Wait for servers to stop
	var wg sync.WaitGroup
	for server := range affectedServers {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			for {
				<-time.After(2 * time.Second)
				status, err := utils.OctyneGetServerStatus(s)
				if err != nil {
					log.Fatalf("Error getting status for server %s: %v\n", s, err)
				} else if status == utils.OctyneServerStatusStopped ||
					status == utils.OctyneServerStatusCrashed {
					break
				}
			}
		}(server)
	}
	select {
	case <-time.After(60 * time.Second):
		log.Fatalln("Failed to stop all servers after 60 seconds! Exiting...")
	case <-utils.WaitGroupDoneAsChannel(&wg):
	}

	interactivelyApplyUpdates(configs, softwareUpdates, serverPropertiesUpdates, pluginUpdates)

	// Send start signals to all servers
	fmt.Println("Starting affected servers via Octyne...")
	for server := range affectedServers {
		err := utils.OctyneStartServer(server)
		if err != nil {
			log.Printf("Error starting server %s: %v\n", server, err)
		}
	}

	fmt.Println("All updates have been applied successfully!")
}
