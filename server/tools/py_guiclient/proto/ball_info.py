import proto.player_ball_info


class BallInfo(proto.player_ball_info.PlayerBallInfo):
    def __init__(self, ball):
        proto.player_ball_info.PlayerBallInfo.myinit(self)
        self.id = ball.id
        self.type = ball.type
        self.x = ball.x
        self.y = ball.y
        
        self.isplayer = None
        
    def print(self):
        print("==============================")
        print("id = ", self.id)
        print("type = ", self.type)
        print("x = ", self.x)
        print("y = ", self.y)
        print("")