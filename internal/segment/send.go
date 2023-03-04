package segment

import (
	"fmt"

	"github.com/kubefirst/kubefirst/pkg"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/analytics-go"
)

func (c *SegmentClient) SendCountMetric(
	cliVersion string,
	cloudProvider string,
	clusterId string,
	clusterType string,
	domainName string,
	gitProvider string,
	kubefirstTeam string,
	metricName string,
) string {

	defer func(c *SegmentClient) {
		err := c.Client.Close()
		if err != nil {
			log.Info().Msgf("error sending identify to segment %s", err.Error())
		}
	}(c)

	if metricName == pkg.MetricInitStarted {
		err := c.Client.Enqueue(analytics.Identify{
			UserId: domainName,
			Type:   "identify",
		})
		if err != nil {
			return fmt.Sprintf("error sending identify to segment %s", err.Error())
		}
	}

	err := c.Client.Enqueue(analytics.Track{
		UserId: domainName,
		Event:  metricName,
		Properties: analytics.NewProperties().
			Set("cli_version", cliVersion).
			Set("cloud_provider", cloudProvider).
			Set("cluster_id", clusterId).
			Set("cluster_type", clusterType).
			Set("domain", domainName).
			Set("git_provider", gitProvider).
			Set("kubefirst_team", kubefirstTeam),
	})
	if err != nil {
		return fmt.Sprintf("error sending track to segment %s", err.Error())
	}

	return ""
}