package main

import (
	"fmt"
	"regexp"
)

func main() {
	id := "d432087a-d17e-4de1-9c2e-ae568c175244"
	id = "12345678-1234-1234-1234-123456789012"
	matched, _ := regexp.MatchString(`^.{8}-.{4}-.{4}-.{4}-.{12}$`, id)
	fmt.Println(matched)
}
