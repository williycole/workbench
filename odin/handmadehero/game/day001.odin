package main
// Handmade Heroe with Odin and SDL2
//
// I decided to do the handmade hero course: https://youtu.be/I5fNrmQYeuI?si=Z8FxQmT-a6QYMBrR
// Reason being is most of the programmers I admire have done or partaken in this course.
// That said, I've already taken a c course through boot.dev, I have a linux machine, and time is short.
// For that reason I am using -> https://davidgow.net/handmadepenguin/ to follow along but work with linux.
// Furthermore I want to actually make some games, like quickly before the end of 2026 and again time is short.
// So I will also be using odin for this reason, it has sdl2 baked in as well as raylib. I will eventually use
// raylib for the actual games, but for now I will be using sdl2 for the handmade hero course as its a bit more low level.

import "core:fmt"
import "vendor:sdl2"

z: string : "Hello Hero!"
main :: proc() {
	fmt.println(z)
	popup := sdl2.ShowSimpleMessageBox(
		sdl2.MESSAGEBOX_INFORMATION,
		"Handmade Hero",
		"This is Handmade Hero",
		nil,
	)
}
