//go:generate go run github.com/rakyll/statik -m -src=./ -Z -include=*.yaml -dest ../ -p stiebeleltron

package stiebeleltron

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/rakyll/statik/fs"
	log "github.com/sirupsen/logrus"
	"github.com/sp0x/surf"
	"github.com/sp0x/surf/browser"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type (
	SystemInfo struct {
		General struct {
			// OutsideTemperature in Celsius
			OutsideTemperature float64
			// CondenserTemperature in Celsius
			CondenserTemperature float64
			// Flow rate in liters/min
			Flow float64
			// HeatingCircuitPressure in bar
			HeatingCircuitPressure float64
			OutputActivity         struct {
				// HeatPump activity in %
				HeatPump float64
				// WaterPump activity in %
				WaterPump float64
			}
		}
		RoomTemperature struct {
			HeatingCircuit1 Temperature
			HeatingCircuit2 Temperature
		}
		DomesticHotWater Temperature
		Heating          struct {
			HeatingCircuit1 Temperature
			HeatingCircuit2 Temperature
			Flow            struct {
				// ActualTemperatureHeatPump in Celsius
				ActualTemperatureHeatPump float64
				// ActualTemperatureReheating in Celsius
				ActualTemperatureReheating float64
				// ActualPreFlowTemperature in Celsius
				ActualPreFlowTemperature float64
			}
			Buffer           Temperature
			FixedTemperature struct {
				// SetTemperature in Celsius
				SetTemperature float64
			}
		}
		ElectricReheating struct {
			DualModeReheating struct {
				// DomesticHotWaterTemperature in Celsius
				DomesticHotWaterTemperature float64
				// HeatingTemperature in Celsius
				HeatingTemperature float64
			}
		}
	}
	HeatPumpInfo struct {
		Runtime struct {
			// Heating runtime in hours
			Heating float64
			// DomesticHotWater runtime in hours
			DomesticHotWater float64
			// Reheating1 runtime in hours
			Reheating1 float64
			// Reheating2 runtime in hours
			Reheating2 float64
		}
		// RemainingCompressorRestingTime in seconds
		RemainingCompressorRestingTime float64
		// Energy accumulated
		Energy struct {
			Heating struct {
				// Day compressor energy so far
				Day float64
				// Total compressor energy so far
				Total float64
			}
			DomesticHotWater struct {
				// Day compressor energy so far
				Day float64
				// Total compressor energy so far
				Total float64
			}
			Reheating struct {
				// Total compressor energy so far
				Total float64
			}
		}
	}
	Temperature struct {
		// ActualTemperature in Celsius
		ActualTemperature float64
		// SetTemperature in Celsius
		SetTemperature float64
	}
	ISGClient struct {
		Options ClientOptions
		browser *browser.Browser
	}
	ClientOptions struct {
		URL     string
		Headers http.Header
	}
	Assignments map[string]Property
	Property    interface {
		GetGroup() string
		GetSearchString() string
		GetValue() float64
		SetValue(v float64)
	}
	PropertyImpl struct {
		Group        string
		SearchString string
		value        float64
	}
	ParseError struct {
		Group    string
		Property string
		Value    string
		Error    error
	}
)

var (
	SystemInfoPageQuery          = "?s=1,0"
	HeatPumpInfoPageQuery        = "?s=1,1"
	PropertyTableQueryExpression = "form#werte table.info tbody"
)

const (
	GeneralOutsideTemperature      = "General.OutsideTemperature"
	GeneralCondenserTemperature    = "General.CondenserTemperature"
	GeneralFlow                    = "General.Flow"
	GeneralHeatingCircuitPressure  = "General.HeatingCircuitPressure"
	GeneralOutputActivityHeatPump  = "General.OutputActivity.HeatPump"
	GeneralOutputActivityWaterPump = "General.OutputActivity.WaterPump"

	RoomTemperatureHeatingCircuit1ActualTemperature = "RoomTemperature.HeatingCircuit1.ActualTemperature"
	RoomTemperatureHeatingCircuit1SetTemperature    = "RoomTemperature.HeatingCircuit1.SetTemperature"
	RoomTemperatureHeatingCircuit2ActualTemperature = "RoomTemperature.HeatingCircuit2.ActualTemperature"
	RoomTemperatureHeatingCircuit2SetTemperature    = "RoomTemperature.HeatingCircuit2.SetTemperature"

	DomesticHotWaterActualTemperature                             = "DomesticHotWater.ActualTemperature"
	DomesticHotWaterSetTemperature                                = "DomesticHotWater.SetTemperature"
	ElectricReheatingDualModeReheatingHeatingTemperature          = "ElectricReheating.DualModeReheating.HeatingTemperature"
	ElectricReheatingDualModeReheatingDomesticHotWaterTemperature = "ElectricReheating.DualModeReheating.DomesticHotWaterTemperature"

	HeatingHeatingCircuit1ActualTemperature = "Heating.HeatingCircuit1.ActualTemperature"
	HeatingHeatingCircuit1SetTemperature    = "Heating.HeatingCircuit1.SetTemperature"
	HeatingHeatingCircuit2ActualTemperature = "Heating.HeatingCircuit2.ActualTemperature"
	HeatingHeatingCircuit2SetTemperature    = "Heating.HeatingCircuit2.SetTemperature"
	HeatingFlowActualTemperatureHeatPump    = "Heating.Flow.ActualTemperatureHeatPump"
	HeatingFlowActualTemperatureReheating   = "Heating.Flow.ActualTemperatureReheating"
	HeatingFlowActualPreFlowTemperature     = "Heating.Flow.ActualPreFlowTemperature"
	HeatingBufferActualTemperature          = "Heating.Buffer.ActualTemperature"
	HeatingBufferSetTemperature             = "Heating.Buffer.SetTemperature"
	HeatingFixedTemperatureSetTemperature   = "Heating.FixedTemperature.SetTemperature"

	RemainingCompressorRestingTime = "RemainingCompressorRestingTime"
	RuntimeHeating                 = "Runtime.Heating"
	RuntimeDomesticHotWater        = "Runtime.DomesticHotWater"
	RuntimeReheating1              = "Runtime.Reheating1"
	RuntimeReheating2              = "Runtime.Reheating2"

	EnergyHeatingDay            = "Energy.Heating.Day"
	EnergyHeatingTotal          = "Energy.Heating.Total"
	EnergyDomesticHotWaterDay   = "Energy.DomesticHotWater.Day"
	EnergyDomesticHotWaterTotal = "Energy.DomesticHotWater.Total"
	EnergyReheatingTotal        = "Energy.Reheating.Total"
)

