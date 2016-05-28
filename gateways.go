package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type GatewayDb map[string]bool

func NewGatewayDb(meshif string) (gdb GatewayDb, err error) {
	file, err := os.Open(fmt.Sprintf("/sys/kernel/debug/batman_adv/%s/gateways", meshif))
	if err != nil {
		return
	}

	gdb = GatewayDb{}

	scanner := bufio.NewScanner(file)
	scanner.Scan() // Discard headers
	for scanner.Scan() {
		gdb[strings.SplitN(strings.TrimSpace(scanner.Text()), " ", 2)[0]] = true
	}
	err = scanner.Err()

	return
}
