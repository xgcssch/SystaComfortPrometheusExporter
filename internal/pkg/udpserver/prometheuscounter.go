//
// Derived from a perl script developed by Klaus.Schmidinger@tvdr.de

package udpserve

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Außentemperatur Fühler
	systacomfortOutsideTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_outside_temperature_celsius",
		Help: "The outside temperature",
	})
	// Heizung Vorlauf
	systacomfortHeatercircuitSupplyTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatercircuit_supply_temperature_celsius",
		Help: "The boiler supply temperature",
	})
	// Heizung Ruecklauf
	systacomfortHeatercircuitReturnTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatercircuit_return_temperature_celsius",
		Help: "The boiler return temperature",
	})
	// Brauchwasser
	systacomfortTapwaterTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_tapwater_temperature_celsius",
		Help: "The tapwater temperature",
	})
	// Zirkulation
	systacomfortCirculationTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_watercirculation_temperature_celsius",
		Help: "The temperature in the water circulation",
	})
	// Raumtemperatur
	systacomfortInsideTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_temperature_celsius",
		Help: "The inside temperature",
	})
	// Kollektor
	// KesselVorlauf
	systacomfortBoilerSupplyTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_supply_temperature_celsius",
		Help: "The boiler supply temperature",
	})
	// KesselRuecklauf
	systacomfortBoilerReturnTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_return_temperature_celsius",
		Help: "The boiler return temperature",
	})
	// BrauchwasserSoll
	systacomfortTapwaterTargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_tapwater_target_temperature_celsius",
		Help: "The tapwater target temperature",
	})
	// InnenSoll
	systacomfortInsideTargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_target_temperature_celsius",
		Help: "The inside target temperature",
	})
	// KesselSoll
	systacomfortBoilerTargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_target_temperature_celsius",
		Help: "The boiler target temperature",
	})
	// Betriebsart 0=Heizprogramm 1, 1=Heizprogramm 2, 2=Heizprogramm 3, 3=Dauernd Normal, 4=Dauernd Komfort, 5=Dauernd Absenken
	systacomfortModeInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_mode_info",
		Help: "The current mode",
	})
	// 39=Raumtemperatur normal (soll)
	// 40=Raumtemperatur komfort (soll)
	// 41=Raumtemperatur abgesenkt (soll)
	// Fusspunkt
	systacomfortHeatCurveRootPointCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatcurve_rootpoint_celsius",
		Help: "The rootpoint of the heating curve",
	})
	// Steilheit
	systacomfortHeatCurveGradientCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatcurve_gradient_celsius",
		Help: "The gradient of the heating curve",
	})
	// Max. Vorlauftemperatur
	systacomfortBoilerSupplyUpperLimitCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
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
	systacomfortBoilerRuntimeSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_runtime_seconds",
		Help: "The total runtime of the boiler",
	})
	// Anzahl Brennerstarts
	systacomfortBoilerStartsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_starts_total",
		Help: "The number of boiler starts",
	})
	// 220=Aktive Relais (Bitpattern)
	//   RelaisHeizkreispumpe    = 0x0001
	//   RelaisLadepumpe         = 0x0080
	//   RelaisZirkulationspumpe = 0x0100
	//   RelaisKessel            = 0x0200
	// Brenner aktiv wenn RelaisKessel && (KesselVorlauf - KesselRuecklauf > 2);

	systacomfortBoilerHeatercircuitPumpInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_heatercircuitpump_info",
		Help: "Is the heatercircuit pump running",
	})
	systacomfortBoilerLoadpumpInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_loadpump_info",
		Help: "Is the boiler loading running",
	})
	systacomfortBoilerCirculationPumpInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_circulationpump_info",
		Help: "Is the boiler circulation pump running",
	})
	systacomfortBoilerActiveInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_active_info",
		Help: "Is the boiler is running",
	})
	systacomfortBoilerTorchInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_torch_info",
		Help: "Is the torch of the boiler active",
	})
	// Fehlerstatus (255 = OK)
	systacomfortBoilerErrorInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_error_info",
		Help: "Errorinformation of the boiler",
	})
	// 232=Status
	systacomfortBoilerStatusInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_status_info",
		Help: "State of the boiler",
	})

	// Kollektortempratur TSA
	systacomfortSolarpanelTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solarpanel_temperature_celsius",
		Help: "The temperature of the liquid in solar heating panel",
	})
	// Solar Rücklauf
	systacomfortSolarpanelReturnTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solarpanel_return_temperature_celsius",
		Help: "The temperature of the liquid streaming from the solar heating panel",
	})
	// Solar Vorlauf
	systacomfortSolarpanelSupplyTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solarpanel_supply_temperature_celsius",
		Help: "The temperature of the liquid streaming into the solar heating panel",
	})
	// Außentemperatur Kollektor TAM
	systacomfortSolarpanelOutsideTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solarpanel_outside_temperature_celsius",
		Help: "The outside temperature measured on the solar heating panel",
	})
	// Maximale Kollektortemperatur
	systacomfortSolarpanelMaximumTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solarpanel_maximum_temperature_celsius",
		Help: "Maximum temperature of the liquid in the solar panal today",
	})
)

func registerCounter(registry *prometheus.Registry) {
	// Metrics have to be registered to be exposed:
	registry.MustRegister(
        systacomfortOutsideTemperatureCelsius,
        systacomfortHeatercircuitSupplyTemperatureCelsius,
        systacomfortHeatercircuitReturnTemperatureCelsius,
        systacomfortTapwaterTemperatureCelsius,
        systacomfortCirculationTemperatureCelsius,
        systacomfortInsideTemperatureCelsius,
        systacomfortBoilerSupplyTemperatureCelsius,
        systacomfortBoilerReturnTemperatureCelsius,
        systacomfortTapwaterTargetTemperatureCelsius,
        systacomfortInsideTargetTemperatureCelsius,
        systacomfortBoilerTargetTemperatureCelsius,
        systacomfortModeInfo,
        systacomfortHeatCurveRootPointCelsius,
        systacomfortHeatCurveGradientCelsius,
        systacomfortBoilerSupplyUpperLimitCelsius,
        systacomfortBoilerRuntimeSeconds,
        systacomfortBoilerStartsTotal,
        systacomfortBoilerHeatercircuitPumpInfo,
        systacomfortBoilerLoadpumpInfo,
        systacomfortBoilerCirculationPumpInfo,
        systacomfortBoilerActiveInfo,
        systacomfortBoilerTorchInfo,
        systacomfortBoilerErrorInfo,
        systacomfortBoilerStatusInfo,
        systacomfortSolarpanelTemperatureCelsius,
        systacomfortSolarpanelReturnTemperatureCelsius,
        systacomfortSolarpanelSupplyTemperatureCelsius,
        systacomfortSolarpanelOutsideTemperatureCelsius,
        systacomfortSolarpanelMaximumTemperatureCelsius,
    )
}
