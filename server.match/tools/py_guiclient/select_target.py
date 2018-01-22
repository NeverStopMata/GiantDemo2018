
'''
选择屏幕范围内最近的玩家
'''
def near_target(user, w, h):
    myball = user.get_ball()
    if myball == None:
        return
        
    xbegin = myball.x - w/2
    xend = myball.x + w/2
    ybegin = myball.y - h/2
    yend = myball.y + h/2
    
    target = 0
    mindistance = 9999999
    user.mutex.acquire()
    for k, v in user.playerballs.items():
        if v.id == user.ball_id:
            continue
        if v.x <= xbegin or v.x >= xend or v.y <= ybegin or v.y >= yend:
            continue
        
        tmpd = (v.x - myball.x) * (v.x - myball.x) + (v.y - myball.y) * (v.y - myball.y)
        if tmpd < mindistance:
            target = k
            mindistance = tmpd        
    user.mutex.release()
    user.face = target