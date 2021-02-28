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
	// Heizung Vorlauf HK1
	systacomfortHeatercircuitSupplyTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatercircuit_supply_temperature_celsius",
		Help: "The boiler supply temperature",
	})
	// Heizung Ruecklauf HK1
	systacomfortHeatercircuitReturnTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatercircuit_return_temperature_celsius",
		Help: "The boiler return temperature",
	})
	// Brauchwasser TWO
	systacomfortTapwaterTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_tapwater_temperature_celsius",
		Help: "The tapwater temperature",
	})
	// Puffertemperatur oben TPO
	systacomfortBufferToplayerTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_buffer_toplayer_temperature_celsius",
		Help: "The temperature in the top layer of the buffer",
	})
	// Puffertemperatur unten TPU
	systacomfortBufferBottomlayerTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_buffer_bottomlayer_temperature_celsius",
		Help: "The temperature in the bottom layer of the buffer",
	})
	// Zirkulation
	systacomfortCirculationTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_watercirculation_temperature_celsius",
		Help: "The temperature in the water circulation",
	})
	// Heizung Vorlauf HK2
	systacomfortHeatercircuit2SupplyTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatercircuit2_supply_temperature_celsius",
		Help: "The boiler supply temperature - Circuit 2",
	})
	// Heizung Ruecklauf HK2
	systacomfortHeatercircuitReturn2TemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatercircuit2_return_temperature_celsius",
		Help: "The boiler return temperature",
	})
	// Raumtemperatur
	systacomfortInsideTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_temperature_celsius",
		Help: "The inside temperature",
	})
	// Raumtemperatur HK2
	systacomfortInside2TemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside2_temperature_celsius",
		Help: "The inside temperature - Circuit 2",
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
	// Vorlauf Heizung Soll HK1
	systacomfortHeatercircuitTargetSupplyTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatercircuit_target_supply_temperature_celsius",
		Help: "The target supply temperature",
	})
	// Vorlauf Heizung Soll HK2
	systacomfortHeatercircuit2TargetSupplyTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatercircuit2_target_supply_temperature_celsius",
		Help: "The target supply temperature Circuit 2",
	})
	// 33 - Puffertemperatur Soll
	systacomfortBufferTargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_buffer_target_temperature_celsius",
		Help: "The target temperature for the buffer",
	})
	// 34 - KesselSoll
	systacomfortBoilerTargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_target_temperature_celsius",
		Help: "The boiler target temperature",
	})
	// 36 - Betriebsart 0=Heizprogramm 1, 1=Heizprogramm 2, 2=Heizprogramm 3, 3=Dauernd Normal, 4=Dauernd Komfort, 5=Dauernd Absenken
	systacomfortModeInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_mode_info",
		Help: "The current mode",
	})
	// 39=Raumtemperatur normal (soll)
	systacomfortInsideNormalTargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_normal_target_temperature_celsius",
		Help: "The normal inside target temperature",
	})
	// 40=Raumtemperatur komfort (soll)
	systacomfortInsideComfortTargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_comfort_target_temperature_celsius",
		Help: "The comfort inside target temperature",
	})
	// 41=Raumtemperatur abgesenkt (soll)
	systacomfortInsideReducedTargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_reduced_target_temperature_celsius",
		Help: "The reduced inside target temperature",
	})
	// 48 - Fusspunkt
	systacomfortHeatCurveRootPointCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatcurve_rootpoint_celsius",
		Help: "The rootpoint of the heating curve",
	})
	// 50 - Steilheit
	systacomfortHeatCurveGradientCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_heatcurve_gradient_celsius",
		Help: "The gradient of the heating curve",
	})
	// 52 - Max. Vorlauftemperatur
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
	// 70=Speicher Sollwert
	systacomfortReservoirTargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_reservoir_target_temperature_celsius",
		Help: "The reservoir target temperature",
	})
	// 78=Raumtemperatur normal - HK2 (soll)
	systacomfortInsideNormal2TargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_normal2_target_temperature_celsius",
		Help: "The normal inside target temperature - Circuit 2",
	})
	// 79=Raumtemperatur komfort - HK2 (soll)
	systacomfortInsideComfort2TargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_comfort2_target_temperature_celsius",
		Help: "The comfort inside target temperature - Circuit 2",
	})
	// 80=Raumtemperatur abgesenkt - HK2 (soll)
	systacomfortInsideReduced2TargetTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_inside_reduced2_target_temperature_celsius",
		Help: "The reduced inside target temperature - Circuit 2",
	})
	// 149=Brauchwassertemperatur normal
	// 150=Brauchwassertemperatur komfort
	// 155=Brauchwassertemperatur Schaltdifferenz
	// 158=Nachlauf Pumpe PK/LP
	// 162=Min. Laufzeit Kessel
	// 169=Nachlaufzeit Pumpe PZ
	// 171=Zirkulation Schaltdifferenz
	// 180 - Betriebszeit Kessel
	systacomfortBoilerRuntimeSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_runtime_seconds",
		Help: "The total runtime of the boiler",
	})
	// 181 - Anzahl Brennerstarts
	systacomfortBoilerStartsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_boiler_starts_total",
		Help: "The number of boiler starts",
	})
	// 182=Solare Leistung
	systacomfortSolarOutputKWh = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solar_output_kwh",
		Help: "The solar output in kWh",
	})
	// 183=Tagesgewinn Solare Leistung
	systacomfortSolarOutputDayKWh = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solar_output_day_kwh",
		Help: "The daily solar output in kWh",
	})
	// 184=Gesamtgewinn Solare Leistung
	systacomfortSolarOutputTotalKWh = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systacomfort_solar_output_total_kwh",
		Help: "The total solar output in kWh",
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
	// 156=Maximale Kollektortemperatur
	//systacomfortSolarpanelMaximumTemperatureCelsius = prometheus.NewGauge(prometheus.GaugeOpts{
	//	Name: "systacomfort_solarpanel_maximum_temperature_celsius",
	//	Help: "Daily Maxima of the solarpanel temperature",
	//})

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
