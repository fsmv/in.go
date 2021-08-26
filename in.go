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
  "bytes"
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

var bar bytes.Buffer

func progressbar(elapsed, total float64) []byte {
  realBarLen := *BarLen-2 // to account for the [ and ]
  if realBarLen < 1 {
    return nil // The user asked for no progress bar
  }
  progress := int((elapsed / total) * float64(realBarLen))
  if progress < 0 || progress > realBarLen {
    fmt.Fprintln(os.Stderr, "Bad miscaculation:", progress, elapsed, total, realBarLen)
    os.Exit(2)
  }
  bar.Reset()
  bar.WriteRune('[')
  for i := 0; i < progress; i+=1 {
    bar.WriteRune(BarFilled)
  }
  for i := 0; i < (realBarLen - progress); i+=1 {
    bar.WriteRune(BarEmpty)
  }
  bar.WriteRune(']')
  bar.WriteRune(' ') // the space is here so we don't print it when realBarLen < 1
  return bar.Bytes()
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
    fmt.Fprintf(os.Stderr, "\r%s%v/%v%s",
      progressbar(elapsed.Seconds(), sleeptime.Seconds()),
      elapsed.Truncate(*Step), sleeptime, ClearLine)
    if remaining := sleeptime - elapsed; remaining < *Step {
      time.Sleep(remaining)
    } else {
      time.Sleep(*Step)
    }
  }
  fmt.Fprintf(os.Stderr, "\rDing!%s\n", ClearLine)
}
