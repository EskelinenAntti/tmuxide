package herdr

import "os"

// InSession reports whether the current process is running inside Herdr.
func InSession() bool {
	return os.Getenv("HERDR_ENV") == "1"
}
