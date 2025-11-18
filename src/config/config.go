package config

import (
	"fmt"
	"oci-exporter/src/utils"
	"os"
	"strconv"
)

var (
	Port                   = os.Getenv("PORT")
	CompartmentId          = os.Getenv("COMPARTMENT_ID")
	CompartmentIdInSubtree = os.Getenv("COMPARTMENT_ID_IN_SUBTREE") == "true"
	OciConfigPath          = os.Getenv("OCI_CONFIG_PATH")
	OciConfigProfile       = os.Getenv("OCI_CONFIG_PROFILE")
	TimeRangeMinute        = 5
)

func InitConfig() {
	if Port == "" {
		Port = "8000"
		utils.Logger.Info(fmt.Sprintf("Using Default PORT %s", Port))
	}

	if CompartmentId == "" {
		utils.Logger.Error("COMPARTMENT_ID environment variable is required")
		os.Exit(1)
	} else {
		utils.Logger.Info(fmt.Sprintf("Using Compartment ID: %s", CompartmentId))
	}

	if CompartmentIdInSubtree {
		utils.Logger.Info("Using COMPARTMENT_ID_IN_SUBTREE: true, searching in all sub-compartments")
	} else {
		utils.Logger.Info("Using COMPARTMENT_ID_IN_SUBTREE: false, searching only in specified compartment")
	}

	if OciConfigPath == "" {
		OciConfigPath = os.Getenv("HOME") + "/.oci/config"
		utils.Logger.Info(fmt.Sprintf("Using Default OCI_CONFIG_PATH Path %s", OciConfigPath))
	} else {
		utils.Logger.Info(fmt.Sprintf("Using OCI_CONFIG_PATH Path: %s", OciConfigPath))
	}

	if OciConfigProfile == "" {
		OciConfigProfile = "DEFAULT"
		utils.Logger.Info(fmt.Sprintf("Using Default OCI_CONFIG_PROFILE %s", OciConfigProfile))
	} else {
		utils.Logger.Info(fmt.Sprintf("Using OCI_CONFIG_PROFILE: %s", OciConfigProfile))
	}

	time, err := strconv.Atoi(os.Getenv("TIME_RANGE_MINUTE"))
	if err != nil {
		utils.Logger.Info(fmt.Sprintf("Using Default TIME_RANGE_MINUTE %d", TimeRangeMinute))
	} else {
		TimeRangeMinute = time
		utils.Logger.Info(fmt.Sprintf("Using TIME_RANGE_MINUTE %d", TimeRangeMinute))
	}

}
