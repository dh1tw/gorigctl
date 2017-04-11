package events

import (
	"bufio"
	"os"
	"strings"

	"github.com/cskr/pubsub"
)

func CaptureKeyboard(evPS *pubsub.PubSub) {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		if scanner.Scan() {
			switch scanner.Text() {
			default:
				evPS.Pub(strings.Fields(scanner.Text()), CliInput)
			}
		}
	}
}
