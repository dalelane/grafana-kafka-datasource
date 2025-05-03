package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/dalelane/grafana-kafka-datasource/pkg/kafka"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var Logger = backend.Logger

var (
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
	_ backend.StreamHandler         = (*Datasource)(nil)
)

func NewDatasource(_ context.Context, s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	Logger.Debug("datasource :: NewDatasource")
	settings, err := getDatasourceSettings(s)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (d *Datasource) Dispose() {
	Logger.Debug("datasource :: Dispose")
}

func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	Logger.Debug("datasource :: CheckHealth")

	var status = backend.HealthStatusOk
	var message = "Connected to Kafka"

	if !d.canConnect() {
		status = backend.HealthStatusError
		message = "Failed to establish connection to Kafka"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (d *Datasource) canConnect() bool {
	Logger.Debug("datasource :: canConnect")
	c, err := kafka.NewClient(d.kafkaconfig)
	if err != nil {
		return false
	}
	return c.Close() == nil
}

func (d *Datasource) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	Logger.Debug("datasource :: SubscribeStream")
	return &backend.SubscribeStreamResponse{
		Status: backend.SubscribeStreamStatusOK,
	}, nil
}

func (d *Datasource) PublishStream(context.Context, *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	Logger.Debug("datasource :: PublishStream")
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}

func (d *Datasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	Logger.Debug("datasource :: RunStream")

	var q Query
	if err := json.Unmarshal(req.Data, &q); err != nil {
		return err
	}
	Logger.Debug("datasource :: RunStream :: Query", "Query", q)

	if q.TopicName == "" || q.TopicName == "TOPIC_NAME" {
		Logger.Debug("datasource :: RunStream :: No topic name provided")
		return context.Cause(ctx)
	}

	Logger.Debug("datasource :: RunStream :: creating Kafka client")
	k, err := kafka.NewClient(d.kafkaconfig)
	if err != nil {
		Logger.Error("datasource :: RunStream :: Error creating Kafka client", "err", err)
		return err
	}

	defer func() {
		if err := k.Close(); err != nil {
			Logger.Error("datasource :: RunStream :: Error closing connection", "err", err)
		}
		Logger.Info("datasource :: RunStream :: Connection closed")
	}()

	Logger.Debug("datasource :: RunStream :: subscribing to topic", "topic", q.TopicName)
	if err := k.Subscribe(q.TopicName); err != nil {
		Logger.Error("datasource :: RunStream :: Error subscribing to topic", "err", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			Logger.Debug("Context cancelled - closing connection")
			return context.Canceled
		default:
			fatalerr, err := sendKafkaMessageToGrafana(k, sender)
			if fatalerr {
				return err
			}
		}
	}
}

func flattenJSON(prefix string, obj map[string]interface{}, fields *[]*data.Field) {
	// sort the keys to ensure a consistent order
	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// iterate over the sorted keys
	for _, key := range keys {
		value := obj[key]
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		switch v := value.(type) {
		case float64:
			*fields = append(*fields, data.NewField(fullKey, nil, []float64{v}))
		case string:
			*fields = append(*fields, data.NewField(fullKey, nil, []string{v}))
		case bool:
			*fields = append(*fields, data.NewField(fullKey, nil, []bool{v}))
		case int:
			*fields = append(*fields, data.NewField(fullKey, nil, []int{v}))
		case int64:
			*fields = append(*fields, data.NewField(fullKey, nil, []int64{v}))
		case []interface{}:
			*fields = append(*fields, data.NewField(fullKey, nil, []string{fmt.Sprintf("%v", v)}))
		case map[string]interface{}:
			flattenJSON(fullKey, v, fields)
		default:
			Logger.Warn("Ignoring unsupported data type", "key", fullKey, "value", value)
		}
	}
}

func sendKafkaMessageToGrafana(k *kafka.Client, sender *backend.StreamSender) (bool, error) {
	var jsonMessage map[string]interface{}

	kafkaMessage, err := k.ReadMessage()
	if err != nil {
		Logger.Error("Error reading message", "err", err)
		return true, err
	}

	rawMsg := kafkaMessage.Value
	Logger.Debug("Message received", "msg", string(rawMsg))

	if err := json.Unmarshal(rawMsg, &jsonMessage); err != nil {
		Logger.Error("Error unmarshalling message", "err", err)
		// not a fatal error - non-JSON messages are ignored
		return false, err
	}

	// create a Grafana data frame to represent the message
	fields := make([]*data.Field, 0, len(jsonMessage)+3)
	fields = append(fields, data.NewField("_eventtime", nil, []time.Time{kafkaMessage.Timestamp}))
	fields = append(fields, data.NewField("_partition", nil, []int32{kafkaMessage.Partition}))
	fields = append(fields, data.NewField("_offset", nil, []int64{kafkaMessage.Offset}))

	// flatten the JSON message
	flattenJSON("", jsonMessage, &fields)

	// submit to Grafana front-end
	frame := data.NewFrame("response", fields...)
	err = sender.SendFrame(frame, data.IncludeAll)
	if err != nil {
		Logger.Error("Error sending frame", "err", err)
		return true, err
	}

	return false, nil
}
