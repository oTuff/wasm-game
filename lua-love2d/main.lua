function love.load()
	Player = {}
	Player.x = 400
	Player.y = 200
	Player.speed = 1
end

function love.update(dt)
	if love.keyboard.isDown("right") then
		Player.x = Player.x + Player.speed + 60 * dt
	end
	if love.keyboard.isDown("left") then
		Player.x = Player.x - Player.speed + 60 * dt
	end
	if love.keyboard.isDown("up") then
		Player.y = Player.y - Player.speed + 60 * dt
	end
	if love.keyboard.isDown("down") then
		Player.y = Player.y + Player.speed + 60 * dt
	end
end

function love.draw()
	love.graphics.print("Hello World", 400, 300)

	love.graphics.setColor(255, 0, 0) -- color of rectangle
	love.graphics.rectangle("fill", Player.x, Player.y, 10, 10)
	love.graphics.setColor(255, 255, 255) -- reset color

	love.graphics.print("Current FPS: " .. tostring(love.timer.getFPS()), 10, 10)
end
