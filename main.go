/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package main

import (
	"github.com/DirusK/utils/printer"

	"authentication-chains/cmd"
)

func main() {
	printer.Colored = true
	cmd.Execute()
}
