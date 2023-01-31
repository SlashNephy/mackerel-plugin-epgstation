# mackerel-plugin-epgstation

A Mackerel plugin to collect EPGStation metrics

## Example Metrics

```console
$ mackerel-plugin-epgstation
epgstation.recording.recording  1       1675135687
epgstation.encode.running       0       1675135687
epgstation.encode.waiting       0       1675135687
epgstation.storages.available   65423953920     1675135687
epgstation.storages.used        3671338823680   1675135687
epgstation.storages.total       3936818806784   1675135687
epgstation.streams.live_stream  0       1675135687
epgstation.streams.live_hls     0       1675135687
epgstation.streams.recorded_stream      0       1675135687
epgstation.streams.recorded_hls 0       1675135687
epgstation.reservation.normal   751     1675135687
epgstation.reservation.skips    0       1675135687
epgstation.reservation.overlaps 82      1675135687
epgstation.reservation.conflicts        0       1675135687
```

## Usage

1. Install via mkr / or download releases directly

```console
$ mkr plugin install SlashNephy/mackerel-plugin-epgstation
```

2. Append following configuration to `mackerel-agent.conf`

```conf
[plugin.metrics.epgstation]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-epgstation"
```

mackerel-plugin-epgstation has some command-line options. Check the help for details.

```console
$ /opt/mackerel-agent/plugins/bin/mackerel-plugin-epgstation --help
```
