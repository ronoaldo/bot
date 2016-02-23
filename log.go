// Copyright 2015 Ronoaldo JLP <ronoaldo@gmail.com>
// Licensed under the Apache License, Version 2.0

package bot

import (
	"log"
	"os"
)

func logf(level string, msg string, args ...interface{}) {
	log.Printf(level+":"+msg, args...)
}

func debugf(msg string, args ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		logf("DEBUG", msg, args...)
	}
}
