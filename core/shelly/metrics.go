package shelly

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// const (
// 	MetricTypeUnknown = iota
// 	//	{
// 	//		"relays_on": 807790
// 	//	}
// 	MetricTypeRelaysOn
// 	//	{
// 	//		"lights_on": 50507
// 	//	}
// 	MetricTypeLightOn
// 	//	{
// 	//		"total_devices": 11421748
// 	//	}
// 	MetricTypeDevices
// 	//	{
// 	//		"total_accounts": 11421748
// 	//	}
// 	MetricTypeAccounts
// 	// in 1W increments
// 	//	{
// 	//		"almost_live_metrics": {
// 	//			"totalPower": 148589111.33000004,
// 	//		}
// 	//	}
// 	MetricTypePower
// 	//	{
// 	//		"fast_api_calls_integrations_today": 1914661
// 	//	}
// 	MetricTypeTodayHomeAssistantCalls
// 	// In 0.1KB increments
// 	//	{
// 	//		"today_traffic": "38092892778454"
// 	//	}
// 	MetricTypeTodayTraffic
// 	//	{
// 	//		"today_new_devices": 11421748
// 	//	}
// 	MetricTypeTodayNewDevices
// 	//	{
// 	//		"today_new_accounts": 2019
// 	//	}
// 	MetricTypeTodayNewAccounts
// 	//	{
// 	//		"executed_scenes_today": 99427048
// 	//	}
// 	MetricTypeTodayExecutedScenes
// 	//	{
// 	//		"jem_telemetry": {
// 	//			"do_execs": 47758
// 	//		}
// 	// }
// 	MetricTypeMinuteExecutedScenes
// 	// in 1Wh increments
// 	//	{
// 	//		"almost_live_metrics": {
// 	//			"totalPowerByType": {
// 	//				"<type>": 15602367.833183136,
// 	//			}
// 	//		}
// 	//	}
// 	MetricTypeMinuteEnergyByType
// 	// in 1Wh increments
// 	//	{
// 	//		"almost_live_metrics": {
// 	//			"totalEnergyConsumed": 120160473.66000001,
// 	//		}
// 	//	}
// 	MetricTypeMinuteEnergyConsumption
// 	// in 1Wh increments
// 	//	{
// 	//		"almost_live_metrics": {
// 	//			"totalEnergyProduced": -16062081.469999999,
// 	//		}
// 	//	}
// 	MetricTypeMinuteEnergyProduction
// 	//	{
// 	//		"new_event": {
// 	//			"type": "flood",
// 	//			"location": {
// 	//				"lat": 32.080299,
// 	// 				"lon": 34.780499
// 	//			},
// 	//		},
// 	//		"time": "2024-11-01T23:32:54.368Z"
// 	//	}
// 	MetricTypeEvent
// 	//	{
// 	//		"new_device": {
// 	//			"type": "SNDM-0013US",
// 	//			"location": {
// 	//				"lat": 32.080299,
// 	// 				"lon": 34.780499,
// 	//				"display": "US, New York"
// 	//			},
// 	//		"time": "2024-11-01T23:32:54.368Z"
// 	//		"friendly_name": "Plus Wall Dimmer",
// 	//		"family": "light"
// 	//		},
// 	MetricTypeNewDevice
// )

type UnknownMetric struct {
	Data      string
	Timestamp time.Time
}

//	{
//		"relays_on": 807790
//	}
type RelaysOnMetric struct {
	RelaysOn  int64
	Timestamp time.Time
}

//	{
//		"lights_on": 50507
//	}
type LightsOnMetric struct {
	LightsOn  int64
	Timestamp time.Time
}

//	{
//		"total_devices": 11421748
//	}
type DevicesMetric struct {
	Devices   int64
	Timestamp time.Time
}

//	{
//		"total_accounts": 11421748
//	}
type AccountsMetric struct {
	Accounts  int64
	Timestamp time.Time
}

// in 1W increments
//
//	{
//		"almost_live_metrics": {
//			"totalPower": 148589111.33000004,
//		}
//	}
type PowerMetric struct {
	Power     float64
	Timestamp time.Time
}

//	{
//		"fast_api_calls_integrations_today": 1914661
//	}
type TodayHomeAssistantCallsMetric struct {
	Calls     int64
	Timestamp time.Time
}

// In 0.1KB increments
//
//	{
//		"today_traffic": "38092892778454"
//	}
type TodayTrafficMetric struct {
	Traffic   float64
	Timestamp time.Time
}

//	{
//		"today_new_devices": 11421748
//	}
type TodayNewDevicesMetric struct {
	Devices   int64
	Timestamp time.Time
}

//	{
//		"today_new_accounts": 2019
//	}
type TodayNewAccountsMetric struct {
	Accounts  int64
	Timestamp time.Time
}

//	{
//		"executed_scenes_today": 99427048
//	}
type TodayExecutedScenesMetric struct {
	Scenes    int64
	Timestamp time.Time
}

//		{
//			"jem_telemetry": {
//				"do_execs": 47758
//			}
//	}
type MinuteExecutedScenesMetric struct {
	Scenes    int64
	Timestamp time.Time
}

