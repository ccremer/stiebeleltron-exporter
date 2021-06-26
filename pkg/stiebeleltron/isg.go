package stiebeleltron

import (
	_ "embed"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/sp0x/surf"
	"github.com/sp0x/surf/browser"
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
	PropertyTableQueryExpression = "form#werte table.info tbody"

	//go:embed isg_english.yaml
	IsgDefinition []byte
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
