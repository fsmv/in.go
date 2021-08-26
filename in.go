// in is a commandline utility that sleeps for a specified amount of time
// while displaying a progress bar.
//
// Ex: in 2m30s && run_a_command
//
// Author: Andrew Kallmeyer
// Copyright 2021 Google LLC.
// SPDX-License-Identifier: Apache-2.0
package main

import (
  "os"
  "fmt"
  "time"
  "strings"
)

const (
  BarLen = 40
  BarFilled = '='
  BarEmpty = '-'
  Step = time.Second
)

const ClearLine = "\033[K" // ESC[K, the ECMA-48 CSI code for Erase line. See man 4 console_codes

func progressbar(elapsed, total int64) string {
  const realBarLen = BarLen-2 // to account for the [ and ]

  // Solve for progress: total / BarLen = elapsed / progress
  progress := int(elapsed / (total / realBarLen))
  if progress < 0 || progress > realBarLen {
    fmt.Fprintln(os.Stderr, "Bad miscaculation:", progress, elapsed, total, realBarLen)
    os.Exit(2)
  }

  var b strings.Builder
  b.WriteRune('[')
  for i := 0; i < progress; i+=1 {
    b.WriteRune(BarFilled)
  }
  for i := 0; i < (realBarLen - progress); i+=1 {
    b.WriteRune(BarEmpty)
  }
  b.WriteRune(']')
  return b.String()
}

func main() {
  if len(os.Args) < 2 {
    fmt.Fprintf(os.Stderr, "Usage example: %v 2m30s\n", os.Args[0])
    os.Exit(127)
  }
  sleeptime, err := time.ParseDuration(os.Args[1])
  if err != nil {
    fmt.Fprintln(os.Stderr, "Invalid duration:", os.Args[1])
    os.Exit(1)
  }
  for elapsed := 0*time.Second; elapsed < sleeptime; elapsed += Step {
    fmt.Fprintf(os.Stderr, "\r%s %v/%v%s",
      progressbar(elapsed.Milliseconds(), sleeptime.Milliseconds()),
      elapsed, sleeptime, ClearLine)
    time.Sleep(Step)
  }
  fmt.Fprintf(os.Stderr, "\rDing!%s\n", ClearLine)
}
