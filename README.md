# in.go

A short little program I wrote to display a progress bar while waiting for a
command to run.

Use it like: `in 2m30s && run_a_command`

My main use case is a server that has a BIOS that is slow to reboot and I want to
ssh back in when it's up.

## Installing

Run `GOGC=off go install ask.systems/in@latest`

Disabling the garbage collector is optional, but it's fun to do it, since this program
only allocates a constant amount of data. If you don't disable it, it will probably still
never be triggered to run.

Alternatively I have added a fancy go shebang line that I adapted so that
`go fmt` wouldn't break it. So you can actually just download the in.go file,
`chmod +x in.go` and run it like a script `./in.go 10s` (as long as you have go
installed).

Run `wget ask.systems/in.go && chmod +x in.go` to get it.

## License

Released under the Apache 2.0 License, which is a permissive license so use and
fork as desired. While the copyright is owned by Google (because I work there).

This is not an official Google product, it is my personal project. Google
disclaims all warranties as to its quality, merchantability, or fitness for a
particular purpose.

## Contributing

Please do not contribute code, only issues. In order for me to accept
contributions you would need to sign the Google CLA and I don't think it is
worth it for this.

Of course, you can make changes in your own repo without the CLA. I like the idea of people customizing it to suit their needs instead of making one over-general thing that does everything.
