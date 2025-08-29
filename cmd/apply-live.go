package cmd

import (
	"fmt"

	"github.com/mythicmc/peridot/deploy"
	"github.com/mythicmc/peridot/utils"
)

func HandleApplyLiveCommand(args []string) {
	_, configs, softwareUpdates, serverPropertiesUpdates, pluginUpdates := loadReposConfigUpdateState()
	previewUpdates(softwareUpdates, serverPropertiesUpdates, pluginUpdates)

	fmt.Print("Proceed to apply updates live? [y/N] ")
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		fmt.Println("Aborting live update.")
		return
	}

	for server, operation := range softwareUpdates {
		fmt.Println("Updating server software for: ", server)

		err := deploy.ApplySoftwareUpdate(operation)
		if err != nil {
			fmt.Println("Error updating server software for '"+server+"':", err)
		}
	}

	for server, operations := range serverPropertiesUpdates {
		fmt.Println("Updating server properties for: ", server)

		err := deploy.ApplyServerPropertiesUpdates(operations, configs[server])
		if err != nil {
			fmt.Println("Error updating server properties for '"+server+"':", err)
		}
	}

	for server, operations := range pluginUpdates {
		for _, operation := range operations {
			prevVersion := utils.PickNonEmptyString(operation.PrevVersion, "(missing)")
			newVersion := utils.PickNonEmptyString(operation.NewVersion, "(removed)")
			fmt.Println("Updating '"+server+"' plugin ", operation.PluginName,
				" ("+prevVersion+" -> "+newVersion+")")

			err := deploy.ApplyPluginUpdate(operation)
			if err != nil {
				fmt.Println("Error updating plugins for '"+server+"':", err)
			}
		}
	}
}
