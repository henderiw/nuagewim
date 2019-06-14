package nuagewim

import (
	"fmt"
	"os"

	"github.com/nuagenetworks/go-bambou/bambou"
)

func handleError(err *bambou.Error, t string, o string) {
	if err != nil {
		fmt.Println("Unable to " + o + " \"" + t + "\": " + err.Description)
		os.Exit(1)
	}
}
