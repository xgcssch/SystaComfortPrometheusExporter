# SystaComfortPrometheusExporter
Intercepts UDP Status Packets from a SystaComfort Heating controller and makes them available to Prometheus

# Motivation
If you are an owner of a *SystaComfort* Heating from [Paradigma](https://www.paradigma.de/), then you have limited possibilities to monitor the
device: You can use a mobile app called [S-Touch](https://www.paradigma.de/produkte/regelungen/s-touch-app/) which just exposes the current values or minimalistic graphics or an webbased remote portal. There are neither Open APIs nor interfaces to common monitoring or control tools like MQTT.

As a result, Klaus Schmidinger decoded the protocol which every *SystaComfort* Controler sends to the Remote portal. He developed a perl script which decoded the data and produced RRD graphs.

As time goes by and technology advances further, i decided to rewrite the perl script and build an exporter for [Prometheus](https://prometheus.io/), a time-series Database specialized in storing metrics. The data contained in Prometheus can later be viewed and analyzed in tools like [Grafana](https://grafana.com).

# Requirements
Before you start, be sure you need some things:
1. A *SystaComfort* Controller configured to send data unencrypted to the Paradigma Remote Portal.
1. A configurable Router where you can change the IP Address of an official DNS entry to a local IP in your network
1. Some kind of server, where you can run this program, which converts the *SystaCoomfort* data into the format used by Prometheus. This program has to run continuously to provide the data. 
1. A permanently running Prometheus instance.
1. A Grafana instance to get the data visualized

All of these software components can be run on small devices like a [Rasperry Pi](https://www.raspberrypi.org/) or a NAS Device like Synology DiskStation capable of running Docker Container.

# Installation
## Configure your DNS
The DNS Name `paradigma.remoteportal.de` must point to the device where the exporter is running. In this sample this is the Address `192.168.100.4`.

Every internet router has it's own method to set such an entry. Some models are even incapble of overriding DNS addresses. If you own such device, there is no easy way to overcome the problem. One way is to install a local DNS server in your network, override the DNS Name there and only point the SystaComfort to that specific DNS Server. Unfortunatly this is not an easy thing to do and falls outside the scope of this document.
See here for doing this with an [Rasperry Pi](https://www.deviceplus.com/raspberry-pi/how-to-use-a-raspberry-pi-as-a-dns-server/) or with a [Container](https://github.com/opsxcq/docker-dnsmasq)

But if you were able to spoof the address, you should see an output like this if you query the address:
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
If you are on a windows machine, use `nslookup paradigma.remoteportal.de` to check.
## Configure *SystaComfort* Controller
There is a *Serviceprogramm Paradigma* which is normally used by maintainance personal. Below the menu group entry *Options* there is an entry called *Networksettings SystaComfort II*. When selected, it displays the Network settings of the controller. In the lower half there is a property sheet with two pages: the second is labeled *Portal*. Select this page. There are two checkboxes: the first, *Remoteportal aktiv* must be checked to enable the data transmissions. The second, *Verschlüsselung aktiv* should not be checked, as we have no information how to decrypt the data.

![Screenshot of the network configuration](https://github.com/xgcssch/SystaComfortPrometheusExporter/raw/main/doc/assets/NetworkSettings.png)

# As a docker container
The easiest way to install the exporter if you already have a device with *Docker* running, is to use a container.

Just issue
```
docker run -d --name systacomfortexporter \
    -p 22460:22460/udp \
    -p 2112:2112/tcp \
    xgcssch/systacomfortprometheusexporter
```
this will run the exporter and expose both the UDP endpoint for the heating controller and the TCP endpoint for Prometheus.

Remember that the address of this machine is the one where `paradigma.remoteportal.de` should point to!
## As a program running on a server machine
Download the binary for your platform from the *Releases+ section of this project.

Just run the executable:
```
luna $ ./SystaComfortPrometheusExporter
2020/10/25 11:23:14 Starting ...
```
Accessing the metrics endpoint with curl should show the current metrics:
```
luna $ curl http://192.168.100.4:2112/metrics
# HELP systacomfort_boiler_active_info Is the boiler is running
# TYPE systacomfort_boiler_active_info gauge
systacomfort_boiler_active_info 0
# HELP systacomfort_boiler_circulationpump_info Is the boiler circulation pump running
# TYPE systacomfort_boiler_circulationpump_info gauge
systacomfort_boiler_circulationpump_info 1
...
```

### Windows Service
To run the program as a windows service, use [NSSM](https://nssm.cc/)
# Prometheus Installation
Installation instructions can be found at the [Prometheus](https://prometheus.io/docs/prometheus/latest/installation/) website.

I would recommend using a docker container too. You will need a configuration file which points to the machine running the exporter like this:
```
global:
  evaluation_interval: 1m
  scrape_interval: 1m
  scrape_timeout: 10s
scrape_configs:
- job_name: systaComfort
  static_configs:
  - targets:
    - 192.168.100.4:2112
```
Make sure you put the data directory on persistent storage, otherwise the recorded data is lost when the container is removed!

Example:
```
docker run -d --name prometheus \
    -p 9090:9090 \
    -v /var/data/prometheus/data:/prometheus \
    -v /var/data/prometheus/config/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus
```
If everything works, you can find the Prometheus UI at `http://192.168.100.4:9090`. Below the menu entry *Status* -> *Targets* there should be a list  of scraped hosts. The host running the SystaComfortExporter should show up as 'State=UP*.
# Grafana Installation
Installation instructions can be found at the [Grafana](https://grafana.com/docs/grafana/latest/) website.
Again i would recommend using a docker container. 

As with the prometheus container, you will need a volume mount to make your configuration persistent.
```
docker run -d --name grafana \
    -p 3000:3000 \
    -v /var/data/grafana:/var/lib/grafana \
    grafana/grafana
```
To further configure *Grafana* follow these steps:
1. Login into Grafana at `http://192.168.100.4:3000`. Initial username and password are `admin`.
1. Goto `Configuration` -> `Datasources`, and press the `Add Datasource` button. Select `Prometheus` and enter `http://192.168.100.4:9090` as the URL. Press `Test and Save`.
1. Download the [dashboard template](https://github.com/xgcssch/SystaComfortPrometheusExporter/blob/main/assets/grafana/SystaComfort-Dashboard.json) from Github.
1. Continue to `Create`-> `Import` and select the downloaded file for import. Choose `Prometheus` as the datasource.
1. *Grafana* should display the data.

![Screenshot of the Grafana dashboard](https://github.com/xgcssch/SystaComfortPrometheusExporter/raw/main/doc/assets/GrafanaDashboard.png)
