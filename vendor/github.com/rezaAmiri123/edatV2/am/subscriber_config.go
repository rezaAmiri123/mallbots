package am

import "time"

type AckType int

const (
	AckTypeAuto AckType = iota
	AckTypeManual
)

var (
	defaultAckWait      = 30 * time.Second
	defaultMaxRedeliver = 5
)

type (
	SubscriberOption interface {
		configureSubscriberConfig(*SubscriberConfig)
	}

	SubscriberConfig struct {
		msgFilter    []string
		groupName    string
		ackType      AckType
		ackWait      time.Duration
		maxRedeliver int
	}
)

func NewSubscriberConfig(options []SubscriberOption) SubscriberConfig {
	cfg := SubscriberConfig{
		msgFilter:    []string{},
		groupName:    "",
		ackType:      AckTypeManual,
		ackWait:      defaultAckWait,
		maxRedeliver: defaultMaxRedeliver,
	}

	for _, option := range options {
		option.configureSubscriberConfig(&cfg)
	}

	return cfg
}

func (c SubscriberConfig) MessageFilters() []string { return c.msgFilter }
func (c SubscriberConfig) GroupName() string        { return c.groupName }
func (c SubscriberConfig) AckType() AckType         { return c.ackType }
func (c SubscriberConfig) AckWait() time.Duration   { return c.ackWait }
func (c SubscriberConfig) MaxRedeliver() int        { return c.maxRedeliver }


func (h AckType) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.ackType = h }


type MessageFilter []string

func (h MessageFilter) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.msgFilter = h }


type GroupName string

func (h GroupName) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.groupName = string(h) }


type AckWait time.Duration

func (h AckWait) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.ackWait = time.Duration(h) }


type MaxDeliver int

func (h MaxDeliver) configureSubscriberConfig(cfg *SubscriberConfig) { cfg.maxRedeliver = int(h) }
