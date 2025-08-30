package cmd

import (
	"fmt"
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

	interactivelyApplyUpdates(configs, softwareUpdates, serverPropertiesUpdates, pluginUpdates)
}
