package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"stiebeleltron-exporter/pkg/stiebeleltron"
)

var (
	Namespace = "stiebeleltron"
)

type MetricProperty struct {
	GaugeName        string
	Labels           map[string]string
	HelpText         string
	Gauge            prometheus.Gauge
	PropertyGroup    string
	SearchString     string
	ValueTransformer func(v float64) float64
}

func NewDefaultMetricProperties() map[string]*MetricProperty {
	heatingCircuitRoomTemperatureGaugeName := "heatingcircuit_room_temperature"
	heatingCircuitLabel := "heating_circuit"
	heatingCircuitRoomTemperatureHelpText := "Room temperature in heating circuit in Celsius"
	heatingBufferTemperatureHelpText := "Buffer temperature in Celsius"
	temperatureTypeLabel := "type"
	heatingCircuitTemperatureHelpText := "Heating temperatur in Celsius"
	heatingCircuitTemperatureGaugeName := "heating_temperature"

	domesticHotWaterHelpText := "Domestic hot water temperature in Celcius"
	domesticHotWaterGaugeName := "domestic_hot_water_temperature"

	energyHelpText := "Compressor energy in Ws"
	energyGaugeName := "compressor_energy"

	timeframeLabel := "timeframe"
	runtimeHelpText := "Total compressor runtime in seconds"
	runtimeGaugeName := "runtime_seconds_total"
	componentLabel := "component"

	m := map[string]*MetricProperty{
		stiebeleltron.GeneralOutsideTemperature:      {HelpText: "Outside temperature in Celsius", GaugeName: "general_outside_temperature"},
		stiebeleltron.GeneralCondenserTemperature:    {HelpText: "Condenser temperature in Celsius", GaugeName: "general_condenser_temperature"},
		stiebeleltron.GeneralFlow:                    {HelpText: "Flow rate in liters per second", GaugeName: "general_flow_rate", ValueTransformer: divideBy60},
		stiebeleltron.GeneralOutputActivityWaterPump: {HelpText: "Water pump activity in percent", GaugeName: "general_outputactivity_waterpump", ValueTransformer: divideBy100},
		stiebeleltron.GeneralOutputActivityHeatPump:  {HelpText: "Heat pump activity in percent", GaugeName: "general_outputactivity_heatpump", ValueTransformer: divideBy100},
		stiebeleltron.GeneralHeatingCircuitPressure:  {HelpText: "Heating circuit pressure in bar", GaugeName: "general_heatingcircuit_pressure"},

		stiebeleltron.RoomTemperatureHeatingCircuit1SetTemperature: {HelpText: heatingCircuitRoomTemperatureHelpText, GaugeName: heatingCircuitRoomTemperatureGaugeName, Labels: prometheus.Labels{
			heatingCircuitLabel:  "1",
			temperatureTypeLabel: "set",
		}},
		stiebeleltron.RoomTemperatureHeatingCircuit2SetTemperature: {HelpText: heatingCircuitRoomTemperatureHelpText, GaugeName: heatingCircuitRoomTemperatureGaugeName, Labels: prometheus.Labels{
			heatingCircuitLabel:  "2",
			temperatureTypeLabel: "set",
		}},
		stiebeleltron.RoomTemperatureHeatingCircuit1ActualTemperature: {HelpText: heatingCircuitRoomTemperatureHelpText, GaugeName: heatingCircuitRoomTemperatureGaugeName, Labels: prometheus.Labels{
			heatingCircuitLabel:  "1",
			temperatureTypeLabel: "actual",
		}},
		stiebeleltron.RoomTemperatureHeatingCircuit2ActualTemperature: {HelpText: heatingCircuitRoomTemperatureHelpText, GaugeName: heatingCircuitRoomTemperatureGaugeName, Labels: prometheus.Labels{
			heatingCircuitLabel:  "2",
			temperatureTypeLabel: "actual",
		}},

		stiebeleltron.HeatingBufferSetTemperature: {HelpText: heatingBufferTemperatureHelpText, GaugeName: "heating_buffer_temperature", Labels: prometheus.Labels{
			temperatureTypeLabel: "set",
		}},
		stiebeleltron.HeatingBufferActualTemperature: {HelpText: heatingBufferTemperatureHelpText, GaugeName: "heating_buffer_temperature", Labels: prometheus.Labels{
			temperatureTypeLabel: "actual",
		}},
		stiebeleltron.HeatingFixedTemperatureSetTemperature: {HelpText: "Fixed temperature in heating in Celcius", GaugeName: "heating_fixed_temperature", Labels: prometheus.Labels{
			temperatureTypeLabel: "set",
		}},
		stiebeleltron.HeatingHeatingCircuit1SetTemperature: {HelpText: heatingCircuitTemperatureHelpText, GaugeName: heatingCircuitTemperatureGaugeName, Labels: prometheus.Labels{
			heatingCircuitLabel:  "1",
			temperatureTypeLabel: "set",
		}},
		stiebeleltron.HeatingHeatingCircuit2SetTemperature: {HelpText: heatingCircuitTemperatureHelpText, GaugeName: heatingCircuitTemperatureGaugeName, Labels: prometheus.Labels{
			heatingCircuitLabel:  "2",
			temperatureTypeLabel: "set",
		}},
		stiebeleltron.HeatingHeatingCircuit1ActualTemperature: {HelpText: heatingCircuitTemperatureHelpText, GaugeName: heatingCircuitTemperatureGaugeName, Labels: prometheus.Labels{
			heatingCircuitLabel:  "1",
			temperatureTypeLabel: "actual",
		}},
		stiebeleltron.HeatingHeatingCircuit2ActualTemperature: {HelpText: heatingCircuitTemperatureHelpText, GaugeName: heatingCircuitTemperatureGaugeName, Labels: prometheus.Labels{
			heatingCircuitLabel:  "2",
			temperatureTypeLabel: "actual",
		}},

		stiebeleltron.HeatingFlowActualTemperatureReheating: {HelpText: "Reheating flow temperature in Celcius", GaugeName: "reheating_flow_temperature"},
		stiebeleltron.HeatingFlowActualTemperatureHeatPump:  {HelpText: "Heating flow heat pump temperature in Celcius", GaugeName: "heating_flow_heatpump_temperature"},
		stiebeleltron.HeatingFlowActualPreFlowTemperature:   {HelpText: "Heating pre-flow temperature in Celcius", GaugeName: "heating_preflow_temperature"},

		stiebeleltron.ElectricReheatingDualModeReheatingHeatingTemperature:          {HelpText: "Heating temperature in Celcius with eletric reheating", GaugeName: "electric_reheating_temperature"},
		stiebeleltron.ElectricReheatingDualModeReheatingDomesticHotWaterTemperature: {HelpText: "Domestic hot water temperature in Celcius with electric reheating", GaugeName: "electric_reheating_hotwater_temperature"},

		stiebeleltron.DomesticHotWaterSetTemperature: {HelpText: domesticHotWaterHelpText, GaugeName: domesticHotWaterGaugeName, Labels: prometheus.Labels{
			temperatureTypeLabel: "set",
		}},
		stiebeleltron.DomesticHotWaterActualTemperature: {HelpText: domesticHotWaterHelpText, GaugeName: domesticHotWaterGaugeName, Labels: prometheus.Labels{
			temperatureTypeLabel: "actual",
		}},

		stiebeleltron.RemainingCompressorRestingTime: {HelpText: "Remaining compressor delay time in seconds", GaugeName: "remaining_compressor_rest_seconds"},
		stiebeleltron.EnergyHeatingDay: {HelpText: energyHelpText, GaugeName: energyGaugeName, ValueTransformer: kWhToWs, Labels: prometheus.Labels{
			timeframeLabel: "day",
			componentLabel: "heating",
		}},
		stiebeleltron.EnergyHeatingTotal: {HelpText: energyHelpText, GaugeName: energyGaugeName, ValueTransformer: mWhToWs, Labels: prometheus.Labels{
			timeframeLabel: "total",
			componentLabel: "heating",
		}},
		stiebeleltron.EnergyDomesticHotWaterDay: {HelpText: energyHelpText, GaugeName: energyGaugeName, ValueTransformer: kWhToWs, Labels: prometheus.Labels{
			timeframeLabel: "day",
			componentLabel: "domestic_hot_water",
		}},
		stiebeleltron.EnergyDomesticHotWaterTotal: {HelpText: energyHelpText, GaugeName: energyGaugeName, ValueTransformer: mWhToWs, Labels: prometheus.Labels{
			timeframeLabel: "total",
			componentLabel: "domestic_hot_water",
		}},
		stiebeleltron.EnergyReheatingTotal: {HelpText: energyHelpText, GaugeName: energyGaugeName, ValueTransformer: mWhToWs, Labels: prometheus.Labels{
			timeframeLabel: "total",
			componentLabel: "reheating",
		}},
		stiebeleltron.RuntimeHeating: {HelpText: runtimeHelpText, GaugeName: runtimeGaugeName, ValueTransformer: multiplyBy3600, Labels: prometheus.Labels{
			componentLabel: "heating",
		}},
		stiebeleltron.RuntimeDomesticHotWater: {HelpText: runtimeHelpText, GaugeName: runtimeGaugeName, ValueTransformer: multiplyBy3600, Labels: prometheus.Labels{
			componentLabel: "domestic_hot_water",
		}},
		stiebeleltron.RuntimeReheating1: {HelpText: runtimeHelpText, GaugeName: runtimeGaugeName, ValueTransformer: multiplyBy3600, Labels: prometheus.Labels{
			componentLabel: "reheating1",
		}},
		stiebeleltron.RuntimeReheating2: {HelpText: runtimeHelpText, GaugeName: runtimeGaugeName, ValueTransformer: multiplyBy3600, Labels: prometheus.Labels{
			componentLabel: "reheating2",
		}},
	}
	return m
}

func divideBy100(value float64) float64 {
	return value / 100
}

func divideBy60(value float64) float64 {
	return value / 60
}

func multiplyBy3600(value float64) float64 {
	return value * 3600
}

func kWhToWs(value float64) float64 {
	return value * 3.6e+6
}

func mWhToWs(value float64) float64 {
	return value * 3.6e+9
}