func (a *PropertyImpl) GetGroup() string {
	return a.Group
}

func (a *PropertyImpl) GetSearchString() string {
	return a.SearchString
}

func (a *PropertyImpl) GetValue() float64 {
	return a.value
}

func (a *PropertyImpl) SetValue(v float64) {
	a.value = v
}

// NewISGClient constructs a client for interacting with Stiebel Eltron ISG.
func NewISGClient(options ClientOptions) (*ISGClient, error) {
	var br *browser.Browser
	br = surf.NewBrowser()
	for key, header := range options.Headers {
		for _, value := range header {
			br.AddRequestHeader(key, value)
		}
	}

	return &ISGClient{
		Options: options,
		browser: br,
	}, nil
}

func extractValue(v string) (float64, error) {
	arr := strings.Split(v, " ")
	rawValue := strings.ReplaceAll(arr[0], ",", ".")
	return strconv.ParseFloat(rawValue, 64)
}

func getProperties() []byte {
	statikFs, err := fs.New()
	if err != nil {
		log.WithError(err).Fatal("Cannot create internal filesystem")
	}
	r, err := statikFs.Open("/isg_english.yaml")
	if err != nil {
		log.WithError(err).Fatal("Cannot open embedded file")
	}
	defer r.Close()
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		log.WithError(err).Fatal(err)
	}
	return contents
}
func assignValue(a Assignments, group, property string, v float64) {
	for key, as := range a {
		if as.GetGroup() == group && as.GetSearchString() == property {
			as.SetValue(v)
			log.WithFields(log.Fields{
				"key":   key,
				"value": v,
			}).Debug("Assigned value")
			return
		}
	}
	log.WithFields(log.Fields{
		"group":    group,
		"property": property,
		"value":    v,
	}).Warn("Could not find a matching API property")
}

