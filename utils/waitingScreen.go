package utils

import (
	"fmt"
	"time"
)
func WaitingScreen(done chan bool, desc string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from error:", r)
				done <- true
			}
		}()

		fmt.Print("\033[2K\r")

		for {
			select {
			case <-done:
				fmt.Print("\033[2K\r") 
				return
			default:
				for _, r := range `-\|/` {
					fmt.Printf("\r%c %s", r, desc)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()
}
