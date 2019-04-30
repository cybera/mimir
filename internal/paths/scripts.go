package paths

import "path/filepath"

func Scripts() string {
	return filepath.Join(ProjectRoot(), "src/scripts")
}
