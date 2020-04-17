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
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestCheckWithEmptyArgs(t *testing.T) {
	Check(nil, "hello")
}

func TestCheckWithByPassExit(t *testing.T) {
	// capture stdOut
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = originalStdout }()

	// setup to bypass exit calls
	BypassExiter()
	defer func() { ResetExiter() }()

	Check(fmt.Errorf("an error"), "hello")

	assert.Nil(t, w.Close())
	out, _ := ioutil.ReadAll(r)
	os.Stdout = originalStdout
	output := string(out)
	assert.True(t, strings.Contains(output, "an error"))
	assert.True(t, strings.Contains(output, "hello"))
	assert.True(t, strings.HasSuffix(output, "*** bypassing exit, code: 3 ***\n"))
}

func TestGetCallerFunctionName(t *testing.T) {
	assert.True(t, strings.HasSuffix(getCallerFunction(2), t.Name()))
}
