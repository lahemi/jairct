package main

import (
	"fmt"
	"os"
	"time"
)

func clock() string {
	h, m, s := time.Now().Clock()
	return fmt.Sprintf("[%d %d %d]", h, m, s)
}

func stdout(str ...interface{}) {
	fmt.Fprintf(os.Stdout, "%s %v\n", clock(), str)
}

func stderr(str ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %v\n", clock(), str)
}

func die(str ...interface{}) {
	stderr(str)
	os.Exit(1)
}
