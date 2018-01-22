
import time

class FollowSync:
    def __init__(self):
        pass

    def SyncMove(self, ball, m, svrTime, matchTime):
        ball.vx = (m.x - ball.x) * 50 / 100
        ball.vy = (m.y - ball.y) * 50 / 100
        ball.this_x = m.x
        ball.this_y = m.y

    def UpdateMove(self, ball, mapsize, now, roomWindow):
    
        if ball.isplayer == True:
            s = roomWindow.res.animal[ball.level]["scale"]
            radius = int(s * roomWindow.NetScale() * roomWindow.camera.scale)
        else:
            radius = 0.25 * roomWindow.NetScale() * roomWindow.camera.scale  #TODO:不是很重要，暂时硬编码
            
        passTime = now - ball.last_update_time
        ball.x += ball.vx*passTime/50
        ball.y += ball.vy*passTime/50

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