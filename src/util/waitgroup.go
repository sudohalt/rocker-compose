/*-
 * Copyright 2015 Grammarly, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

// TODO: document and write tests

import (
	"fmt"
	"time"
)

// ErrorWaitGroup is a wait group that returns the first error occured
type ErrorWaitGroup struct {
	ch chan error
}

// NewErrorWaitGroup makes a new ErrorWaitGroup of size n
func NewErrorWaitGroup(size int) *ErrorWaitGroup {
	return &ErrorWaitGroup{
		ch: make(chan error, size),
	}
}

// Done pops a number from the wait group
func (wg *ErrorWaitGroup) Done(err error) {
	wg.ch <- err
	return
}

// Wait blogs until the wait group finished and returns the error if any
func (wg *ErrorWaitGroup) Wait() (err error) {
	n := cap(wg.ch)
	if n == 0 {
		return nil
	}
	for {
		if resErr := <-wg.ch; resErr != nil && err == nil {
			err = resErr
		}
		if n--; n == 0 {
			break
		}
	}
	return err
}

// WaitFor same as Wait but with the timeout
func (wg *ErrorWaitGroup) WaitFor(timeout time.Duration) error {
	n := cap(wg.ch)
	if n == 0 {
		return nil
	}
	for {
		select {
		case err := <-wg.ch:
			if err != nil {
				return err
			}
		case <-time.After(timeout):
			return fmt.Errorf("timeout %s", timeout)
		}
		if n--; n == 0 {
			break
		}
	}
	return nil
}
