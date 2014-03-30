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
	ORG_URI = "https://api.travis-ci.org/"
	PRO_URI = "https://api.travis-ci.com/"
)

type EndpointConfig struct {
	AccessToken string "access_token"
}

type LastCheckConfig struct {
	Etag    string "etag"
	Version string "version"
	At      int    "at"
}

type Config struct {
	LastCheck         LastCheckConfig            "last_check"
	CheckedCompletion bool                       "checked_completion"
	CompletionVersion string                     "completion_version"
	Endpoints         map[string]*EndpointConfig "endpoints"
	DefaultEndpoint   string                     "default_endpoint,omitempty"
}

func (config Config) Copy() Config {
	return config
}

var config = Config{}
var original_config = Config{}
var access_token string
var isDebug bool
var api_endpoint = ORG_URI
var explicit_api_endpoint = false

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
			apiExecute(RunToken)
		case "version":
			RunVersion()
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
	authenticate()
	fmt.Printf("Your access token is %s\n", ColotImportant(access_token))
}

func authenticate() {
	if access_token == "" {
		fmt.Fprintf(os.Stderr,
			ColorError("not logged in, please run %s\n"),
			command(fmt.Sprintf("login%s", endpointOption())))
		os.Exit(1)
	}
}

func isOrg() bool {
	return api_endpoint == ORG_URI
}

func isPro() bool {
	return api_endpoint == PRO_URI
}

func isDetectedEndpoint() bool {
	return api_endpoint == detectedEndpoint()
}

func detectedEndpoint() string {
	endpoint := defaultEndpoint()
	if endpoint == "" {
		endpoint = ORG_URI
	}
	return endpoint
}

func endpointOption() string {
	if isOrg() && isDetectedEndpoint() {
		return ""
	}
	if isOrg() {
		return " --org"
	}
	if isPro() {
		return " --pro"
	}

	// TODO: add option for config['enterprise']
	return fmt.Sprintf(" -e \"%s\"", api_endpoint)
}

func fetchToken() string {
	token := os.Getenv("TRAVIS_TOKEN")
	if token == "" {
		return endpoint_config().AccessToken
	}
	return token
}

func endpoint_config() *EndpointConfig {
	if config.Endpoints == nil {
		config.Endpoints = make(map[string]*EndpointConfig)
	}

	if endpoint, found := config.Endpoints[api_endpoint]; found {
		return endpoint
	}

	config.Endpoints[api_endpoint] = &EndpointConfig{}
	return config.Endpoints[api_endpoint]
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

		original_config = config.Copy()
	}
}

func loadFile(name string) ([]byte, error) {
	config_path := configPath(name)
	debug(fmt.Sprintf("Loading \"%s\"", config_path))
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

func execute(cmd func()) {
	// setup_trap
	// check_arity(method(:run), *arguments)
	LoadConfig()
	// check_version
	// check_completion
	// setup
	cmd()
	// clear_error
	// store_config
	// rescue Travis::Client::NotLoggedIn => e
	// raise(e) if explode?
	// error "#{e.message} - try running #{command("login#{endpoint_option}")}"
	// rescue StandardError => e
	// raise(e) if explode?
	// message = e.message
	// message += color("\nfor a full error report, run #{command("report#{endpoint_option}")}", :error) if interactive?
	// store_error(e)
	// error(message)

}

func apiSetup() {
	// setup_enterprise
	endpoint := defaultEndpoint()
	if endpoint != "" && !explicit_api_endpoint {
		setApiEndpoint(endpoint)
	}
	access_token = fetchToken()
	if endpoint_config().AccessToken == "" {
		endpoint_config().AccessToken = access_token
	}
	// endpoint_config['insecure']       = insecure unless insecure.nil?
	// self.insecure                     = endpoint_config['insecure']
	// session.ssl                       = { :verify => false } if insecure?
	// authenticate if pro? or enterprise?

	data, err := yaml.Marshal(&config)
	if err != nil {
		// log.Fatalf("error: %v", err)
	}
	fmt.Print(string(data))
}

func apiExecute(cmd func()) {
	execute(func() {
		apiSetup()
		cmd()
	})
}

func setApiEndpoint(uri string) {
	explicit_api_endpoint = true
	api_endpoint = uri
}

func defaultEndpoint() string {
	endpoint := os.Getenv("TRAVIS_ENDPOINT")
	if endpoint == "" {
		endpoint = config.DefaultEndpoint
	}
	return endpoint
}
