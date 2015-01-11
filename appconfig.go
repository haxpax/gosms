package main

import (
	"errors"
	"strings"
	ini "github.com/vaughan0/go-ini"
)

/* ===== Application Configuration ===== */
type setting []string

func getConfig(configFilePath string) (ini.File, error) {
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

	required_fields := []setting{{"SETTINGS", "COMPORT"},
		{"SETTINGS", "BAUDRATE"},
		{"SETTINGS", "RETRIES"},
	}

	for _, c := range required_fields {
		v, ok := appConfig.Get(c[0], c[1])
		if !ok || strings.TrimSpace(v) == "" {
			return false, errors.New("Fatal: " + c[1] + " is not set")
		}

	}
	return true, nil
}

/* ===== Application Configuration ===== */
