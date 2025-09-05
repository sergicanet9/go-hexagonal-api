package observability

import (
	"fmt"
	"io"
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
	// logger.SetOutput(writer) // TODO Put back
	logger.SetOutput(&tempDebugWriter{target: writer}) // TODO remove

	return app, nil
}

// TODO remove
type tempDebugWriter struct {
	target io.Writer
}

func (w *tempDebugWriter) Write(p []byte) (int, error) {
	fmt.Printf("DEBUG log -> %s", string(p))
	fmt.Printf("DEBUG writer type -> %T\n", w.target)
	return w.target.Write(p)
}
