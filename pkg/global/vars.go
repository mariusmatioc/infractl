package global

import (
	"fmt"
	"strings"
)

// Global variables. This binary is a singleton and there are no goroutines, so it's safe to use global variables

// RootFolder is the absolute path of the root folder of the Infractl tree
var RootFolder = "."
var RootParent = ""

// ForceRebuild is a flag to force the rebuild of the docker images even if no code change is detected
var ForceRebuild bool
var Backend *BackendConfig

func SetBackend(backendS string) error {
	ix := strings.Index(backendS, "/")
	if ix == -1 {
		return fmt.Errorf("invalid backend id: %s. Should be <bucket name>/<key>. <key> suggested to be <org>/<project>", backendS)
	}
	Backend = &BackendConfig{
		Bucket: backendS[:ix],
		Key:    backendS[ix+1:],
	}
	return nil
}
