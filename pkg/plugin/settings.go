package plugin

import (
	"encoding/json"

	"github.com/dalelane/grafana-kafka-datasource/pkg/kafka"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Datasource struct {
	kafkaconfig kafka.SimpleClientConfig
}

func getDatasourceSettings(s backend.DataSourceInstanceSettings) (*Datasource, error) {
	Logger.Debug("settings :: getDatasourceSettings", "JSONData", string(s.JSONData))

	cfg := &kafka.SimpleClientConfig{}
	if err := json.Unmarshal(s.JSONData, cfg); err != nil {
		return nil, err
	}

	cfg.ScramPassword = s.DecryptedSecureJSONData["password"]

	return &Datasource{
		kafkaconfig: *cfg,
	}, nil
}
