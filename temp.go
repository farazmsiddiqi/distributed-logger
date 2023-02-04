package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	var t_now float64 = float64(time.Now().UnixNano())

	fmt.Fprintln(os.Stdout, t_now / 1000000000, "-", "testes", "connected")
	fmt.Printf("%f\n", t_now / 1000000000.0)
	fmt.Printf("%f\n", t_now)

	fmt.Printf("%f - %s connected", float64(time.Now().UnixNano()) / 1000000000.0, "sdfdsfa")

}