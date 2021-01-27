package collector

import (
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Metrics is the structure of metric for prometheus
type Metrics struct {
	metrics map[string]*prometheus.Desc
	mutex   sync.Mutex
}

func newGlobalMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}

// NewMetrics is the factory function to initial a Metrics struct
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"gauge_metric": newGlobalMetric(namespace, "gauge_metric", "Temperature of Raspberry Pi", []string{"type"}),
		},
	}
}

// Describe is the API of transfering Metrics struct to channel
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

// Collect is the API of collecting metrics
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	temperatureGaugeMetricData, err := c.GenerateTemperatureMetrics()
	if err != nil {
		log.Warn(err)
	}
	for gaugetype, currentValue := range temperatureGaugeMetricData {
		ch <- prometheus.MustNewConstMetric(c.metrics["gauge_metric"], prometheus.GaugeValue, float64(currentValue), gaugetype)
	}
}

// GenerateTemperatureMetrics is the function to get the CPU temperature of host
func (c *Metrics) GenerateTemperatureMetrics() (map[string]float64, error) {
	temp, err := getTemperature()
	if err != nil {
		return nil, err
	}
	return map[string]float64{"temperature": temp}, nil
}

func getTemperature() (float64, error) {
	s := "cat /sys/devices/virtual/thermal/thermal_zone0/temp"
	cmd := exec.Command("/bin/bash", "-c", s)

	//var out bytes.Buffer
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	tempValue := strings.Replace(string(out), "\n", "", -1)

	result, err := strconv.ParseFloat(tempValue, 64)
	if err != nil {
		return 0, err
	}
	return math.Trunc(result/10+0.5) * 1e-2, nil
}

func appendSlash(s string) string {
	if len(s) == 0 {
		return ""
	}
	b := []byte(s)
	if b[len(b)-1] != '/' {
		b = append(b, '/')
	}
	return string(b)
}

func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
