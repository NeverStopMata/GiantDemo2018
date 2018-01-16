
import time
import math

class FollowSync:
    def __init__(self):
        pass

    def SyncMove(self, ball, m, svrTime, matchTime):
        if ball.this_vx < 1 and ball.this_vy < 1 and (m.vx > 1 or m.vy > 1):
            ball.last_update_time = int(time.time() * 1000)
        ball.this_vx = m.vx
        ball.this_vy = m.vy
        ball.this_x = m.x
        ball.this_y = m.y
        ball.netx = m.x
        ball.nety = m.y

    def UpdateMove(self, ball, mapsize, now, roomWindow):
        if int(ball.x) == int(ball.netx) and int(ball.y) == int(ball.nety):
            ball.last_update_time = now
            return
        if ball.isplayer == True:
            s = roomWindow.res.animal[ball.level]["scale"]
            radius = int(s * roomWindow.NetScale() * roomWindow.camera.scale)
        else:
            radius = 0.25 * roomWindow.NetScale() * roomWindow.camera.scale  #TODO:不是很重要，暂时硬编码

        passTime = now - ball.last_update_time
        ball.netx += ball.this_vx * passTime / 50 / 50
        ball.nety += ball.this_vy * passTime / 50 / 50
        dis = self.Dis(ball.x, ball.y, ball.netx, ball.nety)
        if dis > 100:
            self.Lerp(ball, 50 * passTime / 1000)
        elif dis > 50:
            self.Lerp(ball, 15 * passTime / 1000)
        elif dis > 5:
            self.Lerp(ball, 10 * passTime / 1000)
        elif dis > 0.0001:
            self.Lerp(ball, 4 * passTime / 1000)
        else:
            ball.x = ball.netx
            ball.y = ball.nety

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

    def Dis(self, x1, y1, x2, y2):
        return math.sqrt((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2))

    def Lerp(self, ball, lerp):
        # print("lerp: ", lerp)
        if lerp > 1:
            lerp = 1
        ball.x += (ball.netx - ball.x) * lerp
        ball.y += (ball.nety - ball.y) * lerp