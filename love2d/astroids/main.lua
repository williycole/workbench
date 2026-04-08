-- TODO: do my natural way then play around with fat structs for optimization
-- Helper func for managing bullet and astroid collisions
-- if any bullets xys are the same as a bullets, then turn the astroid red

-- NOTE: Everything moving in 2D is the same pattern:
--   velocity_x = cos(angle) * speed
--   velocity_y = sin(angle) * speed
--   position_x += velocity_x * dt
--   position_y += velocity_y * dt
-- Asteroid, bullet, ship thrust - all of it. Learn this reflex.

-- TODO: When stuck on "how do I move this thing", name what it needs FIRST:
--   where does it start? what direction? how fast?
--   Then the code writes itself. Don't reach for the API docs first.

-- NOTE: initial state values
function love.load()
	GameScoreState = {
		ScoreLimit = 15,
		Score = 0,
		GameMessage = "SCORE: ",
	}

	Astroid = {
		x_position = nil,
		y_position = nil,
		angle = nil,
		speed = nil,
		x_velocity = nil,
		y_velocity = nil,
		radius = nil,
	}

	Astroids = SpawnAstroids(10)

	COOLDOWN_TIMER = 0.5
	Ship = {
		x_position = 400,
		y_position = 300,
		angle = -math.pi / 2,
		thrusters_acceleration = 100,
		y_velocity = 0,
		x_velocity = 0,
		length = 20,
		width = 10,
		gun = {
			angle = 0,
			x_postion = 0,
			y_postion = 0,
		},
		bullets = {},
		rateOfFire = COOLDOWN_TIMER,
	}
end

--NOTE: For camera and inital drawings
function love.draw()
	love.graphics.print(GameScoreState.GameMessage .. GameScoreState.Score, 0, 0)
	DrawAstroids()
	DrawShip()
	DrawBullets()
end

-- helper func for bullet collisons
local function checkBulletAstroidCollision(bullet, asteroid)
	local dx = bullet.x - asteroid.x_position
	local dy = bullet.y - asteroid.y_position
	local distance = math.sqrt(dx * dx + dy * dy)
	return distance < asteroid.radius
end

-- helper func for astroid collisons
-- local function checkShipAstroidCollision(ship, asteroid)
-- 	local dx = ship.x - asteroid.x_position
-- 	local dy = bullet.y - asteroid.y_position
-- 	local distance = math.sqrt(dx * dx + dy * dy)
-- 	return distance < asteroid.radius
-- end

-- NOTE: For game state changes
function love.update(dt)
	UpdateAstroidsPosition(dt)
	UpdateShipPosition(dt)
	UpdateBulletsPosition(dt)

	Ship.rateOfFire = Ship.rateOfFire - dt

	for i = #Ship.bullets, 1, -1 do
		for j = #Astroids, 1, -1 do
			if checkBulletAstroidCollision(Ship.bullets[i], Astroids[j]) then
				table.remove(Ship.bullets, i)
				table.remove(Astroids, j)
				GameScoreState.Score = GameScoreState.Score + 1
				break
			end
		end
	end

	-- Add score check here also to win after a certain score by
	-- stopping astroid spawning
	if #Astroids <= 5 and GameScoreState.Score < GameScoreState.ScoreLimit then
		Astroids = SpawnAstroids(10)
	end

	if GameScoreState.Score >= GameScoreState.ScoreLimit then
		GameScoreState.GameMessage = "YOU WIN!!!"
	end
end

function SpawnAstroids(num_astroids)
	Astroids = {}
	for i = 1, num_astroids do
		local speed = math.random(1, 50)
		local angle = math.random() * 2 * math.pi -- 2pi ~ 6.28, i.e. full circle in radians

		Astroids[i] = {
			x_position = math.random(0, love.graphics.getWidth()),
			y_position = math.random(0, love.graphics.getHeight()),
			radius = math.random(10, 50),
			x_velocity = math.cos(angle) * speed,
			y_velocity = math.sin(angle) * speed,
		}
		-- math.sin(angle) * speed
		print(
			"Astroid spawned at"
				.. " X Position: "
				.. Astroids[i].x_position
				.. " Y Position: "
				.. Astroids[i].y_position
				.. " Radius: "
				.. Astroids[i].radius
				.. " X Velocity: "
				.. Astroids[i].x_velocity
				.. " Y Velocity: "
				.. Astroids[i].y_velocity
		)
	end

	return Astroids
