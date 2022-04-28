package instrumentation

import (
	"strings"
)

func InitSelfInstrumentation(name, license, apmEndpoint, telemetryEndpoint string) {
	if strings.ToLower(name) == apmInstrumentationName {
		apmSelfInstrumentation, err := NewAgentInstrumentationApm(
			license,
			apmEndpoint,
			telemetryEndpoint,
		)
		if err == nil {
			SelfInstrumentation = apmSelfInstrumentation
		}
	}
}
