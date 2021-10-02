package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "starlink"
)

var (
	dishUp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "up"),
		"Was the last query of Starlink dish successful.",
		nil, nil,
	)
	dishScrapeDurationSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "scrape_duration_seconds"),
		"Time to scrape metrics from starlink dish",
		nil, nil,
	)

	// collectDishContext
	dishInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "info"),
		"Running software versions and IDs of hardware",
		[]string{"device_id", "hardware_version", "software_version", "country_code", "utc_offset"}, nil,
	)
	dishUptimeSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "uptime_seconds"),
		"Dish running time",
		nil, nil,
	)
	dishCellId = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "cell_id"),
		"Cell ID dish is located in",
		nil, nil,
	)
	dishPopRackId = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "pop_rack_id"),
		"pop rack id",
		nil, nil,
	)
	dishInitialSatelliteId = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "initial_satellite_id"),
		"initial satellite id",
		nil, nil,
	)
	dishInitialGatewayId = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "initial_gateway_id"),
		"initial gateway id",
		nil, nil,
	)
	dishOnBackupBeam = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "backup_beam"),
		"connected to backup beam",
		nil, nil,
	)
	dishSecondsToSlotEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "time_to_slot_end_seconds"),
		"Seconds left on current slot",
		nil, nil,
	)

	// collectDishStatus
	dishState = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "state"),
		"The current dishState of the Dish (Unknown, Booting, Searching, Connected).",
		nil, nil,
	)
	dishSecondsToFirstNonemptySlot = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "first_nonempty_slot_seconds"),
		"Seconds to next non empty slot",
		nil, nil,
	)
	dishPopPingDropRatio = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "pop_ping_drop_ratio"),
		"Percent of pings dropped",
		nil, nil,
	)
	dishPopPingLatencySeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "pop_ping_latency_seconds"),
		"Latency of connection in seconds",
		nil, nil,
	)
	dishSnr = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "snr"),
		"Signal strength of the connection",
		nil, nil,
	)
	dishUplinkThroughputBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "uplink_throughput_bytes"),
		"Amount of bandwidth in bytes per second upload",
		nil, nil,
	)
	dishDownlinkThroughputBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "downlink_throughput_bytes"),
		"Amount of bandwidth in bytes per second download",
		nil, nil,
	)

	dishBoreSightAzimuthDeg = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "bore_sight_azimuth_deg"),
		"azimuth in degrees",
		nil, nil,
	)

	dishBoreSightElevationDeg = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "bore_sight_elevation_deg"),
		"elevation in degrees",
		nil, nil,
	)

	// collectDishObstructions
	dishCurrentlyObstructed = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "currently_obstructed"),
		"Status of view of the sky",
		nil, nil,
	)
	dishFractionObstructionRatio = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "fraction_obstruction_ratio"),
		"Percentage of obstruction",
		nil, nil,
	)
	dishLast24hObstructedSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "last_24h_obstructed_seconds"),
		"Number of seconds view of sky has been obstructed in the last 24hours",
		nil, nil,
	)
	dishValidSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "valid_seconds"),
		"Unknown",
		nil, nil,
	)
	dishProlongedObstructionDurationSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "prolonged_obstruction_duration_seconds"),
		"Average in seconds of prolonged obstructions",
		nil, nil,
	)
	dishProlongedObstructionIntervalSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "prolonged_obstruction_interval_seconds"),
		"Average prolonged obstruction interval in seconds",
		nil, nil,
	)
	dishWedgeFractionObstructionRatio = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "wedge_fraction_obstruction_ratio"),
		"Percentage of obstruction per wedge section",
		[]string{"wedge", "wedge_name"}, nil,
	)
	dishWedgeAbsFractionObstructionRatio = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "wedge_abs_fraction_obstruction_ratio"),
		"Percentage of Absolute fraction per wedge section",
		[]string{"wedge", "wedge_name"}, nil,
	)

	// collectDishAlerts
	dishAlertMotorsStuck = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "alert_motors_stuck"),
		"Status of motor stuck",
		nil, nil,
	)
	dishAlertThermalThrottle = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "alert_thermal_throttle"),
		"Status of thermal throttling",
		nil, nil,
	)
	dishAlertThermalShutdown = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "alert_thermal_shutdown"),
		"Status of thermal shutdown",
		nil, nil,
	)
	dishAlertMastNotNearVertical = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "alert_mast_not_near_vertical"),
		"Status of mast position",
		nil, nil,
	)
	dishUnexpectedLocation = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "alert_unexpected_location"),
		"Status of location",
		nil, nil,
	)
	dishSlowEthernetSpeeds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "dish", "alert_slow_eth_speeds"),
		"Status of ethernet",
		nil, nil,
	)
)
