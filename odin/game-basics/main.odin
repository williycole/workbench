package main

// import "core:fmt"
import la "core:math/linalg"
import rl "vendor:raylib"

//  NOTE: stuff i can reference
// https://github.com/karl-zylinski/odin-raylib-hot-reload-game-template
// https://zylinski.se/posts/no-engine-gamedev-using-odin-and-raylib/#going-3d

main :: proc() {
	rl.InitWindow(1280, 720, "My Odin + Raylib game")
	// load after InitWindow always
	player := rl.LoadTexture("player.png")
	player_pos: [2]f32

	// is this the game loop?
	for !rl.WindowShouldClose() {

		input: [2]f32
		if rl.IsKeyDown(.UP) {
			input.y -= 1
		}
		if rl.IsKeyDown(.DOWN) {
			input.y += 1
		}
		if rl.IsKeyDown(.LEFT) {
			input.x -= 1
		}
		if rl.IsKeyDown(.RIGHT) {
			input.x += 1
		}

		player_mov_rate: [2]f32 : 500

		player_pos += player_mov_rate * la.normalize0(input) * rl.GetFrameTime()


		rl.BeginDrawing()
		rl.ClearBackground({160, 200, 255, 255})

		rl.DrawTextureV(player, player_pos, rl.WHITE)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}

