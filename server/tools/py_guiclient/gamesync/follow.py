
import time

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

        #old
        # v.x += v.vx*detal/50
        # v.y += v.vy*detal/50
        #old

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