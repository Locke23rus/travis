package main

import (
	"fmt"
	"gopkg.in/v1/yaml"
	"io/ioutil"
	"os"
	"path"
)

var commands = map[string]string{
	"help":    "helps you out when in dire need of information",
	"token":   "outputs the secret API token",
	"version": "outputs the client version",
}

const (
	version = "0.0.1dev"
)

type ConfigEndpoint struct {
	AccessToken string "access_token"
}

type ConfigLastCheck struct {
	Etag    string "etag"
	Version string "version"
	At      int    "at"
}

type Config struct {
	LastCheck         ConfigLastCheck           "last_check"
	CheckedCompletion bool                      "checked_completion"
	CompletionVersion string                    "completion_version"
	Endpoints         map[string]ConfigEndpoint "endpoints"
}

func (config Config) Copy() Config {
	return config
}

var config = Config{}
var original_config = Config{}
var access_token string
var isDebug bool

func main() {
	if len(os.Args) > 1 {
		isDebug = containts(os.Args, "--debug")
		switch os.Args[1] {
		default:
			fmt.Printf("unknown command %s\n", os.Args[1])
			os.Exit(1)
		case "help":
			RunHelp()
		case "token":
			RunToken()
		case "version":
			RunVersion()
		case "load_config":
			LoadConfig()
		}
	} else {
		RunHelp()
	}
}

func RunHelp() {
	fmt.Println("Usage: travis COMMAND ...\n\nAvailable commands:\n")
	for command, description := range commands {
		fmt.Printf("\t%s\t%s\n",
			ColorCommand(command),
			ColorInfo(description))
	}
	fmt.Printf("\nrun `%s help COMMAND` for more infos\n", os.Args[0])
}

func RunVersion() {
	fmt.Println(version)
}

func RunToken() {

	if access_token == "" {
		fmt.Fprintf(os.Stderr,
			ColorError("not logged in, please run %s\n"),
			command("login"))
		os.Exit(1)
	} else {
		fmt.Printf("Your access token is %s\n", ColotImportant(access_token))
	}
}

func authenticate() {

}

func fetchToken() string {
	token := os.Getenv("TRAVIS_TOKEN")
	return token
}

func command(name string) string {
	return ColorCommand(fmt.Sprintf("%s %s", os.Args[0], name))
}

func debug(line string) {
	if !isDebug {
		return
	}
	fmt.Fprintln(os.Stderr, ColorDebug(fmt.Sprintf("** %s", line)))
}

func warn(line string) {
	fmt.Fprintln(os.Stderr, ColorError(line))
}

func ColorCommand(s string) string {
	return fmt.Sprintf("\x1b[1m%s\x1b[0m", s)
}

func ColorInfo(s string) string {
	return fmt.Sprintf("\x1b[33m%s\x1b[0m", s)
}

func ColotImportant(s string) string {
	return fmt.Sprintf("\x1b[1;4m%s\x1b[0m", s)
}

func ColorError(s string) string {
	return fmt.Sprintf("\x1b[31m%s\x1b[0m", s)
}

func ColorDebug(s string) string {
	return fmt.Sprintf("\x1b[35m%s\x1b[0m", s)
}
func LoadConfig() {
	data, err := loadFile("config.yml")
	if err != nil {
		if !os.IsNotExist(err) {
			exitIfErr(err)
		}
	} else {
		err = yaml.Unmarshal(data, &config)
		exitIfErr(err)
		fmt.Println(config)
		original_config = config.Copy()
	}
}

func loadFile(name string) ([]byte, error) {
	config_path := configPath(name)
	debug(fmt.Sprintf("Loading ‘%s‘", config_path))
	return ioutil.ReadFile(config_path)
}

func containts(slice []string, element string) bool {
	for _, e := range slice {
		if element == e {
			return true
		}
	}
	return false
}

func configPath(name string) string {
	config_path := os.Getenv("TRAVIS_CONFIG_PATH")
	if config_path == "" {
		config_path = path.Join(os.Getenv("HOME"), ".travis")
	}

	err := os.MkdirAll(config_path, 0700)
	exitIfErr(err)

	return path.Join(config_path, "config.yml")
}

func exitIfErr(err error) {
	if err != nil {
		warn(fmt.Sprint(err))
		os.Exit(1)
	}
}
