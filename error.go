// Copyright (C) 2023 neocotic
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sets

import (
	"errors"
	"fmt"
)

// ErrJSONElementCount is returned by a fixed-size Set implementation of json.Unmarshaler when the number of
// unmarshalled elements do not meet the requirements of the Set.
var ErrJSONElementCount = errors.New("invalid number of elements unmarshalled from json")

// fmtErrJSONElementCount returns an ErrJSONElementCount formatted with the expected and actual number of elements
// unmarshalled from JSON.
func fmtErrJSONElementCount(expect, actual int) error {
	return fmt.Errorf("%w; want %v, got %v", ErrJSONElementCount, expect, actual)
}
