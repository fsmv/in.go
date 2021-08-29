# in.go

A short little program I wrote to display a progress bar while waiting for a
command to run.

Use it like: `in 2m30s && run_a_command`

My main use case is a server that has a BIOS that is slow to reboot and I want to
ssh back in when it's up.

## Installing

Just run `GOGC=off go build in.go` and move the `in` binary into ~/bin or ~/.local/bin

Disabling the garbage collector is optional, but it's fun to do it, since this program
only allocates a constant amount of data. If you don't disable it, it will probably still
never be triggered to run.

## License

Released under the Apache 2.0 License, which is a permissive license so use and
fork as desired. Technically Google owns the copyright since I work there and I
decided to release this under the "below 100 lines" policy that allows me to use
the Apache license without explicit approval.
