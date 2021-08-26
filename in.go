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
  "flag"
  "time"
  "strings"
)

var (
  BarLen = flag.Int("length", 40,
    "The length of the progress bar (including the end chars) in characters")
  Step = flag.Duration("step", time.Second,
    "The amount of time between display updates")
)

const (
  BarFilled = '='
  BarEmpty = '-'
)

const ClearLine = "\033[K" // ESC[K, the ECMA-48 CSI code for Erase line. See man 4 console_codes

func progressbar(elapsed, total int64) string {
  realBarLen := *BarLen-2 // to account for the [ and ]

  // Solve for progress: total / BarLen = elapsed / progress
  progress := int(elapsed / (total / int64(realBarLen)))
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
  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %v [flags] [duration]\nEx: %v 2m30s && run_a_command\n\n",
      os.Args[0], os.Args[0])
    fmt.Fprintf(os.Stderr, "Flags:\n")
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
  start := time.Now()
  for elapsed := 0*time.Second; elapsed < sleeptime; elapsed = time.Since(start) {
    fmt.Fprintf(os.Stderr, "\r%s %v/%v%s",
      progressbar(elapsed.Milliseconds(), sleeptime.Milliseconds()),
      elapsed.Truncate(*Step), sleeptime, ClearLine)
    if remaining := sleeptime - elapsed; remaining < *Step {
      time.Sleep(remaining)
    } else {
      time.Sleep(*Step)
    }
  }
  fmt.Fprintf(os.Stderr, "\rDing!%s\n", ClearLine)
}
