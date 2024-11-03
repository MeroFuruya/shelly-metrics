package shelly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/MeroFuruya/shelly-metrics/core/logging"
)

func ParseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000-0700", timestamp)
}

func ParseInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func Parse(v []byte) (MetricBatch, error) {
	logger := logging.GetLogger("shelly.parser")
	decoder := json.NewDecoder(bytes.NewReader(v))
	data := make(map[string]json.RawMessage)
	batch := MetricBatch{}

	if err := decoder.Decode(&data); err != nil {
		logger.Error().Err(err).Msg("Failed to decode JSON")
		return batch, err
	}

	for key, json_value := range data {
		var err error
		var bytes []byte
		if bytes, err = json_value.MarshalJSON(); err != nil {
			logger.Error().Err(err).Msg("Failed to marshal JSON")
			continue
		}
		value := string(bytes)

		ParseKeyValue(key, value, &batch)
	}
	return batch, nil
}

func ParseKeyValue(key string, value string, batch *MetricBatch) {
	logger := logging.GetLogger("shelly.parser.ParseKeyValue")
	switch key {
	case "relays_on":
		if number, err := ParseInt64(value); err != nil {
			logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to parse int64")
			return
		} else {
			batch.InsertMetric(RelaysOnMetric{
				RelaysOn:  number,
				Timestamp: time.Now(),
			})
		}
	case "lights_on":
		if number, err := ParseInt64(value); err != nil {
			logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to parse int64")
			return
		} else {
			batch.InsertMetric(LightsOnMetric{
				LightsOn:  number,
				Timestamp: time.Now(),
			})
		}
	case "total_devices":
		if number, err := ParseInt64(value); err != nil {
			logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to parse int64")
			return
		} else {
			batch.InsertMetric(DevicesMetric{
				Devices:   number,
				Timestamp: time.Now(),
			})
		}
	case "total_accounts":
		if number, err := ParseInt64(value); err != nil {
			logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to parse int64")
			return
		} else {
			batch.InsertMetric(AccountsMetric{
				Accounts:  number,
				Timestamp: time.Now(),
			})
		}
	case "almost_live_metrics":
		ParseAlmostLiveMetrics(value, batch)
	case "fast_api_calls_integrations_today":
		if number, err := ParseInt64(value); err != nil {
			logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to parse int64")
			return
		} else {
			batch.InsertMetric(TodayHomeAssistantCallsMetric{
				Calls:     number,
				Timestamp: time.Now(),
			})
		}
	case "today_traffic":
		number_str := regexp.MustCompile("[^0-9]").ReplaceAllString(value, "")
		if number, err := strconv.ParseInt(number_str, 10, 64); err != nil {
			logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to parse int64")
			return
		} else {
			batch.InsertMetric(TodayTrafficMetric{
				Traffic:   float64(number) / 10,
				Timestamp: time.Now(),
			})
		}
	case "today_new_devices":
		if number, err := ParseInt64(value); err != nil {
			logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to parse int64")
			return
		} else {
			batch.InsertMetric(TodayNewDevicesMetric{
				Devices:   number,
				Timestamp: time.Now(),
			})
		}
	case "today_new_accounts":
		if number, err := ParseInt64(value); err != nil {
			logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to parse int64")
			return
		} else {
			batch.InsertMetric(TodayNewAccountsMetric{
				Accounts:  number,
				Timestamp: time.Now(),
			})
		}
	case "executed_scenes_today":
		if number, err := ParseInt64(value); err != nil {
			logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to parse int64")
			return
		} else {
			batch.InsertMetric(TodayExecutedScenesMetric{
				Scenes:    number,
				Timestamp: time.Now(),
			})
		}
	case "jem_telemetry":
		telemetry := struct {
			DoExecs int64 `json:"do_execs"`
		}{}
		if err := json.Unmarshal([]byte(value), &telemetry); err != nil {
			logger.Error().Err(err).Msg("Failed to unmarshal JSON")
			return
		} else {
			batch.InsertMetric(MinuteExecutedScenesMetric{
				Scenes:    telemetry.DoExecs,
				Timestamp: time.Now(),
			})
		}
	case "traffic":
		number_str := regexp.MustCompile("[^0-9]").ReplaceAllString(value, "")
		if number, err := strconv.ParseInt(number_str, 10, 64); err != nil {
			logger.Error().Err(err).Str("key", key).Str("value", value).Msg("Failed to parse int64")
			return
		} else {
			batch.InsertMetric(MinuteTrafficMetric{
				Traffic:   float64(number) / 10,
				Timestamp: time.Now(),
			})
		}
	case "new_event":
		ParseNewEvent(value, batch)
	case "new_device":
		ParseNewDevice(value, batch)
	default:
		ParseUnknown(key, value, batch)
	}
}

