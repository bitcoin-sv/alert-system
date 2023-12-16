package base

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/alert-system/app/config"
	"github.com/bitcoin-sv/alert-system/app/models"
	"github.com/bitcoin-sv/alert-system/app/tester"
	"github.com/stretchr/testify/suite"
)

// TestSuite is for testing the entire package using real/mocked services
type TestSuite struct {
	Dependencies *config.Config // App config and services (dependencies)
	suite.Suite                 // Extends the suite.Suite package
}

// SetupSuite runs at the start of the suite
func (ts *TestSuite) SetupSuite() {

	// Set the env to test
	tester.SetupEnv(ts.T())

	// Load the configuration
	var err error
	ts.Dependencies, err = config.LoadConfig(context.Background(), models.BaseModels, true)
	ts.Require().NoError(err)
}

// TearDownSuite runs after the suite finishes
func (ts *TestSuite) TearDownSuite() {

	// Ensure all connections are closed
	if ts.Dependencies != nil {
		ts.Dependencies.CloseAll(context.Background())
	}

	tester.TeardownEnv(ts.T())
}

// SetupTest runs before each test
func (ts *TestSuite) SetupTest() {

	// Set the env to test
	tester.SetupEnv(ts.T())

	// Load the services
	var err error
	ts.Dependencies, err = config.LoadConfig(context.Background(), models.BaseModels, true)
	ts.Require().NoError(err)
}

// TearDownTest runs after each test
func (ts *TestSuite) TearDownTest() {
	if ts.Dependencies != nil {
		ts.Dependencies.CloseAll(context.Background())
	}

	tester.TeardownEnv(ts.T())
}

// TestTestSuiteApp kick-starts all suite tests
func TestTestSuiteApp(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