func collectSystemInfoValues(a Assignments) *SystemInfo {
	s := &SystemInfo{}
	s.General.OutsideTemperature = getValueFromAssignment(a, GeneralOutsideTemperature)
	s.General.CondenserTemperature = getValueFromAssignment(a, GeneralCondenserTemperature)
	s.General.Flow = getValueFromAssignment(a, GeneralFlow)
	s.General.HeatingCircuitPressure = getValueFromAssignment(a, GeneralHeatingCircuitPressure)
	s.General.OutputActivity.HeatPump = getValueFromAssignment(a, GeneralOutputActivityHeatPump)
	s.General.OutputActivity.WaterPump = getValueFromAssignment(a, GeneralOutputActivityWaterPump)

	s.RoomTemperature.HeatingCircuit1.ActualTemperature = getValueFromAssignment(a, RoomTemperatureHeatingCircuit1ActualTemperature)
	s.RoomTemperature.HeatingCircuit1.SetTemperature = getValueFromAssignment(a, RoomTemperatureHeatingCircuit1SetTemperature)
	s.RoomTemperature.HeatingCircuit2.ActualTemperature = getValueFromAssignment(a, RoomTemperatureHeatingCircuit2ActualTemperature)
	s.RoomTemperature.HeatingCircuit2.SetTemperature = getValueFromAssignment(a, RoomTemperatureHeatingCircuit2SetTemperature)

	s.DomesticHotWater.ActualTemperature = getValueFromAssignment(a, DomesticHotWaterActualTemperature)
	s.DomesticHotWater.SetTemperature = getValueFromAssignment(a, DomesticHotWaterSetTemperature)

	s.ElectricReheating.DualModeReheating.DomesticHotWaterTemperature = getValueFromAssignment(a, ElectricReheatingDualModeReheatingDomesticHotWaterTemperature)
	s.ElectricReheating.DualModeReheating.HeatingTemperature = getValueFromAssignment(a, ElectricReheatingDualModeReheatingHeatingTemperature)

	s.Heating.HeatingCircuit1.ActualTemperature = getValueFromAssignment(a, HeatingHeatingCircuit1ActualTemperature)
	s.Heating.HeatingCircuit1.SetTemperature = getValueFromAssignment(a, HeatingHeatingCircuit1SetTemperature)
	s.Heating.HeatingCircuit2.ActualTemperature = getValueFromAssignment(a, HeatingHeatingCircuit2ActualTemperature)
	s.Heating.HeatingCircuit2.SetTemperature = getValueFromAssignment(a, HeatingHeatingCircuit2SetTemperature)
	s.Heating.Flow.ActualPreFlowTemperature = getValueFromAssignment(a, HeatingFlowActualPreFlowTemperature)
	s.Heating.Flow.ActualTemperatureHeatPump = getValueFromAssignment(a, HeatingFlowActualTemperatureHeatPump)
	s.Heating.Flow.ActualTemperatureReheating = getValueFromAssignment(a, HeatingFlowActualTemperatureReheating)
	s.Heating.Buffer.ActualTemperature = getValueFromAssignment(a, HeatingBufferActualTemperature)
	s.Heating.Buffer.SetTemperature = getValueFromAssignment(a, HeatingBufferSetTemperature)
	s.Heating.FixedTemperature.SetTemperature = getValueFromAssignment(a, HeatingFixedTemperatureSetTemperature)
	return s
}

func collectHeatPumpInfoValues(a Assignments) *HeatPumpInfo {
	s := &HeatPumpInfo{}
	s.Energy.DomesticHotWater.Day = getValueFromAssignment(a, EnergyDomesticHotWaterDay)
	s.Energy.DomesticHotWater.Total = getValueFromAssignment(a, EnergyDomesticHotWaterTotal)
	s.Energy.Heating.Day = getValueFromAssignment(a, EnergyHeatingDay)
	s.Energy.Heating.Total = getValueFromAssignment(a, EnergyHeatingTotal)
	s.Energy.Reheating.Total = getValueFromAssignment(a, EnergyReheatingTotal)
	s.RemainingCompressorRestingTime = getValueFromAssignment(a, RemainingCompressorRestingTime)
	s.Runtime.Heating = getValueFromAssignment(a, RuntimeHeating)
	s.Runtime.DomesticHotWater = getValueFromAssignment(a, RuntimeDomesticHotWater)
	s.Runtime.Reheating1 = getValueFromAssignment(a, RuntimeReheating1)
	s.Runtime.Reheating2 = getValueFromAssignment(a, RuntimeReheating2)
	return s
}

func getValueFromAssignment(a Assignments, path string) float64 {
	return a[path].GetValue()
}

func (c *ISGClient) GetSystemInfo(a map[string]Property) (*SystemInfo, []ParseError, error) {
	err := c.LoadSystemInfoPage()
	if err != nil {
		return nil, nil, err
	}
	p := c.ParsePage(func(group, key string, value float64) {
		assignValue(a, group, key, value)
	})
	return collectSystemInfoValues(a), p, nil
}

func (c *ISGClient) GetHeatPumpInfo(a map[string]Property) (*HeatPumpInfo, []ParseError, error) {
	err := c.browser.Open(c.Options.URL + HeatPumpInfoPageQuery)
	if err != nil {
		return nil, nil, err
	}
	p := c.ParsePage(func(group, key string, value float64) {
		assignValue(a, group, key, value)
	})
	return collectHeatPumpInfoValues(a), p, nil
}

func (c *ISGClient) LoadSystemInfoPage() error {
	return c.browser.Open(c.Options.URL + SystemInfoPageQuery)
}

func (c *ISGClient) LoadHeatPumpInfoPage() error {
	return c.browser.Open(c.Options.URL + HeatPumpInfoPageQuery)
}

func (c *ISGClient) ParsePage(callback func(category, key string, value float64)) []ParseError {
	var p []ParseError
	c.browser.Find(PropertyTableQueryExpression).Each(func(i int, selection *goquery.Selection) {
		group := selection.Find("th").Text()
		selection.Find("tr.even,tr.odd").Each(func(i int, selection *goquery.Selection) {
			key := selection.Find("td.key").Text()
			value := strings.TrimSpace(selection.Find("td.value").Text())
			parsed, err := extractValue(value)
			if err != nil {
				p = append(p, ParseError{
					Group:    group,
					Property: key,
					Value:    value,
					Error:    err,
				})
				return
			}
			log.WithFields(log.Fields{
				"group": group,
				"key":   key,
				"value": parsed,
			}).Debug("Found property")
			callback(group, key, parsed)
		})
	})
	return p
}
