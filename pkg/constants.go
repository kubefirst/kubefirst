package pkg

const (
	ArgoCDLocalBaseURL  = "https://localhost:8080/api/v1"
	JSONContentType     = "application/json"
	SoftServerURI       = "ssh://127.0.0.1:8022/config"
	GitHubOAuthClientId = "Iv1.e783285e6ad438ac" // todo: update it
)

// SegmentIO constants
// SegmentIOWriteKey The write key is the unique identifier for a source that tells Segment which source data comes
// from, to which workspace the data belongs, and which destinations should receive the data.
const (
	SegmentIOWriteKey                 = "0gAYkX5RV3vt7s4pqCOOsDb6WHPLT30M"
	MetricInitStarted                 = "kubefirst.init.started"
	MetricInitCompleted               = "kubefirst.init.completed"
	MetricMgmtClusterInstallStarted   = "kubefirst.mgmt_cluster_install.started"
	MetricMgmtClusterInstallCompleted = "kubefirst.mgmt_cluster_install.completed"
)
