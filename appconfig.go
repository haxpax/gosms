package gosms

import (
	"errors"
	"fmt"
	ini "github.com/vaughan0/go-ini"
	"strconv"
	"strings"
)

/* ===== Application Configuration ===== */
type setting []string

func GetConfig(configFilePath string) (ini.File, error) {
	appConfig, err := ini.LoadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	ok, err := testConfig(appConfig)
	if !ok {
		return nil, err
	}
	return appConfig, nil
}

func testConfig(appConfig ini.File) (bool, error) {
	//test if required parameters are present and are valid

	requiredFields := []setting{{"SETTINGS", "SERVERHOST"},
		{"SETTINGS", "SERVERPORT"},
		{"SETTINGS", "RETRIES"},
		{"SETTINGS", "DEVICES"},
	}

	for _, c := range requiredFields {
		v, ok := appConfig.Get(c[0], c[1])
		if !ok || strings.TrimSpace(v) == "" {
			return false, errors.New("Fatal: " + c[1] + " is not set")
		}
	}

	//now make sure all the devices are have required settings
	tno, _ := appConfig.Get("SETTINGS", "DEVICES")
	noOfDevices, _ := strconv.Atoi(tno)
	requiredFields = []setting{}

	for i := 0; i <= noOfDevices; i++ {
		d := fmt.Sprintf("DEVICE%v", i)
		sCom := setting{d, "COMPORT"}
		sBaud := setting{d, "BAUDRATE"}
		requiredFields = append(requiredFields, sCom)
		requiredFields = append(requiredFields, sBaud)
	}

	for _, c := range requiredFields {
		v, ok := appConfig.Get(c[0], c[1])
		if !ok || strings.TrimSpace(v) == "" {
			return false, errors.New(strings.Join([]string{"Fatal:", c[0], c[1], "is not set"}, " "))
		}
	}

	return true, nil
}

/* ===== Application Configuration ===== */
