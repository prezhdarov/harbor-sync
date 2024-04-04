package harbor

import "flag"

var (
	hbSchema  = flag.String("schema", "https", "Use HTTP or HTTPS")
	hbTLS     = flag.Bool("verifyTLS", false, "Verify or trust TLS certificate")
	hbTimeout = flag.Int("timeout", 5, "Time in seconds to wait for harbor to reply")
)
