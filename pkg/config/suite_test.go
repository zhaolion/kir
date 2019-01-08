package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type stubConfig struct {
	AppName  string         `json:"app_name" yaml:"app_name" toml:"app_name"`
	AppEnv   string         `json:"app_env" yaml:"app_env" toml:"app_env" env:"APP_ENV"`
	Debug    bool           `json:"debug" yaml:"debug" toml:"debug" env:"APP_DEBUG"`
	Database databaseConfig `json:"database" yaml:"database" toml:"database"`
}

type databaseConfig struct {
	DSN string `json:"dsn" yaml:"dsn" toml:"dsn" env:"DB_DSN"`
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type Suite struct {
	suite.Suite

	backupEnvs []string

	JSONFile string
	EnvFile  string
	YAMLFile string
	TOMLFile string
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *Suite) SetupSuite() {
	suite.JSONFile = "test-json.json"
	initConfig(suite.JSONFile, jsonExample)

	suite.EnvFile = "test-env.json"
	initConfig(suite.EnvFile, jsonExample)

	suite.YAMLFile = "test-yaml.yml"
	initConfig(suite.YAMLFile, yamlExample)

	suite.TOMLFile = "test-toml.toml"
	initConfig(suite.TOMLFile, tomlExample)
}

func (suite *Suite) TearDownSuite() {
	os.Remove(suite.JSONFile)
	os.Remove(suite.EnvFile)
	os.Remove(suite.YAMLFile)
	os.Remove(suite.TOMLFile)
}

// The SetupTest method will be run before every test in the suite.
func (suite *Suite) SetupTest() {
	// protected env keys
	suite.backupEnvs = os.Environ()
	os.Clearenv()
}

// The TearDownTest method will be run after every test in the suite.
func (suite *Suite) TearDownTest() {
	for _, pair := range suite.backupEnvs {
		kvs := strings.Split(pair, "=")

		if len(kvs) > 2 {
			os.Setenv(kvs[0], kvs[1])
		}
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func initConfig(file, config string) error {
	Reset()
	outputFile, outputError := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		return errors.New("an error occurred with file opening or creation")
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	outputWriter.WriteString(config)
	return outputWriter.Flush()
}

var jsonExample = `
{
	"app_name": "test-app",
	"app_env": "test",
	"debug": true,
	"database": {
	  "dsn": "root@tcp(localhost:3306)/test"
	}
  }
`

var yamlExample = `
app_name: test-app
app_env: test
debug: true
database:
  dsn: "root@tcp(localhost:3306)/test"
`

var tomlExample = `
app_name = "test-app"
app_env = "test"
debug = true

[database]
dsn = "root@tcp(localhost:3306)/test"
`
