package log

import "fmt"

func panicOnErr(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%s: %w", msg, err))
	}
}
