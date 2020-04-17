//
// Copyright 2018-present Sonatype Inc.
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
//

package customerrors

import (
	"fmt"
	"github.com/sonatype-nexus-community/nancy/buildversion"
	. "github.com/sonatype-nexus-community/nancy/logger"
	"os"
	"runtime"
)

type exiter interface {
	Exit(code int)
}

var activeExiter exiter

type defaultExit struct{}

func (e defaultExit) Exit(code int) {
	os.Exit(code)
}
func init() {
	ResetExiter()
}

type bypassExit struct{}

func (e bypassExit) Exit(code int) {
	fmt.Println(GetBypassMessage(code))
}
func GetBypassMessage(code int) string {
	return fmt.Sprintf("*** bypassing exit, code: %d ***", code)
}

func ResetExiter() {
	activeExiter = defaultExit{}
}
func BypassExiter() {
	activeExiter = &bypassExit{}
}

type swError struct {
	Message string
	Err     error
	Exiter  exiter
}

func newSwError(message string, err error) *swError {
	sw := swError{Message: message, Err: err}
	sw.Exiter = activeExiter
	return &sw
}

func (sw swError) Error() string {
	return fmt.Sprintf("%s - error: %s", sw.Message, sw.Err.Error())
}

func (sw swError) Exit() {
	sw.Exiter.Exit(3)
}

func Check(err error, message string) {
	if err != nil {
		myErr := newSwError(message, err)
		LogLady.WithField("error", err).Error(message)
		fmt.Println(myErr.Error())
		fmt.Printf("For more information, check the log file at %s\n", GetLogFileLocation())
		fmt.Println("nancy version:", buildversion.BuildVersion)
		myErr.Exit()
	}
}

func getCallerFunction(skip int) string {
	if skip > 10 {
		LogLady.Errorf("getCallerFunction called with invalid skip value: %d", skip)
	}
	programCounters := make([]uintptr, 10)
	runtime.Callers(0, programCounters)

	// for debugging
	/*	callerNames := [10]string{}
		for idx, pc := range programCounters {
			if pc != 0 {
				callerNames[idx] = runtime.FuncForPC(programCounters[idx]).Name()
			}
		}
	*/
	return runtime.FuncForPC(programCounters[skip]).Name()
}

func Exit(code int) error {
	activeExiter.Exit(code)
	return GetBypassError(code, getCallerFunction(3))
}

func GetBypassError(code int, callerFunction string) error {
	return fmt.Errorf("exit was bypassed, code: %d, called by: %s", code, callerFunction)
}
