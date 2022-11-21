# Linkchecker

## Download

Download [here](https://github.com/tlanfer/linkchecker/releases/tag/latest).

## Usage

Start the executable. Windows may complain. Tell windows to shut up.

The program will create a default config with sensible defaults. You can change it if you want.

```yaml
prefix: monitor_
interval: 15s
threshold:
    packet_loss: 0.01
    rtt: 100ms
hosts:
    - youtube.com
    - twitch.tv 
```

By default, it only logs failures. If you want to log everything, change the thresholds in the config to zero, like this:
```yaml
threshold:
    packet_loss: 0
    rtt: 0ms
```

It will generate a log file per configured host. Something like this:

```csv
2022-11-21 22:39:10,regjeringen.no,100,0
2022-11-21 22:39:25,regjeringen.no,100,0
2022-11-21 22:39:40,regjeringen.no,100,0
2022-11-21 22:39:55,regjeringen.no,100,0
2022-11-21 22:40:10,regjeringen.no,100,0
```

The columns are
* Timestamp
* host
* Packageloss percentage
* Round-trip-time ("ping")