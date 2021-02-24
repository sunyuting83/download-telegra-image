package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Print(strings.IndexAny("Hello,世界 !", " "))
}
