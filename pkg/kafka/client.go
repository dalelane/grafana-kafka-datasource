package kafka

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

var Logger = backend.Logger

type Client struct {
	bootstrapServers string
	config           *sarama.Config
	consumerGroup    sarama.ConsumerGroup
	messageChannel   chan *sarama.ConsumerMessage
	closeChannel     chan struct{}
	closeOnce        sync.Once
}

type SimpleClientConfig struct {
	BootstrapServers string `json:"bootstrapservers"`
	ClientId         string `json:"clientid"`
	GroupId          string `json:"groupid"`
	ScramAuthType    string `json:"authtype"`
	ScramUsername    string `json:"username"`
	ScramPassword    string `json:"password"`
	UseSslTls        bool   `json:"usetls"`
}

func NewClient(cfg SimpleClientConfig) (*Client, error) {
	Logger.Debug("kafka :: client :: NewClient",
		"bootstrap", cfg.BootstrapServers,
		"clientid", cfg.ClientId,
		"groupid", cfg.GroupId,
		"authtype", cfg.ScramAuthType)

	if cfg.ClientId == "" {
		cfg.ClientId = "grafana"
	}
	if cfg.GroupId == "" {
		cfg.GroupId = "grafana"
	}
	if cfg.ScramAuthType == "" {
		cfg.ScramAuthType = "none"
	}

	config := sarama.NewConfig()
	config.ClientID = cfg.ClientId
	if cfg.ScramAuthType != "none" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = cfg.ScramUsername
		config.Net.SASL.Password = cfg.ScramPassword
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{HashGeneratorFcn: GetHashGeneratorFn(cfg.ScramAuthType)}
		}
		config.Net.SASL.Mechanism = sarama.SASLMechanism(cfg.ScramAuthType)
	}
	config.Net.TLS.Enable = cfg.UseSslTls
	if cfg.UseSslTls {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = &tls.Config{InsecureSkipVerify: true}
	}
	config.Consumer.Offsets.AutoCommit.Enable = false

	consumerGroup, err := sarama.NewConsumerGroup([]string{cfg.BootstrapServers}, cfg.GroupId, config)
	if err != nil {
		Logger.Error("kafka ::: client :: NewConsumerGroup :: failed to create consumer group", "error", err)
		return nil, err
	}

	c := &Client{
		bootstrapServers: cfg.BootstrapServers,
		config:           config,
		consumerGroup:    consumerGroup,
		messageChannel:   make(chan *sarama.ConsumerMessage),
		closeChannel:     make(chan struct{}),
	}
	return c, nil
}

func (c *Client) Subscribe(topicname string) error {
	Logger.Debug("kafka :: client :: Subscribe", "topicname", topicname)

	oneHourAgo := time.Now().Add(-1 * time.Hour).UnixMilli()
	startingOffsets := make(map[int32]int64)

	// Attempt to get a suitable starting offset for each partition that gives
	//  a little history to start the visualisations off with
	// It isn't critical to do this, so if it fails, we just log the failure and
	//  start from the latest offset
	client, err := sarama.NewClient([]string{c.bootstrapServers}, c.config)
	if err != nil {
		Logger.Error("kafka :: client :: Subscribe :: failed to create Sarama client", "error", err)
		return nil
	}

	defer client.Close()

	partitions, err := client.Partitions(topicname)
	if err != nil {
		Logger.Error("kafka :: client :: Subscribe :: failed to fetch partitions info", "error", err)
	} else {
		Logger.Debug("kafka :: client :: Subscribe :: partitions", "num", len(partitions))
		for _, partition := range partitions {
			startingOffset, err := client.GetOffset(topicname, partition, oneHourAgo)
			if err != nil {
				Logger.Error("kafka :: client :: Subscribe :: failed to fetch starting offset", "partition", partition, "error", err)
				continue
			}
			if startingOffset == -1 {
				Logger.Error("kafka :: client :: Subscribe :: unexpected starting offset", "partition", partition)
				continue
			}

			Logger.Debug("kafka :: client :: Subscribe :: starting offset",
				"topic", topicname,
				"partition", partition,
				"start", startingOffset)

			startingOffsets[partition] = startingOffset
		}
	}

	messageHandler := &consumerGroupHandler{
		topicName:      topicname,
		startOffsets:   startingOffsets,
		messageChannel: c.messageChannel,
		closeChannel:   c.closeChannel,
	}
	go func() {
		ctx := context.Background()
		for {
			if err := c.consumerGroup.Consume(ctx, []string{topicname}, messageHandler); err != nil {
				Logger.Error("kafka :: client :: Subscribe :: failed to consume", "error", err)
				return
			}
			select {
			case <-c.closeChannel:
				Logger.Debug("kafka :: client :: Subscribe :: close", "topic", topicname)
				return
			default:
				// continue consuming
			}
		}
	}()

	return nil
}

func (c *Client) ReadMessage() (*sarama.ConsumerMessage, error) {
	Logger.Debug("kafka :: client :: ReadMessage")
	select {
	case msg := <-c.messageChannel:
		return msg, nil
	case <-c.closeChannel:
		Logger.Debug("kafka :: client :: ReadMessage :: close")
		return nil, nil
	}
}

func (c *Client) Close() error {
	Logger.Debug("kafka :: client :: Close")
	var err error
	c.closeOnce.Do(func() {
		close(c.closeChannel)
		Logger.Debug("kafka :: client :: Close :: closing consumer group")
		if closeErr := c.consumerGroup.Close(); closeErr != nil {
			Logger.Error("kafka :: client :: failed to close consumer group", "error", closeErr)
			err = closeErr
		}
	})
	return err
}

type consumerGroupHandler struct {
	topicName      string
	startOffsets   map[int32]int64
	messageChannel chan *sarama.ConsumerMessage
	closeChannel   chan struct{}
}

func (h *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	Logger.Debug("kafka :: client :: consumerGroupHandler :: Setup")
	for partition, offset := range h.startOffsets {
		Logger.Debug("kafka :: client :: consumerGroupHandler :: Setup", "partition", partition, "offset", offset)
		session.MarkOffset(h.topicName, partition, offset, "initial")
	}
	return nil
}
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	Logger.Debug("kafka :: client :: consumerGroupHandler :: Cleanup")
	return nil
}
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	Logger.Debug("kafka :: client :: consumerGroupHandler :: ConsumeClaim")
	for {
		select {
		case msg, ok := <-claim.Messages():
			if ok {
				Logger.Debug("kafka :: client :: consumerGroupHandler :: ConsumeClaim :: message", "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)
				h.messageChannel <- msg
				session.MarkMessage(msg, "")
			}
		case <-h.closeChannel:
			Logger.Debug("kafka :: client :: consumerGroupHandler :: ConsumeClaim :: close")
			return nil
		}
	}
}
