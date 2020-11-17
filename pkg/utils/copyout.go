package utils

import (
	"bufio"
	"fmt"
	"io"
)

func CopyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