// In 0.1KB increments
//
//	{
//		"traffic": "38092892778454"
//	}
type MinuteTrafficMetric struct {
	Traffic   float64
	Timestamp time.Time
}

// in 1Wh increments
//
//	{
//		"almost_live_metrics": {
//			"totalPowerByType": {
//				"<type>": 15602367.833183136,
//			}
//		}
//	}
type MinuteEnergyByTypeMetric struct {
	Type      string
	Energy    float64
	Timestamp time.Time
}

// in 1Wh increments
//
//	{
//		"almost_live_metrics": {
//			"totalEnergyConsumed": 120160473.66000001,
//		}
//	}
type MinuteEnergyConsumptionMetric struct {
	Energy    float64
	Timestamp time.Time
}

// in 1Wh increments
//
//	{
//		"almost_live_metrics": {
//			"totalEnergyProduced": -16062081.469999999,
//		}
//	}
type MinuteEnergyProductionMetric struct {
	Energy    float64
	Timestamp time.Time
}

//	{
//		"new_event": {
//			"type": "flood",
//			"location": {
//				"lat": 32.080299,
//				"lon": 34.780499
//			},
//		"time": "2024-11-01T23:32:54.368Z"
//	}
type EventMetric struct {
	Type      string
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

//	{
//		"new_device": {
//			"type": "SNDM-0013US",
//			"location": {
//				"lat": 32.080299,
//				"lon": 34.780499,
//				"display": "US, New York"
//			},
//		"time": "2024-11-01T23:32:54.368Z"
//		"friendly_name": "Plus Wall Dimmer",
//		"family": "light"
//		},
//	}
type NewDeviceMetric struct {
	Type      string
	Latitude  float64
	Longitude float64
	Display   string
	Friendly  string
	Family    string
	Timestamp time.Time
}

var Metrics = []interface{}{
	UnknownMetric{},
	RelaysOnMetric{},
	LightsOnMetric{},
	DevicesMetric{},
	AccountsMetric{},
	PowerMetric{},
	TodayHomeAssistantCallsMetric{},
	TodayTrafficMetric{},
	TodayNewDevicesMetric{},
	TodayNewAccountsMetric{},
	TodayExecutedScenesMetric{},
	MinuteExecutedScenesMetric{},
	MinuteTrafficMetric{},
	MinuteEnergyByTypeMetric{},
	MinuteEnergyConsumptionMetric{},
	MinuteEnergyProductionMetric{},
	EventMetric{},
	NewDeviceMetric{},
}

func CreateMetricTable(metric interface{}, batch *pgx.Batch) {
	var columns []string = []string{}
	t := reflect.TypeOf(metric)
	for _, field := range reflect.VisibleFields(t) {
		sql := fmt.Sprintf("%s %s", field.Name, MapTypeToPgType(field.Type))
		columns = append(columns, sql)
	}
	batch.Queue(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", t.Name(), strings.Join(columns, ",")))
}

func CreateMetricHypertable(metric interface{}, batch *pgx.Batch) {
	t := reflect.TypeOf(metric)
	for _, field := range reflect.VisibleFields(t) {
		if field.Type == reflect.TypeOf(time.Time{}) {
			batch.Queue(fmt.Sprintf("SELECT create_hypertable('%s', 'timestamp', if_not_exists => TRUE)", t.Name()))
		}
	}
}

func CreateMetricTables(batch *pgx.Batch) {
	for _, metric := range Metrics {
		CreateMetricTable(metric, batch)
	}
}

func CreateMetricHyperTablesBatch(batch *pgx.Batch) {
	batch.Queue("CREATE EXTENSION IF NOT EXISTS timescaledb")
	for _, metric := range Metrics {
		CreateMetricHypertable(metric, batch)
	}
}

func MapTypeToPgType(t reflect.Type) string {
	if t == reflect.TypeOf(time.Time{}) {
		return "timestamptz"
	}

	switch t.Kind() {
	case reflect.Int:
		return "bigint"
	case reflect.Int64:
		return "bigint"
	case reflect.Float64:
		return "float8"
	case reflect.String:
		return "text"
	}

	return "text"
}

func InsertMetric(metric interface{}, batch *pgx.Batch) {
	t := reflect.TypeOf(metric)
	v := reflect.ValueOf(metric)
	var columns []string = []string{}
	var values_str []string = []string{}
	var values []any = []any{}
	for i, field := range reflect.VisibleFields(t) {
		value := v.Field(i)
		columns = append(columns, field.Name)
		values_str = append(values_str, fmt.Sprintf("$%d", i+1))
		values = append(values, value.Interface())
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", t.Name(), strings.Join(columns, ","), strings.Join(values_str, ","))
	batch.Queue(sql, values...)
}

type MetricBatch struct {
	Batch pgx.Batch
}

func NewMetricBatch() *MetricBatch {
	return &MetricBatch{
		Batch: pgx.Batch{},
	}
}

func (m *MetricBatch) InsertMetric(metric interface{}) {
	InsertMetric(metric, &m.Batch)
}
