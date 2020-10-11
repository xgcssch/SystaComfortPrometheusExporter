//
// Derived from a perl script developed by Klaus.Schmidinger@tvdr.de

package udpserve

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Außentemperatur Fühler
	systacomfortOutsideTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_outside_temperature_celsius",
		Help: "The outside temperature",
	})
	// Heizung Vorlauf
	systacomfortHeatercircuitSupplyTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatercircuit_supply_temperature_celsius",
		Help: "The boiler supply temperature",
	})
	// Heizung Ruecklauf
	systacomfortHeatercircuitReturnTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatercircuit_return_temperature_celsius",
		Help: "The boiler return temperature",
	})
	// Brauchwasser
	systacomfortTapwaterTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_tapwater_temperature_celsius",
		Help: "The tapwater temperature",
	})
	// Zirkulation
	systacomfortCirculationTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_watercirculation_temperature_celsius",
		Help: "The temperature in the water circulation",
	})
	// Raumtemperatur
	systacomfortInsideTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_temperature_celsius",
		Help: "The inside temperature",
	})
	// Kollektor
	// KesselVorlauf
	systacomfortBoilerSupplyTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_supply_temperature_celsius",
		Help: "The boiler supply temperature",
	})
	// KesselRuecklauf
	systacomfortBoilerReturnTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_return_temperature_celsius",
		Help: "The boiler return temperature",
	})
	// BrauchwasserSoll
	systacomfortTapwaterTargetTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_tapwater_target_temperature_celsius",
		Help: "The tapwater target temperature",
	})
	// InnenSoll
	systacomfortInsideTargetTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_target_temperature_celsius",
		Help: "The inside target temperature",
	})
	// KesselSoll
	systacomfortBoilerTargetTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_target_temperature_celsius",
		Help: "The boiler target temperature",
	})
	// Betriebsart 0=Heizprogramm 1, 1=Heizprogramm 2, 2=Heizprogramm 3, 3=Dauernd Normal, 4=Dauernd Komfort, 5=Dauernd Absenken
	systacomfortModeInfo = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_mode_info",
		Help: "The current mode",
	})
	// 39=Raumtemperatur normal (soll)
	// 40=Raumtemperatur komfort (soll)
	// 41=Raumtemperatur abgesenkt (soll)
	// Fusspunkt
	systacomfortHeatCurveRootPointCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatcurve_rootpoint_celsius",
		Help: "The rootpoint of the heating curve",
	})
	// Steilheit
	systacomfortHeatCurveGradientCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatcurve_gradient_celsius",
		Help: "The gradient of the heating curve",
	})
	// Max. Vorlauftemperatur
	systacomfortBoilerSupplyUpperLimitCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_supply_upperlimit_celsius",
		Help: "The upper limit of the supply temperature",
	})
	// 53=Heizgrenze Heizbetrieb
	// 54=Heizgrenze Absenken
	// 55=Frostschutz Aussentemperatur
	// 56=Vorhaltezeit Aufheizen
	// 57=Raumeinfluss
	// 58=Ueberhoehung Kessel
	// 59=Spreizung Heizkreis
	// 60=Minimale Drehzahl Pumpe PHK
	// 62=Mischer Laufzeit
	// 149=Brauchwassertemperatur normal
	// 150=Brauchwassertemperatur komfort
	// 155=Brauchwassertemperatur Schaltdifferenz
	// 158=Nachlauf Pumpe PK/LP
	// 162=Min. Laufzeit Kessel
	// 169=Nachlaufzeit Pumpe PZ
	// 171=Zirkulation Schaltdifferenz
	// Betriebszeit Kessel
	systacomfortBoilerRuntimeSeconds = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_runtime_seconds",
		Help: "The total runtime of the boiler",
	})
	// Anzahl Brennerstarts
	systacomfortBoilerStartsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_starts_total",
		Help: "The number of boiler starts",
	})
	// 220=Aktive Relais (Bitpattern)
	//   RelaisHeizkreispumpe    = 0x0001
	//   RelaisLadepumpe         = 0x0080
	//   RelaisZirkulationspumpe = 0x0100
	//   RelaisKessel            = 0x0200
	// Brenner aktiv wenn RelaisKessel && (KesselVorlauf - KesselRuecklauf > 2);

	systacomfortBoilerHeatercircuitPumpInfo = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_heatercircuitpump_info",
		Help: "Is the heatercircuit pump running",
	})
	systacomfortBoilerLoadpumpInfo = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_loadpump_info",
		Help: "Is the boiler loading running",
	})
	systacomfortBoilerCirculationPumpInfo = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_circulationpump_info",
		Help: "Is the boiler circulation pump running",
	})
	systacomfortBoilerActiveInfo = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_active_info",
		Help: "Is the boiler is running",
	})
	systacomfortBoilerTorchInfo = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_torch_info",
		Help: "Is the torch of the boiler active",
	})
	// Fehlerstatus (255 = OK)
	systacomfortBoilerErrorInfo = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_error_info",
		Help: "Errorinformation of the boiler",
	})
	// 232=Status
	systacomfortBoilerStatusInfo = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_status_info",
		Help: "State of the boiler",
	})

	// Kollektortempratur TSA
	systacomfortSolarpanelTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solarpanel_temperature_celsius",
		Help: "The temperature of the liquid in solar heating panel",
	})
	// Solar Rücklauf
	systacomfortSolarpanelReturnTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solarpanel_return_temperature_celsius",
		Help: "The temperature of the liquid streaming from the solar heating panel",
	})
	// Solar Vorlauf
	systacomfortSolarpanelSupplyTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solarpanel_supply_temperature_celsius",
		Help: "The temperature of the liquid streaming into the solar heating panel",
	})
	// Außentemperatur Kollektor TAM
	systacomfortSolarpanelOutsideTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solarpanel_outside_temperature_celsius",
		Help: "The outside temperature measured on the solar heating panel",
	})
	// Maximale Kollektortemperatur
	systacomfortSolarpanelMaximumTemperatureCelsius = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solarpanel_maximum_temperature_celsius",
		Help: "Maximum temperature of the liquid in the solar panal today",
	})
)
