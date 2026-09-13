# astroids TODO

Parked. Not porting this to Godot, new game instead. Notes here so future-me
doesn't re-derive them.

## Bugs

### 1. Ship rotation is framerate-dependent (main.lua:202, :205)

```lua
Ship.angle = Ship.angle + math.pi / 180
```

No `dt`. That's 1 degree per *frame*, not per second. 60Hz = 60 deg/sec,
144Hz = 144 deg/sec. Same code, different game.

Fix:

```lua
local TURN_SPEED = math.pi -- radians/sec = 180 deg/sec

if love.keyboard.isDown("right") then
    Ship.angle = Ship.angle + TURN_SPEED * dt
end
if love.keyboard.isDown("left") then
    Ship.angle = Ship.angle - TURN_SPEED * dt
end
```

Everything else in the file already scales by `dt`. `Ship.rateOfFire - dt` is
correct, this was the one that got missed.

### 2. Bullets leak

`Ship.bullets` only shrinks on impact. Miss every shot and the table grows
forever. Despawn offscreen, or give each bullet a lifetime and count it down
with `dt`.

### 3. Typo in the commented-out ship/asteroid collision

`checkShipAstroidCollision` reads `bullet.y` where it means `ship.y`. Also
`ship.x` / `ship.y` don't exist, the fields are `x_position` / `y_position`.

## Missing for it to be a real game

- Screen wrapping. Ship and asteroids fly off forever right now.
- Asteroid splitting. Big one hit = two smaller ones, that's the whole game.
- Ship/asteroid collision + death. Function is written, just commented out.
- `Astroid` prototype table in `love.load()` is declared and never used. Delete it.
