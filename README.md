# anmin
照度センサーを使ってRaspberry PiのLEDをON/OFFします。
現在対応しているモジュールは[TSL2561](https://learn.adafruit.com/tsl2561?view=all)のみです。

## Build

```
$ GOOS=linux GOARCH=arm GOARM=5 go build
```

## Usage

デフォルトでは1秒ごとに照度を取得し表示します。また、照度が10未満になるとLEDをOFFにします。

```
# ./anmin
142.53
142.53
11.80
0.22
66.24
141.51
```

### systemdを使ってバックグラウンド起動する

今のところanminにはデーモンモードはありませんが、systemdを使ってバックグラウンド起動することができます。

1. [サンプルのUnitファイル](systemd/anmin.service)を `/etc/systemd/system` 配下に配置する
2. anminを `/usr/local/sbin` 配下に配置する
3. `systemctl daemon-reload` し `systemctl start anmin` します
