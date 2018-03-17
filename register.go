package metricsink

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	chassisMetrics "github.com/ServiceComb/go-chassis/metrics"
	"github.com/ServiceComb/go-chassis/third_party/forked/afex/hystrix-go/hystrix/metric_collector"
	"github.com/rcrowley/go-metrics"
	"time"
)

func init() {
	chassisMetrics.InstallReporter("CSE Monitoring", reportMetricsToCSEDashboard)
}

//reportMetricsToCSEDashboard use go-metrics to send metrics to cse dashboard
func reportMetricsToCSEDashboard(r metrics.Registry) error {
	metricCollector.Registry.Register(NewCseCollector)

	monitorServerURL, err := chassisMetrics.GetMonitorEndpoint()
	if err != nil {
		lager.Logger.Warn("Get Monitoring URL failed, CSE monitoring function disabled", err)
		return nil
	}

	tlsConfig, tlsError := chassisMetrics.GetTLSForClient(monitorServerURL)
	if tlsError != nil {
		lager.Logger.Errorf(tlsError, "Get %s.%s TLS config failed.", monitorServerURL, common.Consumer)
		return tlsError
	}

	InitializeCseCollector(&CseCollectorConfig{
		CseMonitorAddr: monitorServerURL,
		Header:         chassisMetrics.GetAuthHeaders(),
		TimeInterval:   time.Second * 2,
		TLSConfig:      tlsConfig,
	}, r, config.GlobalDefinition.AppID, config.SelfVersion, config.SelfServiceName,
		config.MicroserviceDefinition.ServiceDescription.Environment)
	lager.Logger.Infof("Started sending metric Data to Monitor Server : %s", monitorServerURL)
	return nil
}
