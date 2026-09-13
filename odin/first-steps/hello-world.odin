package main

import "core:fmt"
import rl "vendor:raylib"

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

	rl.InitWindow(800, 600, "raylib example")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		defer rl.EndDrawing()

		rl.ClearBackground(rl.RAYWHITE)
		rl.DrawText("Hellope from Odin + raylib!", 200, 260, 24, rl.DARKGRAY)
		rl.DrawCircle(400, 350, 40, rl.MAROON)
		rl.DrawRectangle(340, 420, 120, 60, rl.SKYBLUE)
	}
}

