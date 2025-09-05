package observability

import (
	"log"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrwriter"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// SetupNewRelic configures the application and configures the singleton logger to forward logs to New Relic
func SetupNewRelic(appName, newrelicKey string) (*newrelic.Application, error) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(newrelicKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		return app, err
	}
	if err := app.WaitForConnection(5 * time.Second); err != nil {
		return app, err
	}

	writer := nrwriter.New(os.Stdout, app)

	m.Lock()
	defer m.Unlock()
	logger = log.New(writer, "", log.Default().Flags())

	return app, nil
}
