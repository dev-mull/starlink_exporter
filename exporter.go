package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"

	"github.com/dev-mull/starlink_exporter/pkg/spacex.com/api/device"
)

// Exporter collects Starlink stats from the Dish and exports them using
// the prometheus metrics package.
type Exporter struct {
	address     string
	conn        *grpc.ClientConn // keep the conn internal
	client      device.DeviceClient
	DishID      string
	CountryCode string
}

// New returns an initialized Exporter.
func NewExporter(address string) (*Exporter, error) {
	e := &Exporter{
		address: address,
	}
	return e, e.Conn()
}

func (e *Exporter) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

func (e *Exporter) GetState() connectivity.State {
	if e.conn != nil {
		return e.conn.GetState()
	}
	return 1 //If we are not connected, just say so
}

func (e *Exporter) Conn() error {
	if e.conn != nil {
		return nil
	}
	ctx, connCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer connCancel()
	conn, err := grpc.DialContext(ctx, e.address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("error creating underlying gRPC connection to starlink dish: %s", err.Error())
	}

	ctx, HandleCancel := context.WithTimeout(context.Background(), time.Second*1)
	defer HandleCancel()
	resp, err := device.NewDeviceClient(conn).Handle(ctx, &device.Request{
		Request: &device.Request_GetDeviceInfo{},
	})
	if err != nil {
		return fmt.Errorf("could not collect initial information from dish: %s", err.Error())
	}
	e.conn = conn
	e.client = device.NewDeviceClient(conn)
	e.DishID = resp.GetGetDeviceInfo().GetDeviceInfo().GetId()
	e.CountryCode = resp.GetGetDeviceInfo().GetDeviceInfo().GetCountryCode()
	return nil
}

// Describe describes all the metrics ever exported by the Starlink exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- dishUp
	ch <- dishScrapeDurationSeconds

	// collectDishContext
	ch <- dishInfo
	ch <- dishUptimeSeconds
	ch <- dishCellId
	ch <- dishPopRackId
	ch <- dishInitialSatelliteId
	ch <- dishInitialGatewayId
	ch <- dishOnBackupBeam
	ch <- dishSecondsToSlotEnd

	// collectDishStatus
	ch <- dishState
	ch <- dishSecondsToFirstNonemptySlot
	ch <- dishPopPingDropRatio
	ch <- dishPopPingLatencySeconds
	ch <- dishSnr
	ch <- dishUplinkThroughputBytes
	ch <- dishDownlinkThroughputBytes
	ch <- dishBoreSightAzimuthDeg
	ch <- dishBoreSightElevationDeg

	// collectDishObstructions
	ch <- dishCurrentlyObstructed
	ch <- dishFractionObstructionRatio
	ch <- dishLast24hObstructedSeconds
	ch <- dishValidSeconds
	ch <- dishProlongedObstructionDurationSeconds
	ch <- dishProlongedObstructionIntervalSeconds
	ch <- dishWedgeFractionObstructionRatio
	ch <- dishWedgeAbsFractionObstructionRatio

	// collectDishAlerts
	ch <- dishAlertMotorsStuck
	ch <- dishAlertThermalThrottle
	ch <- dishAlertThermalShutdown
	ch <- dishAlertMastNotNearVertical
	ch <- dishUnexpectedLocation
	ch <- dishSlowEthernetSpeeds
}

// Collect fetches the stats from Starlink dish and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if err := e.Conn(); err != nil {
		return
	}
	start := time.Now()

	ok := e.collectDishContext(ch)
	ok = ok && e.collectDishStatus(ch)
	ok = ok && e.collectDishObstructions(ch)
	ok = ok && e.collectDishAlerts(ch)

	if ok {
		ch <- prometheus.MustNewConstMetric(
			dishUp, prometheus.GaugeValue, 1.0,
		)
		ch <- prometheus.MustNewConstMetric(
			dishScrapeDurationSeconds, prometheus.GaugeValue, time.Since(start).Seconds(),
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			dishUp, prometheus.GaugeValue, 0.0,
		)
	}
}

