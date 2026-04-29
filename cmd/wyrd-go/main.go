// SPDX-License-Identifier: AGPL-3.0-or-later

// Command wyrd-go is the eventual server binary for OpenWyrd MOP.
// Currently a stub — see MIGRATION-PLAN.md for the work to flesh this out.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("wyrd-go 0.0.0-dev — pre-alpha, no relay yet")
		return
	}
	fmt.Fprintln(os.Stderr, "wyrd-go: relay not implemented — see MIGRATION-PLAN.md")
	os.Exit(2)
}
