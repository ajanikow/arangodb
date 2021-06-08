//
// DISCLAIMER
//
// Copyright 2021 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Adam Janikowski
//

package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type TimeoutFunc func() error

func NewTimeoutFunc(f func() error) TimeoutFunc {
	return f
}

func (f TimeoutFunc) Execute(timeout, interval time.Duration) error {
	if err := f(); err != nil {
		if IsInterrupt(err) {
			return nil
		}

		return err
	}

	timeoutT := time.NewTimer(timeout)
	defer timeoutT.Stop()

	intervalT := time.NewTicker(interval)
	defer intervalT.Stop()

	for {
		select {
		case <-timeoutT.C:
			return fmt.Errorf("timeout")
		case <-intervalT.C:
			if err := f(); err != nil {
				if IsInterrupt(err) {
					return nil
				}

				return err
			}
		}
	}
}

func (f TimeoutFunc) ExecuteT(t *testing.T, timeout, interval time.Duration) {
	require.NoError(t, f.Execute(timeout, interval))
}

type Interrupt struct {
}

func (i Interrupt) Error() string {
	return "interrupt"
}

func NewInterrupt() error {
	return Interrupt{}
}

func IsInterrupt(err error) bool {
	_, ok := err.(Interrupt)
	return ok
}
