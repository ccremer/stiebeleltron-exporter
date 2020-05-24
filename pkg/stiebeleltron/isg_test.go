package stiebeleltron

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestISGClient_GetSystemInfo_GivenHTML_WhenDefaultAssignment_ThenParseValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		payload, err := ioutil.ReadFile("testdata/systeminfo_1.html")
		require.NoError(t, err)
		rw.Write(payload)
	}))
	c, err := NewISGClient(ClientOptions{
		URL: server.URL,
	})
	log.SetLevel(log.DebugLevel)
	require.NoError(t, err)

	result, parseErrors, err := c.GetSystemInfo(NewSystemInfoDefaultAssignments())
	assert.NoError(t, err)
	assert.Empty(t, parseErrors)
	assert.Equal(t, createExampleSystemInfo(), result)
}

func TestISGClient_GetHeatPumpInfo_GivenHTML_WhenDefaultAssignment_ThenParseValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		payload, err := ioutil.ReadFile("testdata/heatpumpinfo_1.html")
		require.NoError(t, err)
		rw.Write(payload)
	}))
	c, err := NewISGClient(ClientOptions{
		URL: server.URL,
	})
	log.SetLevel(log.DebugLevel)
	require.NoError(t, err)

	result, parseErrors, err := c.GetHeatPumpInfo(NewHeatPumpInfoDefaultAssignments())
	assert.NoError(t, err)
	assert.Empty(t, parseErrors)
	assert.Equal(t, createExampleHeatPumpInfo(), result)
}

func createExampleSystemInfo() *SystemInfo {
	s := &SystemInfo{}
	s.General.OutsideTemperature = 17.9
	s.General.OutputActivity.WaterPump = 2.2
	s.General.OutputActivity.HeatPump = 13
	s.General.Flow = 0.4
	s.General.HeatingCircuitPressure = 1.23
	s.General.CondenserTemperature = 31.4

	s.RoomTemperature.HeatingCircuit1.SetTemperature = 21.6
	s.RoomTemperature.HeatingCircuit1.ActualTemperature = 23.5
	s.RoomTemperature.HeatingCircuit2.SetTemperature = 21.6
	s.RoomTemperature.HeatingCircuit2.ActualTemperature = 23.6

	s.Heating.Flow.ActualTemperatureReheating = 33.9
	s.Heating.Flow.ActualTemperatureHeatPump = 33.5
	s.Heating.Flow.ActualPreFlowTemperature = 30.9
	s.Heating.FixedTemperature.SetTemperature = 42
	s.Heating.Buffer.SetTemperature = 42
	s.Heating.Buffer.ActualTemperature = 46.1
	s.Heating.HeatingCircuit1.ActualTemperature = 46.1
	s.Heating.HeatingCircuit1.SetTemperature = 42
	s.Heating.HeatingCircuit2.ActualTemperature = 25.9
	s.Heating.HeatingCircuit2.SetTemperature = 23.2

	s.ElectricReheating.DualModeReheating.HeatingTemperature = -13
	s.ElectricReheating.DualModeReheating.DomesticHotWaterTemperature = -13

	s.DomesticHotWater.SetTemperature = 44.5
	s.DomesticHotWater.ActualTemperature = 47.4
	return s
}

func createExampleHeatPumpInfo() *HeatPumpInfo {
	h := &HeatPumpInfo{}
	h.RemainingCompressorRestingTime = 1
	h.Runtime.Reheating1 = 0
	h.Runtime.Reheating2 = 2
	h.Runtime.DomesticHotWater = 1771
	h.Runtime.Heating = 523
	h.Energy.Heating.Day = 21.145
	h.Energy.Heating.Total = 56.97
	h.Energy.DomesticHotWater.Day = 5.052
	h.Energy.DomesticHotWater.Total = 12.617
	h.Energy.Reheating.Total = 0.02
	return h
}
