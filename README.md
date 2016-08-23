# goping
Reads a list of addresses, ping frequency, and ping timeout from confing json file and logs ping results.

This provides cross platform pings (linux, windows, mac osx).  Under the hood it uses exec.cmd using the native ping command.  The reason for this is that sending ICMP packets requires root privileges.  But the native ping command uses setid or similar to grant permissions so non-root users (or in this case programs) can ping just fine.

## How to use
```
go build goping.go
./goping
# or on windows:  goping.exe

# default config location is ./address_list.json
# but you can override with
./goping -f /path/to/other/config.json
```

## Sample output
```
2016/08/22 22:50:09 event='program_args' config_filename='address_list.json'
2016/08/22 22:50:09 event='config_values' timeout_sec='1' ping_freq_sec='3'
2016/08/22 22:50:10 event='ping_latency' name='TheGoogle' addresss='www.google.com' latency_ms='16'
2016/08/22 22:50:10 event='ping_latency' name='CNN' addresss='www.cnn.com' latency_ms='21'
2016/08/22 22:50:11 event='ping_latency' name='TheGoogle' addresss='www.google.com' latency_ms='24'
2016/08/22 22:50:11 event='ping_latency' name='CNN' addresss='www.cnn.com' latency_ms='30'
2016/08/22 22:50:12 event='ping_latency' name='TheGoogle' addresss='www.google.com' latency_ms='16'
2016/08/22 22:50:12 event='ping_latency' name='CNN' addresss='www.cnn.com' latency_ms='20'
2016/08/22 22:50:13 event='ping_latency' name='CNN' addresss='www.cnn.com' latency_ms='22'
2016/08/22 22:50:13 event='ping_latency' name='TheGoogle' addresss='www.google.com' latency_ms='27'
2016/08/22 22:50:14 event='ping_latency' name='TheGoogle' addresss='www.google.com' latency_ms='14'
```

## How to configure
```
{
  "ping_frequency_sec": 1,
  "ping_timeout_sec": 3,
  "addresses": [
    {
      "name": "TheGoogle",
      "address": "www.google.com"
    },
    {
      "name": "CNN",
      "address": "www.cnn.com"
    }
  ]
}
```
