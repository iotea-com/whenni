package environment

import (
	"fmt"
	"runtime"
	"strings"
)

type BuildInfo struct {
	ServiceName string `validate:"required"`
	Version     string `validate:"required"`
	LogLevel    string `validate:"required"`
	CommitHash  string `validate:"required"`
	BuildTime   string `validate:"required"`
	Dirty       string `validate:"required"`
	Creator     string `validate:"required"`
	ExtraFields map[string]any
}

func PrintInfoBanner(info BuildInfo) string {
	env := GetFromEnvVar()

	banner := []string{
		"--------------------------------------------------------------------------------",
		fmt.Sprintf("\tService              : %s", info.ServiceName),
		fmt.Sprintf("\tEnvironment          : %s", env),
		fmt.Sprintf("\tLog Level            : %s", info.LogLevel),
		fmt.Sprintf("\tSystem Info          : %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
		fmt.Sprintf("\tVersion              : %s (%s)", info.Version, truncateCommit(info.CommitHash)),
		fmt.Sprintf("\tBuild Time           : %s", info.BuildTime),
		fmt.Sprintf("\tDirty                : %s", info.Dirty),
		fmt.Sprintf("\tBuild Creator        : %s", info.Creator),
	}

	// Add any extra fields
	for k, v := range info.ExtraFields {
		var formattedValue string
		switch val := v.(type) {
		case int:
			formattedValue = fmt.Sprintf("%d", val)
		case int64:
			formattedValue = fmt.Sprintf("%d", val)
		case float64:
			formattedValue = fmt.Sprintf("%.2f", val) // 2 decimal places
		case float32:
			formattedValue = fmt.Sprintf("%.2f", val) // 2 decimal places
		case bool:
			formattedValue = fmt.Sprintf("%t", val)
		case []string:
			formattedValue = strings.Join(val, ", ")
		case []int:
			// Convert []int to []string and join
			strSlice := make([]string, len(val))
			for i, num := range val {
				strSlice[i] = fmt.Sprintf("%d", num)
			}
			formattedValue = strings.Join(strSlice, ", ")
		default:
			formattedValue = fmt.Sprintf("%v", val)
		}
		banner = append(banner, fmt.Sprintf("\t%-20s : %s", k, formattedValue))
	}

	banner = append(banner, "--------------------------------------------------------------------------------")

	return strings.Join(banner, "\n")
}

// truncateCommit safely truncates the commit hash to 8 characters
func truncateCommit(commit string) string {
	if len(commit) > 8 {
		return commit[:8]
	}
	return commit
}
