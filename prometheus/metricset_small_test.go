// Copyright 2018 Jump Trading
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build small

package prometheus_test

import (
	"testing"

	"github.com/jumptrading/influx-spout/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	set := prometheus.NewMetricSet()
	assert.Len(t, set.All(), 0)
}

func TestAddition(t *testing.T) {
	set := prometheus.NewMetricSet()

	m0 := &prometheus.Metric{
		Name:  []byte("foo"),
		Value: 1,
	}
	set.Update(m0)
	assert.Equal(t, set.All(), []*prometheus.Metric{m0})

	m1 := &prometheus.Metric{
		Name:  []byte("bar"),
		Value: 2,
	}
	set.Update(m1)
	assert.ElementsMatch(t, set.All(), []*prometheus.Metric{m0, m1})
}

func TestUpdate(t *testing.T) {
	set := prometheus.NewMetricSet()

	m := &prometheus.Metric{
		Name:         []byte("foo"),
		Value:        1,
		Milliseconds: 100,
	}
	set.Update(m)

	m.Value = 2
	m.Milliseconds = 200
	set.Update(m)
	assert.Equal(t, set.All(), []*prometheus.Metric{m})
}

func TestUpdateWithLabels(t *testing.T) {
	set := prometheus.NewMetricSet()

	m0 := &prometheus.Metric{
		Name: []byte("temp"),
		Labels: prometheus.LabelPairs{
			{Name: []byte("host"), Value: []byte("nyc02")},
			{Name: []byte("core"), Value: []byte("0")},
		},
		Value:        55,
		Milliseconds: 100,
	}
	set.Update(m0)

	m1 := &prometheus.Metric{
		Name: []byte("temp"),
		Labels: prometheus.LabelPairs{
			{Name: []byte("host"), Value: []byte("nyc02")},
			{Name: []byte("core"), Value: []byte("1")},
		},
		Value:        56,
		Milliseconds: 100,
	}
	set.Update(m1)

	m2 := &prometheus.Metric{
		Name: []byte("temp"),
		Labels: prometheus.LabelPairs{
			{Name: []byte("host"), Value: []byte("auk01")},
			{Name: []byte("core"), Value: []byte("0")},
		},
		Value:        44,
		Milliseconds: 100,
	}
	set.Update(m2)

	m3 := &prometheus.Metric{
		Name: []byte("temp"),
		Labels: prometheus.LabelPairs{
			{Name: []byte("host"), Value: []byte("auk02")},
			{Name: []byte("core"), Value: []byte("0")},
		},
		Value:        45,
		Milliseconds: 101,
	}
	set.Update(m3)

	// Labels vary for the same metric name so there are four separate
	// metrics.
	assert.ElementsMatch(t, set.All(), []*prometheus.Metric{m0, m1, m2, m3})

	// Now update some and ensure they are updated.
	m0.Value = 66
	m0.Milliseconds = 150
	set.Update(m0)

	m3.Value = 98
	m3.Milliseconds = 151
	set.Update(m3)

	assert.ElementsMatch(t, set.All(), []*prometheus.Metric{m0, m1, m2, m3})
}

func TestUpdateLabelOrdering(t *testing.T) {
	set := prometheus.NewMetricSet()

	m0 := &prometheus.Metric{
		Name: []byte("temp"),
		Labels: prometheus.LabelPairs{
			{Name: []byte("host"), Value: []byte("nyc02")},
			{Name: []byte("core"), Value: []byte("0")},
		},
		Value:        55,
		Milliseconds: 100,
	}
	set.Update(m0)

	// m1 has the same name and labels as m0 but the label ordering is
	// switched.
	m1 := &prometheus.Metric{
		Name: []byte("temp"),
		Labels: prometheus.LabelPairs{
			{Name: []byte("core"), Value: []byte("0")},
			{Name: []byte("host"), Value: []byte("nyc02")},
		},
		Value:        66,
		Milliseconds: 101,
	}
	set.Update(m1)

	assert.Equal(t, set.All(), []*prometheus.Metric{m1})
}