package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// Log out a string.
func Log(s string) error {
	fmt.Printf("%v %v\n", color.BlueString("->"), color.YellowString(s))
	return nil
}

// Barrowed from:
// https://github.com/builderscon/octav/blob/master/octav/tools/rand.go
var src = rand.NewSource(time.Now().UnixNano())

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandomString creates a random string at a specified length.
func RandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}