//go:build wasm || js
// +build wasm js

package tty

// NewDevTty obtains a default tty from the console or TTY (e.g. /dev/tty) for the process.
func NewDevTty() (Tty, error) { _ = "STUB: not implemented"; return *new(Tty), nil }

// NewDevTtyFromDev obtains a tty from the given device path. Not supported on Windows.
func NewDevTtyFromDev(dev string) (Tty, error) { _ = "STUB: not implemented"; return *new(Tty), nil }

// NewStdIoTty obtains a tty from stdin and stdout.
func NewStdIoTty() (Tty, error) { _ = "STUB: not implemented"; return *new(Tty), nil }
