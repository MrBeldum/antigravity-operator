package platform

import (
	"os"
	"testing"
)

func TestCheckDisplay(t *testing.T) {
	// Darwin always assumes display
	if !checkDisplay("darwin") {
		t.Error("expected darwin to return true for checkDisplay")
	}

	// Linux without DISPLAY or WAYLAND_DISPLAY
	origDisplay := os.Getenv("DISPLAY")
	origWayland := os.Getenv("WAYLAND_DISPLAY")
	defer func() {
		_ = os.Setenv("DISPLAY", origDisplay)
		_ = os.Setenv("WAYLAND_DISPLAY", origWayland)
	}()

	_ = os.Unsetenv("DISPLAY")
	_ = os.Unsetenv("WAYLAND_DISPLAY")
	if checkDisplay("linux") {
		t.Error("expected linux without DISPLAY to return false")
	}

	_ = os.Setenv("DISPLAY", ":0")
	if !checkDisplay("linux") {
		t.Error("expected linux with DISPLAY to return true")
	}

	_ = os.Unsetenv("DISPLAY")
	_ = os.Setenv("WAYLAND_DISPLAY", "wayland-0")
	if !checkDisplay("linux") {
		t.Error("expected linux with WAYLAND_DISPLAY to return true")
	}
}

func TestFindChromeBinary(t *testing.T) {
	// Should not panic on any OS
	_ = findChromeBinary("darwin")
	_ = findChromeBinary("linux")
	_ = findChromeBinary("windows")
	_ = findChromeBinary("unknown-os")
}
