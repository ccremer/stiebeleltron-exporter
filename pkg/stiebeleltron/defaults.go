package stiebeleltron

// NewSystemInfoDefaultAssignmentImpls creates a hardcoded internal PropertyImpl map in English for system info page.
func NewSystemInfoDefaultAssignments() Assignments {
	generalGroup := "GENERAL"
	heatingGroup := "HEATING"
	roomTemperatureGroup := "ROOM TEMPERATURE"
	return Assignments{
		GeneralOutsideTemperature:      &PropertyImpl{heatingGroup, "OUTSIDE TEMPERATURE", 0},
		GeneralCondenserTemperature:    &PropertyImpl{generalGroup, "CONDENSER TEMP.", 0},
		GeneralFlow:                    &PropertyImpl{generalGroup, "FLOW RATE", 0},
		GeneralHeatingCircuitPressure:  &PropertyImpl{generalGroup, "PRESSURE HTG CIRC", 0},
		GeneralOutputActivityHeatPump:  &PropertyImpl{generalGroup, "OUTPUT HP", 0},
		GeneralOutputActivityWaterPump: &PropertyImpl{generalGroup, "INT PUMP RATE", 0},

		RoomTemperatureHeatingCircuit1ActualTemperature: &PropertyImpl{roomTemperatureGroup, "ACTUAL TEMPERATURE HC 1", 0},
		RoomTemperatureHeatingCircuit1SetTemperature:    &PropertyImpl{roomTemperatureGroup, "SET TEMPERATURE HC 1", 0},
		RoomTemperatureHeatingCircuit2ActualTemperature: &PropertyImpl{roomTemperatureGroup, "ACTUAL TEMPERATURE HC 2", 0},
		RoomTemperatureHeatingCircuit2SetTemperature:    &PropertyImpl{roomTemperatureGroup, "SET TEMPERATURE HC 2", 0},

		DomesticHotWaterActualTemperature: &PropertyImpl{"DHW", "ACTUAL TEMPERATURE", 0},
		DomesticHotWaterSetTemperature:    &PropertyImpl{"DHW", "SET TEMPERATURE", 0},

		ElectricReheatingDualModeReheatingHeatingTemperature:          &PropertyImpl{"ELECTRIC REHEATING", "DUAL MODE TEMP HEATING", 0},
		ElectricReheatingDualModeReheatingDomesticHotWaterTemperature: &PropertyImpl{"ELECTRIC REHEATING", "DUAL MODE TEMP DHW", 0},

		HeatingHeatingCircuit1ActualTemperature: &PropertyImpl{heatingGroup, "ACTUAL TEMPERATURE HC 1", 0},
		HeatingHeatingCircuit1SetTemperature:    &PropertyImpl{heatingGroup, "SET TEMPERATURE HC 1", 0},
		HeatingHeatingCircuit2ActualTemperature: &PropertyImpl{heatingGroup, "ACTUAL TEMPERATURE HC 2", 0},
		HeatingHeatingCircuit2SetTemperature:    &PropertyImpl{heatingGroup, "SET TEMPERATURE HC 2", 0},
		HeatingFlowActualTemperatureHeatPump:    &PropertyImpl{heatingGroup, "ACTUAL FLOW TEMPERATURE WP", 0},
		HeatingFlowActualTemperatureReheating:   &PropertyImpl{heatingGroup, "ACTUAL FLOW TEMPERATURE NHZ", 0},
		HeatingFlowActualPreFlowTemperature:     &PropertyImpl{heatingGroup, "ACTUAL RETURN TEMPERATURE", 0},
		HeatingBufferActualTemperature:          &PropertyImpl{heatingGroup, "ACTUAL BUFFER TEMPERATURE", 0},
		HeatingBufferSetTemperature:             &PropertyImpl{heatingGroup, "SET BUFFER TEMPERATURE", 0},
		HeatingFixedTemperatureSetTemperature:   &PropertyImpl{heatingGroup, "SET FIXED TEMPERATURE", 0},
	}
}

// NewHeatPumpInfoDefaultAssignmentImpls creates a hardcoded internal PropertyImpl map in English for heat pump info page.
func NewHeatPumpInfoDefaultAssignments() Assignments {
	runtimeGroup := "RUNTIME"
	amountOfHeatGroup := "AMOUNT OF HEAT"
	return Assignments{
		RemainingCompressorRestingTime: &PropertyImpl{"PROCESS DATA", "COMP DLAY CNTR", 0},

		RuntimeHeating:          &PropertyImpl{runtimeGroup, "RNT COMP 1 HEA", 0},
		RuntimeDomesticHotWater: &PropertyImpl{runtimeGroup, "RNT COMP 1 DHW", 0},
		RuntimeReheating1:       &PropertyImpl{runtimeGroup, "BH 1", 0},
		RuntimeReheating2:       &PropertyImpl{runtimeGroup, "BH 2", 0},

		EnergyHeatingDay:            &PropertyImpl{amountOfHeatGroup, "COMPRESSOR HEATING DAY", 0},
		EnergyHeatingTotal:          &PropertyImpl{amountOfHeatGroup, "COMPRESSOR HEATING TOTAL", 0},
		EnergyDomesticHotWaterDay:   &PropertyImpl{amountOfHeatGroup, "COMPRESSOR DHW DAY", 0},
		EnergyDomesticHotWaterTotal: &PropertyImpl{amountOfHeatGroup, "COMPRESSOR DHW TOTAL", 0},
		EnergyReheatingTotal:        &PropertyImpl{amountOfHeatGroup, "BH HEATING TOTAL", 0},
	}
}
