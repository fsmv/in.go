// in is a commandline utility that sleeps for a specified amount of time
// while displaying a progress bar. Ex: in 2m30s && run_a_command
//
// Author: Andrew Kallmeyer
// Copyright 2021 Google LLC.
// SPDX-License-Identifier: Apache-2.0
package main

import (
  "os"
  "fmt"
  "flag"
  "time"
  "bytes"
  "strconv"
)

var (
  Step = flag.Duration("step", time.Second, "The amount of time between display updates")
  BarStyle = flag.String("style", "[=-]", "The specification for the progress bar style.\nMust be exactly 4 chars {start, filled, unfilled, end}.")
  BarLen = flag.Int("length", 40, "The length of the progress bar (including the end chars) in characters")
  durationBuffer = make([]byte, 32) // max string is "2540400h59m59.000000001"
)

const ClearLine = "\033[K" // ESC[K, the ECMA-48 CSI code for Erase line. See man 4 console_codes

func writeProgressbar(bar *bytes.Buffer, elapsed, total float64) {
  realBarLen := *BarLen-2 // to account for the beginning and end chars
  if realBarLen < 1 {
    return // The user asked for no progress bar
  }
  progress := int((elapsed / total) * float64(realBarLen))
  bar.WriteByte((*BarStyle)[0])
  for i := 0; i < progress; i+=1 {
    bar.WriteByte((*BarStyle)[1])
  }
  for i := 0; i < (realBarLen - progress); i+=1 {
    bar.WriteByte((*BarStyle)[2])
  }
  bar.Write([]byte{(*BarStyle)[3],' '}) // the space is here so we don't print it when realBarLen < 1
}

func writeTimeUnit(out *bytes.Buffer, val *int64, unit time.Duration, label byte) {
  if u := *val / unit.Nanoseconds(); u >= 1 {
    durationBuffer = strconv.AppendInt(durationBuffer, u, 10)
    durationBuffer = append(durationBuffer, label)
    *val -= u * unit.Nanoseconds()
  }
}

// time.Duration.String() allocates the result on the heap, so I reimplemented to avoid allocation
func writeDuration(out *bytes.Buffer, d time.Duration) {
  durationBuffer = durationBuffer[:0]
  val := d.Nanoseconds()
  writeTimeUnit(out, &val, time.Hour, 'h')
  writeTimeUnit(out, &val, time.Minute, 'm')
  durationBuffer = strconv.AppendFloat(durationBuffer, float64(val)/1e9, 'f', -1, 64)
  out.Write(durationBuffer)
  out.WriteRune('s')
}

func main() {
  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %v [flags] [duration]\nEx: %v 2m30s && run_a_command\n\nFlags:\n", os.Args[0], os.Args[0])
    flag.PrintDefaults()
  }
  flag.Parse()
  if len(*BarStyle) != 4 {
    fmt.Fprintln(os.Stderr, "Invalid style argument, must be exactly 4 characters")
    os.Exit(1)
  }
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
  start := time.Now()
  var out bytes.Buffer
  for elapsed := 0*time.Second; elapsed < sleeptime; elapsed = time.Since(start) {
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
  fmt.Fprintf(os.Stderr, "\rDing!%s\n", ClearLine)
}
