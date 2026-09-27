package version

// Build-time variables set via -ldflags.
var (
	Version   = "0.1.0-dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Info returns version metadata for API and UI surfaces.
func Info() map[string]string {
	return map[string]string{
		"version":    Version,
		"commit":     Commit,
		"build_date": BuildDate,
		"product":    "Yggdrasil",
	}
}
