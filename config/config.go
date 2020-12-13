package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/zaffka/aggregator/domain"
)

var (
	once          sync.Once
	configuration *config
)

type config struct {
	Gens           []generator  `json:"generators"`
	Aggrs          []aggregator `json:"aggregators"`
	MsgQueueLength int          `json:"msg_queue_length"`
	StorageType    int          `json:"storage_type"`
}

type generator struct {
	TimeoutS       int          `json:"timeout_s"`
	SendingPeriodS int          `json:"sending_period_s"`
	DataSources    []dataSource `json:"data_sources"`
}

type dataSource struct {
	ID            string `json:"id"`
	InitValue     int    `json:"init_value"`
	MaxChangeStep int    `json:"max_change_step"`
}

type aggregator struct {
	WorkDurationS int      `json:"work_duration_s"`
	SubSources    []string `json:"sub_sources"`
}

// Get is the only function to acquire configuration handler interface(domain.Configurer).
func Get(confFilePath string) (domain.Configurer, error) {
	var err error
	once.Do(func() {
		configuration = &config{}
		err = readConfig(confFilePath, configuration)
	})

	return configuration, err
}

func readConfig(confFilePath string, confStruct *config) error {
	bts, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal(bts, confStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return nil
}

func (c *config) Generators() []domain.Generator {
	res := []domain.Generator{}

	for _, g := range c.Gens {
		for _, src := range g.DataSources {
			res = append(res, domain.Generator{
				Topic:        src.ID,
				StartValue:   src.InitValue,
				MaxStep:      src.MaxChangeStep,
				GenPeriod:    time.Duration(g.SendingPeriodS) * time.Second,
				WorkDuration: time.Duration(g.TimeoutS) * time.Second,
			})
		}
	}

	return res
}

func (c *config) Aggregators() []domain.Aggregator {
	res := []domain.Aggregator{}

	for i, a := range c.Aggrs {
		for _, src := range a.SubSources {
			res = append(res, domain.Aggregator{
				ID:           fmt.Sprintf("aggregator_%d", i+1),
				Topic:        src,
				WorkDuration: time.Duration(a.WorkDurationS) * time.Second,
			})
		}
	}

	return res
}

func (c *config) QueueTopics() []string {
	res := []string{}

	for _, g := range c.Gens {
		for _, src := range g.DataSources {
			res = append(res, src.ID)
		}
	}

	for _, a := range c.Aggrs {
		res = append(res, a.SubSources...)
	}

	return res
}

func (c *config) QueueLen() int { return c.MsgQueueLength }
func (c *config) Storage() int  { return c.StorageType }
