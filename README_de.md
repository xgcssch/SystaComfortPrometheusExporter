# SystaComfortPrometheusExporter
Empfängt UDP Status Pakete von einer SystaComfort II Heizungssteuerung und stellt diese für die Archivierung in Prometheus zur Verfügung

# Motivation
Als Besitzer einer [Paradigma](https://www.paradigma.de/) Heizung mit einer *SystaComfort II* Heizungssteuerung hat man für die Überwachung nur begrenzte Möglichkeiten:
1. man kann die mobile App [S-Touch](https://www.paradigma.de/produkte/regelungen/s-touch-app/) verwenden, welche jedoch nur minimalistische Grafiken und keine Historiendaten bietet.
2. über ein kostenpflichtiges webbasierendes Remoteportal lassen sich die Daten darstellen
Leider gibt es keine offene, dokumentierte API oder Anbindungen an Überwachungs- und Kontrollframeworks wie MQTT.

Als Folge hat ein Anwender, Klaus Schmidinger, die Kommunikation der *SystaComfort II* Steuerung zum Remoteportal analysiert und ein Perl Skript geschrieben, welches die Daten dekodiert und in RRD Grafiken umwandelt.

Da dieses sich nicht gut mit üblichen Monitoringtools auswerten lässt, habe ich mich entschlossen einen Exporter für [Prometheus](https://prometheus.io/) zu schreiben. Prometheus ist eine Datenbank spezialisiert auf die Erfassung von Zeitreihen. Die dort abgelegten Daten können dann später mit Tools wie [Grafana](https://grafana.com) analysiert und visualisiert werden.

# *SystaComfort* und Modbus TCP
Das Projekt verwendt ein kompliziertes und undokumentiertes Verfahren um an die Monitordaten zu kommen. Ein deutlich besserer Ansatz wäre, auf einem offiziell unterstützen Interface aufzusetzen. Glücklicherweise steht dieser Weg allen Anwendern mit einer Hardware Version 2 und höher - Herstellungsdatum ca. 2016 - zur Verfügung. Die Version lässt sich leicht an der Platine erkennen: die Version hat nur einen *LAN* Port, während die Version 2 aufweist. Wenn sie die Version 2 besitzen, dann lässt sich die Steuerung komplett über den *Modbus TCP* Controller abfragen und steuern. 

Weitere Informationen über das *SystaComfort* *Modbus TCP* Interface findet sich [hier](https://mam-prod.paradigma.de/pinaccess/pinaccess.do?pinCode=6om1uCSh0OKd). Um das Interface zu nutzen ist ggf. eine Aktualisierung der Software auf eine akutelle Version notwendig.

Dank an Boris Bartenstein von *Paradigma* für diese Information.
# Voraussetzungen
Bevor die Installation durchgeführt werden kann, sollten folgende Punkte sichergestellt sein:
1. Es muss eine *SystaComfort II* Steuerung vorhanden sein, welche Daten unverschlüsselt Daten an das Remoteportal von Paradigma sendet.
1. Es muss ein Router zur Verfügung stehen, mit dem ein existierender DNS Eintrag übersteuert werden kann
1. Ein Server muss zur Verfügung stehen, auf welchem permanent dieses Programm läuft um die Monitoringwerte im Prometheus Format zur Verfügung zu stellen.
1. Weiterhin wird ein Server benötigt, auf welchem Prometheus läuft, um die Monitoringwerte zu erfassen und archivieren.
1. Schlussendlich benötigt man noch eine Grafana Instanz um die arichivierten Prometheusdaten zu visualisieren.

Alle diese Komponenten können auf kleinen und sparsamen Geräten, wie etwa einem [Rasperry Pi](https://www.raspberrypi.org/) oder einem NAS wie der Synology DiskStation laufen, gerade wenn diese
auch den Betrieb von Docker Containern erlauben.
## Getestete Versionen
Die folgenden Versionen des *SystaComfort* Controllers konnte erfolgreich mit den folgenden Hardwar/Software Kombinationen betrieben werden:
- HW 1 - V1.24
- HW 1 - V1.26
- HW 2 - V2.14
# Installation
## Konfiguration des  DNS
Der DNS Name `paradigma.remoteportal.de` muss auf den Server zeigen, auf dem der Exporter läuft. In diesem Beispiel ist das die Adresse `192.168.100.4`.

Jeder Internet Router hat sein eigenes Verfahren um einen DNS Eintrag zu übersteuern. Manchen Modellen fehlt auch die Möglichkeit, DNS Adressen zu übersteuern.
Wenn Sie einen solchen Router besitzen, ist es schwierig dieses Problem zu umgehen. Eine Möglichkeit ist, einen eigenständigen DNS Server im lokalen Netzwerk zu installieren. Insgesamt sind solche 
Lösungen nicht trivial zu implementieren und erfordern erweiterte Kenntnisse über die Netzwerkonfiguration.

Wenn Sie aber in der Lage sind, dem Namen eine interne Adresse zuzuordnen, dann erhalten Sie bei der Namensauflösung eine Ausgabe ähnlich dieser:
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
Haben Sie einen Windows Rechner, so verwenden Sie `nslookup paradigma.remoteportal.de` zur Abfrage.
## Konfiguration der  *SystaComfort II* Heizungssteuerung
Es gibt ein *Serviceprogramm Paradigma*, welches vom Heizungsbauer zur Konfiguration verwendet wird. Unterhalb des Menüeintrages *Einstellungen* findet sich ein Eintrag *Netzwerkeinstellungen SystaComfort II*. Wählt man diesen aus, so werden die Netzwerkeinstellungen angezeigt. In der unteren Hälfte befindet sich ein Bereich mit zwei Reitern, der zweite ist mit *Portal* bezeichnet. 
Diesen wählt man aus. Dort wiederum befinden sich zwei Checkboxen: *Remoteportal aktiv* und *Verschlüsselung aktiv*.
>Die erste Option muss dabei ausgewählt sein, die Zweite nicht, da wir keine 
Möglichkeit haben die Daten zu entschlüsseln!

![Screenshot der Netzwerkkonfiguration](https://github.com/xgcssch/SystaComfortPrometheusExporter/raw/main/doc/assets/NetworkSettings.png)

# Exporter als Docker Container
Die einfachste Methode den Exporter zu installieren ist, wenn Sie bereits ein Gerät haben auf welchem eine Containerumgebung wie Docker läuft.

In diesem Fall reicht es den Container bspw. mit
```Shell
docker run -d --name systacomfortexporter \
    -p 22460:22460/udp \
    -p 2112:2112/tcp \
    xgcssch/systacomfortprometheusexporter
```
zu starten. Dieser veröffentlich dann den UDP Endpunkt für die Heizungssteuerung und den TCP Endpunkt für Prometheus.

Achtung: die Adresse der Maschine wo der Exporter läuft muss diejenige Adresse sein, auf die `paradigma.remoteportal.de` zeigt!
## Exporter als Programm auf einem Server
Die aktuelle Fassung des Exporters kann aus der *Releases+ section des Projektes für diverse Systeme heruntergeladen werden.

Danach einfach das Programm starten:
```
luna $ ./SystaComfortPrometheusExporter
2020/10/25 11:23:14 Starting ...
```
Um zu überprüfen, ob das Programm korrekt läuft und Werte erfasst, können die Daten mit einem Browser oder über ein Kommandozeilentool angezeigt werden:
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
Um das Programm als Windows Service laufen zu lassen, kann man [NSSM](https://nssm.cc/) verwenden.
# Prometheus Installation
Die Installationsanweisungen finden sich auf der [Prometheus](https://prometheus.io/docs/prometheus/latest/installation/) Website.

Ich empfehle, Prometheus auch als Container zu starten. Dafür benötigt man eine Konfigruationsdatei, welche auf die Maschine zeigt wo der
Exporter installiert ist. Diese sieht dann in etwa so aus:
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
Achten Sie darauf, dass das *data* Verzeichnis auf persistenten Speicher liegt, andernfalls gehen die erfassten Daten zusammen mit dem Container verloren!

Beispiel:
```
docker run -d --name prometheus \
    -p 9090:9090 \
    -v /var/data/prometheus/data:/prometheus \
    -v /var/data/prometheus/config/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus
```
Wenn alles korrekt läuft, kann man die Oberfläche von Prometheus auf `http://192.168.100.4:9090` finden. Unterhalb des Menüeintrags *Status* -> *Targets* sollte sich eine Liste der abgefragten Endpunkte befinden. Der Endpunkt mit der Bezeichnung *systaComfort* sollte als 'State=UP' dargestell sein.
# Grafana Installation
Installationsanweisungen finden sich auf der [Grafana](https://grafana.com/docs/grafana/latest/) Website.
Und auch hier empfehle ich die Installation als Container. 

Wie mit dem Prometheus Container, benötigen Sie auch hier einen persistenten Speicher um die Einstellungen dauerhaft zu behalten.
```
docker run -d --name grafana \
    -p 3000:3000 \
    -v /var/data/grafana:/var/lib/grafana \
    grafana/grafana
```
Um *Grafana* weiter zu Konfigurieren führen Sie diese Schritte aus:
1. Anmelden bei Grafana auf `http://192.168.100.4:3000`. Initialer Username und Kennwort sind `admin`.
1. Nach `Configuration` -> `Datasources` gehen, dort den `Add Datasource` Knopf wählen. `Prometheus` auswählen und `http://192.168.100.4:9090` als Server URL eintrgen. Danach `Test and Save` drücken.
1. Das [Dashboard Template](https://github.com/xgcssch/SystaComfortPrometheusExporter/blob/main/assets/grafana/SystaComfort-Dashboard.json) von Github herunterladen.
1. Nach `Create`-> `Import` gehen und dort die heruntergeladene Datei zum Import auswählen. `Prometheus` als Datenquelle angeben.
1. *Grafana* sollte danach bereits das Dashboard mit den aktuellen Daten anzeigen.

![Screenshot des Grafana Dashboards](https://github.com/xgcssch/SystaComfortPrometheusExporter/raw/main/doc/assets/GrafanaDashboard.png)
# Fehlersuche
1. Sicherstellen, dass die DNS Auflösung von der *SystaComfort II* aus korrekt funktionert.

   Dafür muss man im *Serviceprogramm Paradigma* unter `Optionen`-> `Netzwerkeinstellungen SystaComfort II` die aktuellen Netzwerkeinstellungen anschauen:

   ![Screenshot der Basis-Netzwerkkonfiguration](https://github.com/xgcssch/SystaComfortPrometheusExporter/raw/main/doc/assets/NetworkSettings-DNS.png)

   Der dort eingetragene DNS-Server muss die Adresse `paradigma.remoteportal.de` desjenigen Servers zurückliefern, auf welchem der Exporter läuft.
   
   In dem folgenden Beispiel läuft der Exporter auf dem Rechner mit der IP `192.168.100.4` und der DNS Server ist wie in der Abbildung ersichtlich die `192.168.100.3`

   >Ist der DNS Server nicht der gewünschte Server => diesen korrigieren!

   Beispiel für Windows:
   ```Batchfile
   C:\>nslookup paradigma.remoteportal.de 192.168.100.3
   Server:  mikrotik-main.int.schau.org
   Address:  192.168.100.3
   
   Nicht autorisierende Antwort:
   Name:    paradigma.remoteportal.de
   Address:  192.168.100.4
   ```
   Beispiel für Unix:
   ```Shell
   [root@minerva ~]# host paradigma.remoteportal.de 192.168.100.3
   Using domain server:
   Name: 192.168.100.3
   Address: 192.168.100.3#53
   Aliases:
   
   paradigma.remoteportal.de has address 192.168.100.4
   ```

   >Liefert die Abfrage nicht das gewünschte Ergebnis: DNS Konfiguration prüfen!
1. Wenn ein Dockercontainer verwendet wird: überprüfen, ob das Portmapping korrekt spezifiziert wurde.
   
   Dafür auf dem Docker Host folgendes Kommando absetzen:

   ```Shell
   soenke@nas01:/$ docker inspect --format='{{range $p, $conf := .NetworkSettings.Ports}} {{$p}} -> {{(index $conf 0).HostPort}} {{end}}' systacomfortexporter
   2112/tcp -> 2112  22460/udp -> 22460
   ```
   Das Ergebnis sollte 
   
   a) `22460/udp -> 22460` enhalten, damit die UDP Pakete von der Heizung an den Container weitergeleitet werden und 
   
   b) `2112/tcp -> 2112` damit Anforderungen an den Prometheus Exporter weitergeleitet werden.

   >Sind die Portmappings nicht vorhanden => `docker run ...` Kommando überprüfen!

1. Sicherstellen, dass UDP Daten von der Heizung auf dem (Docker-)Host ankommen.
   
   Falls der (Docker-)Host ein *ix basierendes System ist, kann `tcpdump`verwendet werden:
   ```Shell
   bash-4.3# tcpdump udp port 22460
   tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
   listening on eth0, link-type EN10MB (Ethernet), capture size 262144 bytes
   20:43:14.579978 IP 192.168.100.250.8002 > nas01.int.schau.org.22460: UDP, length 1048
   20:43:14.580366 IP nas01.int.schau.org.22460 > 192.168.100.250.8002: UDP, length 20
   20:43:14.620267 IP 192.168.100.250.8002 > nas01.int.schau.org.22460: UDP, length 1048
   20:43:14.620644 IP nas01.int.schau.org.22460 > 192.168.100.250.8002: UDP, length 20
   20:43:14.660629 IP 192.168.100.250.8002 > nas01.int.schau.org.22460: UDP, length 1048
   20:43:14.660894 IP nas01.int.schau.org.22460 > 192.168.100.250.8002: UDP, length 20
   20:43:14.700913 IP 192.168.100.250.8002 > foresnas01.int.schau.org.22460: UDP, length 192
   ^C
   7 packets captured
   8 packets received by filter
   0 packets dropped by kernel
   ```

   >Werden keine Pakete angezeigt => Netzwerkkonfiguration zwischen Heizung und (Docker-)Host prüfen!

1. Sicherstellen, dass UDP Daten von der Heizung in dem Dockercontainer ankommen.
   
   Dafür müssen wir eine Shell in den Container öffnen, `tcpdump` installieren und ausführen:

   ```Shell
   bash-4.3# docker exec -it systacomfortexporter ash
   ~ # tcpdump
   ash: tcpdump: not found
   ~ # apk update
   fetch http://dl-cdn.alpinelinux.org/alpine/v3.12/main/x86_64/APKINDEX.tar.gz
   fetch http://dl-cdn.alpinelinux.org/alpine/v3.12/community/x86_64/APKINDEX.tar.gz
   v3.12.3-76-ged1baecfad [http://dl-cdn.alpinelinux.org/alpine/v3.12/main]
   v3.12.3-74-g09e375413f [http://dl-cdn.alpinelinux.org/alpine/v3.12/community]
   OK: 12747 distinct packages available
   ~ # apk add tcpdump
   (1/2) Installing libpcap (1.9.1-r2)
   (2/2) Installing tcpdump (4.9.3-r2)
   Executing busybox-1.31.1-r19.trigger
   OK: 7 MiB in 17 packages
   ~ # tcpdump udp port 22460
   tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
   listening on eth0, link-type EN10MB (Ethernet), capture size 262144 bytes
   19:55:08.157164 IP 172.17.0.1.35359 > 8115f53e1933.22460: UDP, length 1048
   19:55:08.157594 IP 8115f53e1933.22460 > 172.17.0.1.35359: UDP, length 20
   19:55:08.197455 IP 172.17.0.1.35359 > 8115f53e1933.22460: UDP, length 1048
   19:55:08.197663 IP 8115f53e1933.22460 > 172.17.0.1.35359: UDP, length 20
   19:55:08.237742 IP 172.17.0.1.35359 > 8115f53e1933.22460: UDP, length 1048
   19:55:08.237892 IP 8115f53e1933.22460 > 172.17.0.1.35359: UDP, length 20
   19:55:08.278079 IP 172.17.0.1.35359 > 8115f53e1933.22460: UDP, length 192
   ^C
   7 packets captured
   7 packets received by filter
   0 packets dropped by kernel
   ~ exit
   ```
   >Werden keine Pakete angezeigt => Netzwerkkonfiguration des Dockerhostes in Zusammenhang mit dem Container prüfen!
1. Sicherstellen, dass die exportierten Daten des Containers abgefragt werden können.
   
   Dazu mit Hilfe von `curl` den Endpunkt abfragen:

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
   >Werden keine Werte oder eine Fehlermeldung angezeigt => Netzwerkkonfiguration des Containers prüfen!

1. Sicherstellen, dass die Heizungswerte auch exportiert werden.
   
   Sind bei dem vorigen Schritt Daten zu sehen gewesen, welche ungleich `0` waren, dann kann man davon ausgehen, dass die Daten vom Exporter empfangen und verarbeitet wurden.

   Bspw.

   ```
   ...
   # HELP systacomfort_heatercircuit_return_temperature_celsius The boiler return temperature
   # TYPE systacomfort_heatercircuit_return_temperature_celsius gauge
   systacomfort_heatercircuit_return_temperature_celsius 27.4
   # HELP systacomfort_heatercircuit_supply_temperature_celsius The boiler supply temperature
   # TYPE systacomfort_heatercircuit_supply_temperature_celsius gauge
   systacomfort_heatercircuit_supply_temperature_celsius 30.6
   ...
   ```

   Ist dieses nicht der Fall, so kann man das Progamm auch mit erweiterten Ausgaben starten. Am Besten startet man den Container auch interaktiv, damit man dessen Ausgaben direkt sehen kann:

   ```Shell
   docker run -it -p 22460:22460/udp -p 2112:2112/tcp xgcssch/systacomfortprometheusexporter /root/SystaComfortPrometheusExporter -v 4
   ```
   
   Relevant ist hier das `-v` Flag. Eine Stufe von `4` gibt das Programm bei jedem UDP Paket einen Hinweis aus. Bei Stufe `5` dumpt er zusätzlich die Werte.

   > Bitte daran denken, dass nur ein Exportercontainer gleichzeitig laufen kann. Ggf. ist ein anderer Container vorher zu beenden.

# FAQ
1. Du hast eine Frage oder einen Fehler gefunden? Möchtest du Verbesserungsvorschläge machen?

   Dann beginne eine Diskussion oder lege einen Issue an. Gerne auch in Deutscher Sprache. Ich habe die Hauptseite zwar in Englisch erstellt, jedoch sind bisher alle mir bekannten Anwender aus dem Deutschen Sprachbereich. Sollte sich das als Fehlannahme erweisen, so können wir dieses dann ändern ;-) 

1. Wo bekomme ich das *Serviceprogramm Paradigma* her?

   ~~Das Programm kann man bei *Paradigma* direkt herunterladen: https://www.paradigma.de/software/. Es handelt sich dabei um das Programm `SystaService` in der *PC-Software* Kategorie. Ich habe dieses in der Version 1.60 verwendet.~~
   
   **Update:** Die Webseite ist nur noch für angemeldete Fachbetriebe zugänglich. Allerdings kann man die Software weiterhin über folgendes Downloadportal herunterladen:
   
   https://mam-downloadcenter.paradigma.de/page/faq  
   Suchbegriff: "SystaService CD", bzw [hier](https://mam-downloadcenter.paradigma.de/page/download?assetid=40688&assetname=SystaService%20CD_V160.zip) der Direktdownload.
   
   (Hinweis: Firefox lädt die Datei fälschlicherweise als PDF herunter. Falls das passiert einfach die Dateiendung nach dem Download zu `.zip` ändern)
   
   
1. Mit welcher *SystaComfort II* Version läuft der Exporter?

   Ich kann im Augenblick nur von meiner Anlage definitiv sagen, dass es läuft:
   
   ![Screenshot der Systemversion](https://github.com/xgcssch/SystaComfortPrometheusExporter/raw/main/doc/assets/Systemversion.jpg)

   Die Version wurde letztes Jahr von meinem Heizungsbauer aktualisiert. Ich kann nicht sagen welche Version vorher aktiv war, nur dass diese ebenfalls korrekt Daten geliefert hat.
   
1. Ich habe einen Router, mit dem ich keine DNS Namen anpassen kann (bspw. Fritzbox)

   Eine Möglichkeit besteht darin, auf dem Dockerhost zustätzlich einen Container mit einem DNS Server zu fahren und den DNS Server in der *SystaComfort II* Steuerung explizit auf diesen zu verweisen.
   
   Bspw. 

   ```Shell
   docker run -p 53:53/tcp -p 53:53/udp --cap-add=NET_ADMIN andyshinn/dnsmasq:2.75 -S 192.168.100.1 --address=/paradigma.remoteportal.de/192.168.100.4
   ```

   wobei `192.168.100.4` die IP des Dockerhostes ist, auf dem der Exporter Container läuft und `192.168.100.1` die IP des Routers bzw. des lokalen DNS Servers.

1. Ich hätte gerne noch den Wert xyz exportiert?

   Prima! Bekommen wir hin, wenn du weißt wo in den übermittelten Daten sich der relevante Wert befindet. Mach einen Issue auf, oder noch besser: erstelle einen Pull Request mit dem entsprechenden Code.