func ParseAlmostLiveMetrics(value string, batch *MetricBatch) {
	logger := logging.GetLogger("shelly.parser.ParseAlmostLiveMetrics")
	var err error
	almost_live_metrics := struct {
		TotalPower          float64            `json:"totalPower"`
		TotalPowerByType    map[string]float64 `json:"totalPowerByType"`
		TotalEnergyConsumed float64            `json:"totalEnergyConsumed"`
		TotalEnergyProduced float64            `json:"totalEnergyProduced"`
	}{}

	if err = json.Unmarshal([]byte(value), &almost_live_metrics); err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal JSON")
		return
	}

	batch.InsertMetric(PowerMetric{
		Power:     almost_live_metrics.TotalPower,
		Timestamp: time.Now(),
	})
	batch.InsertMetric(MinuteEnergyConsumptionMetric{
		Energy:    almost_live_metrics.TotalEnergyConsumed,
		Timestamp: time.Now(),
	})
	batch.InsertMetric(MinuteEnergyProductionMetric{
		Energy:    almost_live_metrics.TotalEnergyProduced,
		Timestamp: time.Now(),
	})

	for key, value := range almost_live_metrics.TotalPowerByType {
		batch.InsertMetric(MinuteEnergyByTypeMetric{
			Type:      key,
			Energy:    value,
			Timestamp: time.Now(),
		})
	}
}

func ParseNewEvent(value string, batch *MetricBatch) {
	logger := logging.GetLogger("shelly.parser.ParseNewEvent")
	new_event := struct {
		EventType string `json:"type"`
		Location  struct {
			Latitude  float64 `json:"lat"`
			Longitude float64 `json:"lon"`
		} `json:"location"`
		Time time.Time `json:"time"`
	}{}

	if err := json.Unmarshal([]byte(value), &new_event); err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal JSON")
		return
	}

	batch.InsertMetric(EventMetric{
		Type:      new_event.EventType,
		Latitude:  new_event.Location.Latitude,
		Longitude: new_event.Location.Longitude,
		Timestamp: new_event.Time,
	})
}

func ParseNewDevice(value string, batch *MetricBatch) {
	logger := logging.GetLogger("shelly.parser.ParseNewDevice")
	new_device := struct {
		Type     string `json:"type"`
		Location struct {
			Latitude  float64 `json:"lat"`
			Longitude float64 `json:"lon"`
			Display   string  `json:"display"`
		} `json:"location"`
		Friendly string    `json:"friendly_name"`
		Family   string    `json:"family"`
		Time     time.Time `json:"time"`
	}{}

	if err := json.Unmarshal([]byte(value), &new_device); err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal JSON")
		return
	}

	batch.InsertMetric(NewDeviceMetric{
		Type:      new_device.Type,
		Latitude:  new_device.Location.Latitude,
		Longitude: new_device.Location.Longitude,
		Display:   new_device.Location.Display,
		Friendly:  new_device.Friendly,
		Family:    new_device.Family,
		Timestamp: new_device.Time,
	})
}

func ParseUnknown(key string, value string, batch *MetricBatch) {
	logger := logging.GetLogger("shelly.parser.ParseUnknown")
	logger.Warn().Str("key", key).Str("value", value).Msg("Unknown metric")
	batch.InsertMetric(UnknownMetric{
		Data:      fmt.Sprintf("%s: %s", key, value),
		Timestamp: time.Now(),
	})
}
