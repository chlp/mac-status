package main

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v3/host"
	"log"
	"net/http"
	"sync"
	"time"
)

type Metric struct {
	Name   string `json:"name"`
	Value  int    `json:"value"`
	Unit   string `json:"unit"`
	Status string `json:"status"`
}

type MetricsData struct {
	MaxTemp      int
	SystemUptime int
}

var (
	mu             sync.RWMutex
	metrics        = MetricsData{}
	startTime      = time.Now()
	updateInterval = 5 * time.Second
)

func main() {
	go monitorTemperature()
	go monitorSystemUptime()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()

		appUptime := int(time.Since(startTime).Seconds())

		response := []Metric{
			{
				Name:   "max_temperature",
				Value:  metrics.MaxTemp,
				Unit:   "celsius",
				Status: determineTempStatus(metrics.MaxTemp),
			},
			{
				Name:   "system_uptime",
				Value:  metrics.SystemUptime,
				Unit:   "seconds",
				Status: "none",
			},
			{
				Name:   "app_uptime",
				Value:  appUptime,
				Unit:   "seconds",
				Status: "none",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	})

	log.Println("HTTP server started on http://localhost:8093")
	log.Fatal(http.ListenAndServe(":8093", nil))
}

// monitorTemperature periodically updates the MaxTemp value
func monitorTemperature() {
	for {
		updateMaxTemperature()
		time.Sleep(updateInterval)
	}
}

// monitorSystemUptime periodically updates the SystemUptime value
func monitorSystemUptime() {
	for {
		updateSystemUptime()
		time.Sleep(updateInterval)
	}
}

// updateMaxTemperature scans sensor temperatures and updates MaxTemp if a higher value is found
func updateMaxTemperature() {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		log.Printf("Error getting temperatures: %v", err)
		return
	}

	var localMax float64
	for _, t := range temps {
		if t.Temperature > localMax && t.Temperature < 100 {
			localMax = t.Temperature
		}
	}

	mu.Lock()
	metrics.MaxTemp = int(localMax)
	mu.Unlock()
}

// updateSystemUptime retrieves and updates the system uptime in seconds
func updateSystemUptime() {
	bootTime, err := host.BootTime()
	if err != nil {
		log.Printf("Error getting boot time: %v", err)
		return
	}
	uptime := int(time.Since(time.Unix(int64(bootTime), 0)).Seconds())

	mu.Lock()
	metrics.SystemUptime = uptime
	mu.Unlock()
}

// determineTempStatus returns a status string based on the temperature value
func determineTempStatus(temp int) string {
	switch {
	case temp >= 85:
		return "critical"
	case temp >= 70:
		return "warn"
	case temp > 0:
		return "ok"
	default:
		return "none"
	}
}
