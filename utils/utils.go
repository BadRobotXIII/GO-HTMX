// Generic utilities
package utils

import (
	"log"
	"os/exec"
)

func OpenBrowser(url string) {
	err := exec.Command("rundll32", url).Start()
	if err != nil {
		log.Fatal(err)
	}
}
