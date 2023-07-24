/*?sr/bin/env go run "$0" "$@"; exit $? #*/
// Author: Andrew Kallmeyer <ask@ask.systems>
// Copyright 2021-2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// in is a commandline utility that sleeps for a specified amount of time
// while displaying a progress bar. Ex: in 2m30s && run_a_command
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	Step = flag.Duration("step", time.Second,
		"The amount of time between display updates")
	BarStyle = flag.String("style", "[=-]", ""+
		"The specification for the progress bar style.\n"+
		"Must be exactly 4 chars {start, filled, unfilled, end}.")
	BarLen = flag.Int("length", 40,
		"The length of the progress bar (including the end chars) in characters")
)

var durationBuffer = make([]byte, 32) // max string is "2540400h59m59.000000001"
const ClearLine = "\033[K"            // ESC[K, the ECMA-48 CSI code for Erase line. See man 4 console_codes

func writeProgressbar(out *bytes.Buffer, elapsed, total float64) {
	if realBarLen := *BarLen - 2; realBarLen > 0 { // handle the start and end chars
		progress := int((elapsed / total) * float64(realBarLen))
		out.WriteByte((*BarStyle)[0])
		for i := 0; i < progress; i += 1 {
			out.WriteByte((*BarStyle)[1])
		}
		for i := 0; i < (realBarLen - progress); i += 1 {
			out.WriteByte((*BarStyle)[2])
		}
		out.Write([]byte{(*BarStyle)[3], ' '}) // the space is here so we don't print it when realBarLen < 1
	}
}

func writeTimeUnit(out *bytes.Buffer, remainingNanos *int64, unit time.Duration, label byte) {
	if u := *remainingNanos / unit.Nanoseconds(); u >= 1 {
		durationBuffer = strconv.AppendInt(durationBuffer, u, 10)
		durationBuffer = append(durationBuffer, label)
		*remainingNanos -= u * unit.Nanoseconds()
	}
}

// time.Duration.String() allocates the result on the heap, so I reimplemented to avoid allocation
func writeDuration(out *bytes.Buffer, d time.Duration) {
	durationBuffer = durationBuffer[:0]
	remainingNanos := d.Nanoseconds()
	writeTimeUnit(out, &remainingNanos, time.Hour, 'h')
	writeTimeUnit(out, &remainingNanos, time.Minute, 'm')
	if len(durationBuffer) == 0 || remainingNanos != 0 {
		durationBuffer = strconv.AppendFloat(durationBuffer, float64(remainingNanos)/1e9, 'f', -1, 64)
		durationBuffer = append(durationBuffer, 's')
	}
	out.Write(durationBuffer)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, ""+
			"Usage: %v [flags] [duration]\n"+
			"Ex: %v 2m30s && run_a_command\n\n"+
			"Flags:\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(127)
	}
	sleeptime, err := time.ParseDuration(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid duration:", args[0])
		os.Exit(1)
	}
	if len(*BarStyle) != 4 {
		fmt.Fprintln(os.Stderr, "Invalid style argument, must be exactly 4 characters")
		os.Exit(1)
	}
	start := time.Now()
	var out bytes.Buffer
	for elapsed := 0 * time.Second; elapsed < sleeptime; elapsed = time.Since(start) {
		out.WriteRune('\r')
		writeProgressbar(&out, elapsed.Seconds(), sleeptime.Seconds())
		writeDuration(&out, elapsed.Truncate(*Step))
		out.WriteRune('/')
		writeDuration(&out, sleeptime)
		out.WriteString(ClearLine)
		os.Stderr.Write(out.Bytes())
		out.Reset()
		if remaining := sleeptime - elapsed; remaining < *Step {
			time.Sleep(remaining)
		} else {
			time.Sleep(*Step)
		}
	}
	fmt.Fprintf(os.Stderr, "\r%s", ClearLine)
}
