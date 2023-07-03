package magic

import (
	"github.com/Instantan/watcher"
)

// Call this function as the first function in your code.
// It recompiles your program after a change in a go file and starts it after stopping the previous application.
//
//	func main() {
//		magic.HotReload()
//		// your code
//	}
//
// This function is imported from github.com/Instantan/watcher
func HotReload(command ...string) {
	watcher.HotReload(command...)
}
