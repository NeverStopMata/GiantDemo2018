
import time
import math

class FollowSync:
    def __init__(self):
        pass

    def SyncMove(self, ball, m, svrTime, matchTime):
        if svrTime == None:
            return
            
        if ball.this_vx < 1 and ball.this_vy < 1 and (m.vx > 1 or m.vy > 1):
            ball.last_update_time = int(time.time() * 1000)

        pastTime = matchTime.GetTimeDelta(ball.last_update_time, svrTime)
        resTime = 100 - pastTime
        if resTime == 0:
            resTime = 1

        ball.vx = (m.x - ball.x) * 50 / resTime
        ball.vy = (m.y - ball.y) * 50 / resTime

        ball.this_cli_time = matchTime.ToClientTime(svrTime) + 100
        ball.this_vx = m.vx
        ball.this_vy = m.vy
        ball.this_x = m.x
        ball.this_y = m.y

        if abs(m.x - ball.x) < 2:
            ball.vx = 0
            ball.this_vx = 0
        if abs(m.y - ball.y) < 2:
            ball.vy = 0
            ball.this_vy = 0

    def UpdateMove(self, ball, mapsize, now, roomWindow):
        if ball.isplayer == True:
            s = roomWindow.res.animal[ball.level]["scale"]
            radius = int(s * roomWindow.NetScale() * roomWindow.camera.scale)
        else:
            radius = 0.25 * roomWindow.NetScale() * roomWindow.camera.scale  #TODO:不是很重要，暂时硬编码

        passTime = now - ball.last_update_time
        if now < ball.this_cli_time:
            ball.x += ball.vx * passTime / 50
            ball.y += ball.vy * passTime / 50
        elif ball.last_update_time >= ball.this_cli_time:
            ball.x += ball.this_vx * passTime / 50
            ball.y += ball.this_vy * passTime / 50
        else:
            t1 = ball.this_cli_time - ball.last_update_time
            ball.x += ball.vx * t1 / 50
            ball.y += ball.vy * t1 / 50

            t2 = now - ball.this_cli_time
            ball.x += ball.this_vx * t2 / 50
            ball.y += ball.this_vy * t2 / 50

        for k, v in roomWindow.user.balls.items():
            if v.type == 21:
                (ok, xd, yd) = self.CheckCollision(ball.x, ball.y, radius, v, roomWindow)
                if ok:
                    ball.x += xd
                    ball.y += yd

        for k, v in roomWindow.user.playerballs.items():
            if v.id != ball.id:
                (ok, xd, yd) = self.CheckOtherBallCollision(ball.x, ball.y, radius, v, roomWindow)
                if ok:
                    ball.x += xd
                    ball.y += yd

        for block in roomWindow.res.map["nodes"]:
            (ok, xd, yd) = self.CheckRectangle(ball, radius, block, roomWindow)
            if ok:
                ball.x += xd
                ball.y += yd

        if ball.x > mapsize - radius:
            ball.x = mapsize - radius
        elif ball.x < radius:
            ball.x = radius
        if ball.y > mapsize - radius:
            ball.y = mapsize - radius
        elif ball.y < radius:
            ball.y = radius

        if abs(ball.x - ball.this_x) < 1:
            ball.x = ball.this_x
        if abs(ball.y - ball.this_y) < 1:
            ball.y = ball.this_y

        ball.last_update_time = now

    def CheckCollision(self, x, y, radius, ball, roomWindow):
        ballRadius = int(roomWindow.res.food[ball.type] * roomWindow.NetScale() * roomWindow.camera.scale)
        dis = math.sqrt((x - ball.x) * (x - ball.x) + (y - ball.y)*(y - ball.y))
        if dis >= radius + ballRadius:
            return (False, 0, 0)
        delta = abs(dis - radius - ballRadius)
        xdelta = (x - ball.x) * delta / dis
        ydelta = (y - ball.y) * delta / dis
        return (True, xdelta, ydelta)

    def CheckOtherBallCollision(self, x, y, radius, ball, roomWindow):
        s = roomWindow.res.animal[ball.level]["scale"]
        ballRadius = int(s * roomWindow.NetScale() * roomWindow.camera.scale)
        dis = math.sqrt((x - ball.x) * (x - ball.x) + (y - ball.y)*(y - ball.y))
        if dis >= radius + ballRadius:
            return (False, 0, 0)
        delta = abs(dis - radius - ballRadius)
        xdelta = (x - ball.x) * delta / dis
        ydelta = (y - ball.y) * delta / dis
        return (True, xdelta, ydelta)

    def CheckRectangle(self, ball, radius, block, roomWindow):
        w = int(block["radius"]* roomWindow.NetScale() * roomWindow.camera.scale)
        bx = block["px"] * roomWindow.NetScale() * roomWindow.camera.scale
        by = block["py"] * roomWindow.NetScale() * roomWindow.camera.scale
        if abs(bx - ball.x) >= radius + w or abs(by - ball.y) >= radius + w:
            return (False, 0, 0)

        dx = radius + w - abs(bx - ball.x) 
        dy = radius + w - abs(by - ball.y)
        xdelta = 0
        ydelta = 0

        if abs(ball.this_vy) < 0.1 or (abs(ball.vy) >= abs(ball.vx) and dx < dy):
            if ball.x > bx:
                xdelta = dx
            else:
                xdelta = -dx
        elif abs(ball.this_vx) < 0.1 or (abs(ball.vx) >= abs(ball.vy) and dy < dx):
            if ball.y > by:
                ydelta = dy
            else:
                ydelta = -dy
        return (True, xdelta, ydelta)

