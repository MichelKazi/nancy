// Copyright 2020 Sonatype Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package logger has functions to obtain a logger, and helpers for setting up where the logger writes
package logger

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogger(t *testing.T) {
	LogLady.Level = logrus.DebugLevel
	// Do initial write to have a file
	LogLady.Debug("Test")

	if !strings.Contains(GetLogFileLocation(), "nancy.test.log") {
		t.Error("Nancy test file not in log file location")
	}

	err := os.Truncate(GetLogFileLocation(), 0)

	LogLady.Debug("Test")

	dat, err := ioutil.ReadFile(GetLogFileLocation())
	if err != nil {
		t.Error("Unable to open log file")
	}

	var logTest LogTest

	err = json.Unmarshal(dat, &logTest)
	if err != nil {
		t.Error("Improperly written log, should be valid json")
	}

	if logTest.Level != "debug" {
		t.Error("Log level not set properly")
	}

	if logTest.Msg != "Test" {
		t.Error("Message not written to log correctly")
	}
}

type LogTest struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
}