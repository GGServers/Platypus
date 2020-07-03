package pluginhandler

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"github.com/go-yaml/yaml"
)

var Plugins map[string]Plugin

// @TODO: This will cache for available hooks for a hook.
var PluginsForHooks map[string]HookFunc

type Plugin struct {
	Name            string `yaml:"name"`
	Version         string `yaml:"version"`
	Homepage        string `yaml:"homepage"`
	ImplementsHooks []HookFunc
}

type HookFunc struct {
	Type     string
	Name     string
	Func     string
	Contents string
}

// RegisterPlugins searches through plugins directory and registers plugins 
// with a valid plugin.yml file. Should only be run at startup.
// @TODO: Investigate reloading?
func RegisterPlugins() {
	Plugins = make(map[string]Plugin)

	plugins, err := ioutil.ReadDir("plugins")
	if err != nil {
		fmt.Println("Unable to read plugins directory")
		return
	}

	for _, plugin := range plugins {
		pluginName := plugin.Name()
		registeredPlugin, ok := Plugins[pluginName]
		// Already registered.
		if ok {
			continue
		}

		registeredPlugin = Plugin{}
		pluginInfo, err := ioutil.ReadFile("plugins/" + pluginName + "/plugin.yml")
		err = yaml.Unmarshal(pluginInfo, &registeredPlugin)
		// Malformed plugin.yml, @TODO log this somewhere for debugging
		if err != nil {
			continue
		}

		pluginContents, err := ioutil.ReadFile("plugins/" + pluginName + "/plugin.go")
		if err != nil {
			continue
		}
		pluginContent := string(pluginContents)

		hookRe := regexp.MustCompile(`(func )[A-Z][a-z]+\(.+\)\s?{(.|\n)*?\n}`)
		funcs := hookRe.FindAllStringSubmatch(pluginContent, -1)

		for range funcs {
			re := regexp.MustCompile(`((func )[A-Z][a-z]+)`)
			foundFuncName := re.FindStringSubmatch(pluginContent)
			realFuncName := strings.TrimLeft(foundFuncName[1], "func ")

			hookFunc := HookFunc{
				Type:     realFuncName,
				Func:     pluginName + "." + realFuncName,
				Contents: pluginContent,
			}

			registeredPlugin.ImplementsHooks = append(registeredPlugin.ImplementsHooks, hookFunc)
		}
		Plugins[pluginName] = registeredPlugin
		fmt.Println("Registered plugin " + pluginName)
	}

}

// Loop through registered plugins and execute any that match the hook.
func ExecuteHook(original string, hook string) string {
	data := original
	for _, plugin := range Plugins {
		for _, hookImplementor := range plugin.ImplementsHooks {
			if hookImplementor.Type == hook {
				// @TODO: Handle errors :(
				data, err := executePlugin(original, hookImplementor)
				if err != nil {
					fmt.Println(err.Error())
				}
				if data == "" {
					continue
				}
			}
		}
	}
	// If the plugins returned us nothing, return our original data.
	if data == "" {
		return original
	}
	return data
}

// Execute a plugins hook function, and return the string result.
func executePlugin(data string, hook HookFunc) (string, error) {
	result := data
	interpreter := interp.New(interp.Options{})
	interpreter.Use(stdlib.Symbols)
	_, err := interpreter.Eval(hook.Contents)
	if err != nil {
		return result, err
	}
	function, err := interpreter.Eval(hook.Func)
	if err != nil {
		return "", err
	}

	callable := function.Interface().(func(string) (string, error))
	result, err = callable(data)
	if err != nil {
		return "", err
	}

	return result, nil
}
