package config

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zaffka/aggregator/domain"
)

func Test_readConfig(t *testing.T) {
	expectedConf := &config{
		Gens: []generator{
			{
				TimeoutS:       30,
				SendingPeriodS: 1,
				DataSources: []dataSource{
					{ID: "data_1", InitValue: 50, MaxChangeStep: 5},
					{ID: "data_2", InitValue: 60, MaxChangeStep: 7},
					{ID: "data_5", InitValue: 10, MaxChangeStep: 3},
				},
			},
			{
				TimeoutS:       20,
				SendingPeriodS: 1,
				DataSources: []dataSource{
					{ID: "data_7", InitValue: 20, MaxChangeStep: 4},
					{ID: "data_8", InitValue: 30, MaxChangeStep: 6},
				},
			},
		},
		Aggrs: []aggregator{
			{
				WorkDurationS: 10,
				SubSources:    []string{"data_1", "data_7"},
			},
			{
				WorkDurationS: 20,
				SubSources:    []string{"data_2", "data_5"},
			},
		},
		MsgQueueLength: 50,
		StorageType:    0,
	}

	tests := []struct {
		name       string
		path       string
		confStruct *config
		errMsg     string
	}{
		{
			"no file",
			"nofile.json",
			nil,
			"failed to read config file",
		},
		{
			"unmarshalling fail",
			"bad_conf_file.json",
			&config{},
			"failed to unmarshal config file",
		},
		{
			"ok config",
			"ok_conf_file.json",
			&config{},
			"",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := readConfig(tt.path, tt.confStruct)
			if err != nil {
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.Equal(t, expectedConf, tt.confStruct)
			}
		})
	}
}

func TestGet(t *testing.T) {
	_, err := Get("ok_conf_file.json")
	assert.NoError(t, err)
	pntr1 := fmt.Sprintf("%p", configuration)
	t.Log(pntr1)

	_, err = Get("ok_conf_file.json")
	assert.NoError(t, err)
	pntr2 := fmt.Sprintf("%p", configuration)
	t.Log(pntr1)

	assert.Equal(t, pntr1, pntr2)
}

func Test_config_Generators(t *testing.T) {
	type fields struct {
		Gens []generator
	}
	tests := []struct {
		name   string
		fields fields
		want   []domain.Generator
	}{
		{
			"empty generators slice",
			fields{
				Gens: []generator{},
			},
			[]domain.Generator{},
		},
		{
			"nil generators slice",
			fields{
				Gens: nil,
			},
			[]domain.Generator{},
		},
		{
			"ok generators slice",
			fields{
				Gens: []generator{
					{
						10,
						20,
						[]dataSource{
							{"test1", 5, 3},
							{"test3", 6, 4},
						},
					},
					{
						30,
						40,
						[]dataSource{
							{"test2", 5, 3},
							{"test4", 6, 4},
						},
					},
				},
			},
			[]domain.Generator{
				{Topic: "test1", StartValue: 5, MaxStep: 3, GenPeriod: 20 * time.Second, WorkDuration: 10 * time.Second},
				{Topic: "test3", StartValue: 6, MaxStep: 4, GenPeriod: 20 * time.Second, WorkDuration: 10 * time.Second},
				{Topic: "test2", StartValue: 5, MaxStep: 3, GenPeriod: 40 * time.Second, WorkDuration: 30 * time.Second},
				{Topic: "test4", StartValue: 6, MaxStep: 4, GenPeriod: 40 * time.Second, WorkDuration: 30 * time.Second},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := &config{
				Gens: tt.fields.Gens,
			}
			if got := c.Generators(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("config.Generators() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_config_Aggregators(t *testing.T) {
	type fields struct {
		Aggrs []aggregator
	}
	tests := []struct {
		name   string
		fields fields
		want   []domain.Aggregator
	}{
		{
			"empty aggregators slice",
			fields{
				Aggrs: []aggregator{},
			},
			[]domain.Aggregator{},
		},
		{
			"nil aggregators slice",
			fields{
				Aggrs: nil,
			},
			[]domain.Aggregator{},
		},
		{
			"ok aggregators slice",
			fields{
				Aggrs: []aggregator{
					{20, []string{"test2", "test1"}},
					{10, []string{"test3"}},
				},
			},
			[]domain.Aggregator{
				{
					ID:           "aggregator_1",
					Topic:        "test2",
					WorkDuration: 20 * time.Second,
				},
				{
					ID:           "aggregator_1",
					Topic:        "test1",
					WorkDuration: 20 * time.Second,
				},
				{
					ID:           "aggregator_2",
					Topic:        "test3",
					WorkDuration: 10 * time.Second,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c := &config{
				Aggrs: tt.fields.Aggrs,
			}
			if got := c.Aggregators(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("config.Aggregators() = %v, want %v", got, tt.want)
			}
		})
	}
}