end

function DrawAstroids()
	for i = 1, #Astroids do
		love.graphics.circle("line", Astroids[i].x_position, Astroids[i].y_position, Astroids[i].radius)
	end
end

function UpdateAstroidsPosition(dt)
	for i = 1, #Astroids do
		Astroids[i].x_position = Astroids[i].x_position + Astroids[i].x_velocity * dt
		Astroids[i].y_position = Astroids[i].y_position + Astroids[i].y_velocity * dt
	end
	return Astroids
end

function DrawShip()
	local mode = "fill"
	love.graphics.push()
	love.graphics.translate(Ship.x_position, Ship.y_position)
	love.graphics.rotate(Ship.angle)
	love.graphics.polygon(mode, -Ship.length / 2, -Ship.width / 2, -Ship.length / 2, Ship.width / 2, Ship.length / 2, 0)
	love.graphics.pop()
end

local function getFrontPosition(ship)
	local frontX = ship.x_position + (ship.length / 2) * math.cos(ship.angle)
	local frontY = ship.y_position + (ship.length / 2) * math.sin(ship.angle)
	return frontX, frontY
end

function UpdateShipPosition(dt)
	-- update gun postions
	Ship.gun.x_postion, Ship.gun.y_postion = getFrontPosition(Ship)

	if love.keyboard.isDown("space") and Ship.rateOfFire <= 0 then
		local newBullet = {
			x = Ship.gun.x_postion,
			y = Ship.gun.y_postion,
			angle = Ship.angle,
			speed = 100,
			x_velocity = math.cos(Ship.angle) * 100,
			y_velocity = math.sin(Ship.angle) * 100,
		}
		table.insert(Ship.bullets, newBullet)
		Ship.rateOfFire = COOLDOWN_TIMER
	end

	-- Increase velocity on up key press
	if love.keyboard.isDown("up") then
		Ship.x_velocity = Ship.x_velocity + Ship.thrusters_acceleration * math.cos(Ship.angle) * dt
		Ship.y_velocity = Ship.y_velocity + Ship.thrusters_acceleration * math.sin(Ship.angle) * dt
	end

	-- Decrease velocity on down key press
	if love.keyboard.isDown("down") then
		Ship.x_velocity = Ship.x_velocity - Ship.thrusters_acceleration * math.cos(Ship.angle) * dt
		Ship.y_velocity = Ship.y_velocity - Ship.thrusters_acceleration * math.sin(Ship.angle) * dt
	end

	-- Rotate ship angle
	if love.keyboard.isDown("right") then
		Ship.angle = Ship.angle + math.pi / 180
	end
	if love.keyboard.isDown("left") then
		Ship.angle = Ship.angle - math.pi / 180
	end

	-- Update ship postion
	Ship.x_position = Ship.x_position + Ship.x_velocity * dt
	Ship.y_position = Ship.y_position + Ship.y_velocity * dt
end

-- A bullet is just an asteroid with a known starting angle (ship.angle).
-- Took ~2hrs and outside help to see this. It was already written above in SpawnAstroids.
function UpdateBulletsPosition(dt)
	for _, bullet in ipairs(Ship.bullets) do
		bullet.x = bullet.x + bullet.x_velocity * dt
		bullet.y = bullet.y + bullet.y_velocity * dt
	end
end

function DrawBullets()
	for i, bullet in ipairs(Ship.bullets) do
		love.graphics.rectangle("fill", bullet.x, bullet.y, 3, 3)
	end
end
