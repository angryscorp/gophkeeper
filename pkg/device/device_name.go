package device

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
)

// GenerateDeviceName builds a unique device name from OS, arch,
// current username, and hostname. Spaces are replaced with dashes
// and the result is lowercased.
func GenerateDeviceName() string {
	username := "unknown"
	if currentUser, err := user.Current(); err == nil {
		username = currentUser.Username
	}

	hostname, _ := os.Hostname()

	deviceName := fmt.Sprintf("%s-%s-%s-%s",
		runtime.GOOS,
		runtime.GOARCH,
		username,
		hostname,
	)

	deviceName = strings.ReplaceAll(deviceName, " ", "-")
	deviceName = strings.ToLower(deviceName)

	return deviceName
}
