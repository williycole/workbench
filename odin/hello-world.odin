package main

import "core:fmt"

a :: "I'm a const"
b: string : "also const"
x, y: int
z: string = "Hellope!"

main :: proc() {
	fmt.println(z)
	x, y, z := 1, 1, "Odin"
	fmt.printf("y: %d, x: %d, z: %s, a: %s, b: %s\n", x, y, z, a, b)
	x, y, z = 2, 2, "Rules"
	fmt.printf("y: %d, x: %d, z: %s, a: %s, b: %s\n", x, y, z, a, b)
}
