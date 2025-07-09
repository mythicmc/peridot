package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/dop251/goja"
)

func ExecuteConfig(configFolder, configFile string) (Config, error) {
	vm := goja.New()

	vm.GlobalObject().Set("module", vm.NewObject())
	vm.GlobalObject().Set("require", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) != 1 {
			return vm.NewTypeError("require() expects a single argument")
		}
		modulePath := call.Argument(0).String()
		// Resolve module path
		if !filepath.IsAbs(modulePath) {
			modulePath = filepath.Join(configFolder, modulePath)
		}
		// Read the module file
		moduleData, err := os.ReadFile(modulePath)
		if err != nil {
			return vm.NewGoError(err)
		}
		// Execute the module code
		_, err = vm.RunString(string(moduleData))
		if err != nil {
			return vm.NewGoError(err)
		}
		// Return the exports of the module
		value := vm.GlobalObject().Get("module").ToObject(vm).Get("exports")
		return value
	})
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	_, err := vm.RunString(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	exportsJson, err := vm.RunString("JSON.stringify(module.exports)")
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal([]byte(exportsJson.String()), &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
