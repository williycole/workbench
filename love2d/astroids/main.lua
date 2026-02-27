function love.load()
  GREET = "Hello Astroids!!"

  local num_astroids = 10
  Astroids = SpawnAstroids(num_astroids)
end

function love.draw()
  love.graphics.print(GREET, 0, 0)
  for i = 1, #Astroids do
    love.graphics.circle("line", Astroids[i].x, Astroids[i].y, Astroids[i].radius)
  end
end

function love.update(dt)
  for i = 1, #Astroids do
    Astroids[i].x = Astroids[i].x + Astroids[i].speed * dt
    Astroids[i].y = Astroids[i].y + Astroids[i].speed * dt
  end
end

function SpawnAstroids(num_astroids)
  Astroids = {}
  for i = 1, num_astroids do
    Astroids[i] = {
      x = math.random(0, love.graphics.getWidth()),
      y = math.random(0, love.graphics.getHeight()),
      radius = math.random(10, 50),
      speed = math.random(1, 25)
    }

    print("Astroid spawned at X: " .. Astroids[i].x
      .. " Y: " .. Astroids[i].y
      .. " Radius: " .. Astroids[i].radius
      .. " Speed: " .. Astroids[i].speed)
  end

  return Astroids
end
