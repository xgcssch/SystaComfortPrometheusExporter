# SystaComfortPrometheusExporter
Intercepts UDP Status Packets from a SystaComfort Heating controller and makes them available to Prometheus

# Motivation
If you are a owner of a *SystaComfort* Heating from [Paradigma](https://www.paradigma.de/), then you have limited possibilities to monitor the
device: You can use a mobile app called [S-Touch](https://www.paradigma.de/produkte/regelungen/s-touch-app/) which just exposes the current values or minimalistic graphics or an webbased remote portal. There are neither Open APIs nor interfaces to common monitoring or control tools like MQTT.

As a result, Klaus Schmidinger decoded the protocol which every *SystaComfort* Controler sends to the Remote portal. He developed a perl script which decoded the data and produced RRD graphs.

As time goes by and technology advances further, i decided to rewrite the perl script and build an exporter for [Prometheus](https://prometheus.io/), a time-series Database specialized in storing metrics. The data contained in Prometheus can later be viewed and analyzed in tools like [Grafana](https://grafana.com).

# Requirements
Your need some things:
1. A configurable Router where you can change the IP Address of an official DNS entry to a local IP in your network
1. Some kind of server, where you can permanently run this program, which converts the *SystaCoomfort* data into the format used by Prometheus
1. A permanently running Prometheus instance.
1. A Grafana instance to get the data visualized

All of this, except the router configuration, can be run on small devices like a [Rasperry Pi](https://www.raspberrypi.org/) or a NAS Device like Synology DiskStation.

# Installation
## As a docker container
## As a program running on a server machine
# DNS Configuration
The DNS Name `paradigma.remoteportal.de` must point to the device where the exporter is running. In this sample it is the `192.168.100.4`.

```
luna $ dig paradigma.remoteportal.de

; <<>> DiG 9.11.3-1ubuntu1.13-Ubuntu <<>> paradigma.remoteportal.de
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 55342
;; flags: qr rd ad; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;paradigma.remoteportal.de.     IN      A

;; ANSWER SECTION:
paradigma.remoteportal.de. 0    IN      A       192.168.100.4

;; Query time: 3 msec
;; SERVER: 172.24.176.1#53(172.24.176.1)
;; WHEN: Sun Oct 25 00:34:00 CEST 2020
;; MSG SIZE  rcvd: 84
```
# Prometheus Installation
# Grafana Installation
