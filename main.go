package main

import (
	"fmt"
	"log"

	"github.com/jessevdk/go-flags"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

type commandLine struct {
	EPGStationHost string `long:"host" description:"EPGStation host" default:"localhost"`
	EPGStationPort int    `long:"port" description:"EPGStation port" default:"8888"`
	Prefix         string `long:"prefix" description:"Metric key prefix" default:"epgstation"`
	Tempfile       string `long:"tempfile" description:"Temp filename"`
}

type plugin struct {
	api     *EPGStationAPI
	options *commandLine
}

func (u plugin) MetricKeyPrefix() string {
	return u.options.Prefix
}

const (
	keyStreamsLiveStream     = "live_stream"
	keyStreamsLiveHLS        = "live_hls"
	keyStreamsRecordedStream = "recorded_stream"
	keyStreamsRecordedHLS    = "recorded_hls"
	keyReserveNormal         = "normal"
	keyReserveSkips          = "skips"
	keyReserveOverlaps       = "overlaps"
	keyReserveConflicts      = "conflicts"
	keyRecordingCount        = "recording"
	keyEncodeRunning         = "running"
	keyEncodeWaiting         = "waiting"
	keyStoragesAvailable     = "available"
	keyStoragesUsed          = "used"
	keyStoragesTotal         = "total"
)

func (u plugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"streams": {
			Label: "EPGStation Streams",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{
					Name:    keyStreamsLiveStream,
					Label:   "Live Stream",
					Stacked: true,
				},
				{
					Name:    keyStreamsLiveHLS,
					Label:   "Live HLS",
					Stacked: true,
				},
				{
					Name:    keyStreamsRecordedStream,
					Label:   "Recorded Stream",
					Stacked: true,
				},
				{
					Name:    keyStreamsRecordedHLS,
					Label:   "Recorded HLS",
					Stacked: true,
				},
			},
		},
		"reservation": {
			Label: "EPGStation Reservation",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{
					Name:    keyReserveNormal,
					Label:   "Normal",
					Stacked: true,
				},
				{
					Name:    keyReserveSkips,
					Label:   "Skips",
					Stacked: true,
				},
				{
					Name:    keyReserveOverlaps,
					Label:   "Overlaps",
					Stacked: true,
				},
				{
					Name:    keyReserveConflicts,
					Label:   "Conflicts",
					Stacked: true,
				},
			},
		},
		"recording": {
			Label: "EPGStation Recording",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{
					Name:  keyRecordingCount,
					Label: "Count",
				},
			},
		},
		"encode": {
			Label: "EPGStation Encoding",
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				{
					Name:    keyEncodeRunning,
					Label:   "Running",
					Stacked: true,
				},
				{
					Name:    keyEncodeWaiting,
					Label:   "Waiting",
					Stacked: true,
				},
			},
		},
		"storages": {
			Label: "EPGStation Storages",
			Unit:  mp.UnitBytes,
			Metrics: []mp.Metrics{
				{
					Name:  keyStoragesTotal,
					Label: "Total",
				},
				{
					Name:    keyStoragesAvailable,
					Label:   "Available",
					Stacked: true,
				},
				{
					Name:    keyStoragesUsed,
					Label:   "Used",
					Stacked: true,
				},
			},
		},
	}
}

func (u plugin) FetchMetrics() (map[string]float64, error) {
	metrics := map[string]float64{}

	if err := u.appendStreamsMetrics(metrics); err != nil {
		return nil, err
	}

	if err := u.appendReserveCountsMetrics(metrics); err != nil {
		return nil, err
	}

	if err := u.appendRecordingMetrics(metrics); err != nil {
		return nil, err
	}

	if err := u.appendEncodeMetrics(metrics); err != nil {
		return nil, err
	}

	if err := u.appendStoragesMetrics(metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (u plugin) appendStreamsMetrics(metrics map[string]float64) error {
	streams, err := u.api.GetStreams()
	if err != nil {
		return err
	}

	if streams.Code != 0 {
		return fmt.Errorf("failed to get streams: %d: %s, %s", streams.Code, streams.Message, streams.Errors)
	}

	var (
		ls float64
		lh float64
		rs float64
		rh float64
	)
	for _, stream := range streams.Items {
		switch stream.Type {
		case "LiveStream":
			ls++
		case "LiveHLS":
			lh++
		case "RecordedStream":
			rs++
		case "RecordedHLS":
			rh++
		}
	}

	metrics[keyStreamsLiveStream] = ls
	metrics[keyStreamsLiveHLS] = lh
	metrics[keyStreamsRecordedStream] = rs
	metrics[keyStreamsRecordedHLS] = rh
	return nil
}

func (u plugin) appendReserveCountsMetrics(metrics map[string]float64) error {
	counts, err := u.api.GetReserveCounts()
	if err != nil {
		return err
	}

	if counts.Code != 0 {
		return fmt.Errorf("failed to get reserve counts: %d: %s, %s", counts.Code, counts.Message, counts.Errors)
	}

	metrics[keyReserveNormal] = float64(counts.Normal)
	metrics[keyReserveSkips] = float64(counts.Skips)
	metrics[keyReserveOverlaps] = float64(counts.Overlaps)
	metrics[keyReserveConflicts] = float64(counts.Conflicts)
	return nil
}

func (u plugin) appendRecordingMetrics(metrics map[string]float64) error {
	recording, err := u.api.GetRecording()
	if err != nil {
		return err
	}

	if recording.Code != 0 {
		return fmt.Errorf("failed to get recording: %d: %s, %s", recording.Code, recording.Message, recording.Errors)
	}

	metrics[keyRecordingCount] = float64(len(recording.Records))
	return nil
}

func (u plugin) appendEncodeMetrics(metrics map[string]float64) error {
	encode, err := u.api.GetEncode()
	if err != nil {
		return err
	}

	if encode.Code != 0 {
		return fmt.Errorf("failed to get encode: %d: %s, %s", encode.Code, encode.Message, encode.Errors)
	}

	metrics[keyEncodeRunning] = float64(len(encode.RunningItems))
	metrics[keyEncodeWaiting] = float64(len(encode.WaitItems))
	return nil
}

func (u plugin) appendStoragesMetrics(metrics map[string]float64) error {
	storages, err := u.api.GetStorages()
	if err != nil {
		return err
	}

	if storages.Code != 0 {
		return fmt.Errorf("failed to get storages: %d: %s, %s", storages.Code, storages.Message, storages.Errors)
	}

	var (
		available int
		used      int
		total     int
	)
	for _, storage := range storages.Items {
		available += storage.Available
		used += storage.Used
		total += storage.Total
	}

	metrics[keyStoragesAvailable] = float64(available)
	metrics[keyStoragesUsed] = float64(used)
	metrics[keyStoragesTotal] = float64(total)
	return nil
}

func main() {
	options := commandLine{}
	parser := flags.NewParser(&options, flags.Default)

	_, err := parser.Parse()
	if err != nil {
		log.Fatalln(err)
	}

	helper := mp.NewMackerelPlugin(plugin{
		api:     NewEPGStationAPI(options.EPGStationHost, options.EPGStationPort),
		options: &options,
	})
	helper.Tempfile = options.Tempfile
	helper.Run()
}