func (e *Exporter) collectDishContext(ch chan<- prometheus.Metric) bool {
	req := &device.Request{
		Request: &device.Request_DishGetContext{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	resp, err := e.client.Handle(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() != 7 {
			log.Errorf("failed to collect dish context: %s", err.Error())
			return false
		}
	}

	dishC := resp.GetDishGetContext()
	dishI := dishC.GetDeviceInfo()
	dishS := dishC.GetDeviceState()

	ch <- prometheus.MustNewConstMetric(
		dishInfo, prometheus.GaugeValue, 1.00,
		dishI.GetId(),
		dishI.GetHardwareVersion(),
		dishI.GetSoftwareVersion(),
		dishI.GetCountryCode(),
		fmt.Sprint(dishI.GetUtcOffsetS()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishUptimeSeconds, prometheus.GaugeValue, float64(dishS.GetUptimeS()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishCellId, prometheus.GaugeValue, float64(dishC.GetCellId()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishPopRackId, prometheus.GaugeValue, float64(dishC.GetPopRackId()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishInitialSatelliteId, prometheus.GaugeValue, float64(dishC.GetInitialSatelliteId()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishInitialGatewayId, prometheus.GaugeValue, float64(dishC.GetInitialGatewayId()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishOnBackupBeam, prometheus.GaugeValue, flool(dishC.GetOnBackupBeam()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishSecondsToSlotEnd, prometheus.GaugeValue, float64(dishC.GetSecondsToSlotEnd()),
	)

	return true
}

func (e *Exporter) collectDishStatus(ch chan<- prometheus.Metric) bool {
	req := &device.Request{
		Request: &device.Request_GetStatus{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	resp, err := e.client.Handle(ctx, req)
	if err != nil {
		log.Errorf("failed to collect status from dish: %s", err.Error())
		return false
	}

	dishStatus := resp.GetDishGetStatus()

	ch <- prometheus.MustNewConstMetric(
		dishState, prometheus.GaugeValue, float64(dishStatus.GetState().Number()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishSecondsToFirstNonemptySlot, prometheus.GaugeValue, float64(dishStatus.GetSecondsToFirstNonemptySlot()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishPopPingDropRatio, prometheus.GaugeValue, float64(dishStatus.GetPopPingDropRate()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishPopPingLatencySeconds, prometheus.GaugeValue, float64(dishStatus.GetPopPingLatencyMs()/1000),
	)

	ch <- prometheus.MustNewConstMetric(
		dishSnr, prometheus.GaugeValue, float64(dishStatus.GetSnr()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishUplinkThroughputBytes, prometheus.GaugeValue, float64(dishStatus.GetUplinkThroughputBps()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishDownlinkThroughputBytes, prometheus.GaugeValue, float64(dishStatus.GetDownlinkThroughputBps()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishBoreSightAzimuthDeg, prometheus.GaugeValue, float64(dishStatus.GetBoresightAzimuthDeg()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishBoreSightElevationDeg, prometheus.GaugeValue, float64(dishStatus.GetBoresightElevationDeg()),
	)

	return true
}

func (e *Exporter) collectDishObstructions(ch chan<- prometheus.Metric) bool {
	req := &device.Request{
		Request: &device.Request_GetStatus{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	resp, err := e.client.Handle(ctx, req)
	if err != nil {
		log.Errorf("failed to collect obstructions from dish: %s", err.Error())
		return false
	}

	obstructions := resp.GetDishGetStatus().GetObstructionStats()

	ch <- prometheus.MustNewConstMetric(
		dishCurrentlyObstructed, prometheus.GaugeValue, flool(obstructions.GetCurrentlyObstructed()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishFractionObstructionRatio, prometheus.GaugeValue, float64(obstructions.GetFractionObstructed()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishLast24hObstructedSeconds, prometheus.GaugeValue, float64(obstructions.GetLast_24HObstructedS()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishValidSeconds, prometheus.GaugeValue, float64(obstructions.GetValidS()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishProlongedObstructionDurationSeconds, prometheus.GaugeValue, float64(obstructions.GetAvgProlongedObstructionDurationS()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishProlongedObstructionIntervalSeconds, prometheus.GaugeValue, float64(obstructions.GetAvgProlongedObstructionIntervalS()),
	)

	for i, v := range obstructions.GetWedgeFractionObstructed() {
		ch <- prometheus.MustNewConstMetric(
			dishWedgeFractionObstructionRatio, prometheus.GaugeValue, float64(v),
			strconv.Itoa(i),
			fmt.Sprintf("%d_to_%d", i*30, (i+1)*30),
		)
	}

	for i, v := range obstructions.GetWedgeAbsFractionObstructed() {
		ch <- prometheus.MustNewConstMetric(
			dishWedgeAbsFractionObstructionRatio, prometheus.GaugeValue, float64(v),
			strconv.Itoa(i),
			fmt.Sprintf("%d_to_%d", i*30, (i+1)*30),
		)
	}

	return true
}

func (e *Exporter) collectDishAlerts(ch chan<- prometheus.Metric) bool {
	req := &device.Request{
		Request: &device.Request_GetStatus{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	resp, err := e.client.Handle(ctx, req)
	if err != nil {
		log.Errorf("failed to collect alerts from dish: %s", err.Error())
		return false
	}

	alerts := resp.GetDishGetStatus().GetAlerts()

	ch <- prometheus.MustNewConstMetric(
		dishAlertMotorsStuck, prometheus.GaugeValue, flool(alerts.GetMotorsStuck()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishAlertThermalThrottle, prometheus.GaugeValue, flool(alerts.GetThermalThrottle()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishAlertThermalShutdown, prometheus.GaugeValue, flool(alerts.GetThermalShutdown()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishAlertMastNotNearVertical, prometheus.GaugeValue, flool(alerts.GetMastNotNearVertical()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishUnexpectedLocation, prometheus.GaugeValue, flool(alerts.GetUnexpectedLocation()),
	)

	ch <- prometheus.MustNewConstMetric(
		dishSlowEthernetSpeeds, prometheus.GaugeValue, flool(alerts.GetSlowEthernetSpeeds()),
	)

	return true
}

func flool(b bool) float64 {
	if b {
		return 1.00
	}
	return 0.00
}
