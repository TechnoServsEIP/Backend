package model

import (
	"testing"
	"os/exec"
)

func TestkebabToCamelCase(t *testing.T) {
	kebabString := []string{
		"generate-structures",
		"online-mode",
		"max-build-height",
		"level-seed",
		"use-native-transport",
		"prevent-proxy-connections"
		"enable-jmx-monitoring",
	}

	if (kebabToCamelCase(kebabString[0]) != "GenerateStructures") {
		t.Errorf("Error for string GenerateStructures")
	}

	if (kebabToCamelCase(kebabString[1]) != "OnlineMode") {
		t.Errorf("Error for string OnlineMode")
	}

	if (kebabToCamelCase(kebabString[2]) != "MaxBuildHeight") {
		t.Errorf("Error for string MaxBuildHeight")
	}

	if (kebabToCamelCase(kebabString[3]) != "LevelSeed") {
		t.Errorf("Error for string LevelSeed")
	}

	if (kebabToCamelCase(kebabString[4]) != "UseNativeTransport") {
		t.Errorf("Error for string UseNativeTransport")
	}

	if (kebabToCamelCase(kebabString[5]) != "PreventProxyConnections") {
		t.Errorf("Error for string PreventProxyConnections")
	}

	if (kebabToCamelCase(kebabString[6]) != "EnableJmxMonitoring") {
		t.Errorf("Error for string EnableJmxMonitoring")
	}
}