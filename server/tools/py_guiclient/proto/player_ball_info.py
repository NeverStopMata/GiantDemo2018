import time

class PlayerBallInfo():
    def __init__(self, ball):
        self.id = ball.id
        self.level = ball.level
        self.hp = ball.hp
        self.curmp = 0
        self.x = ball.x
        self.y = ball.y
        self.state = 0
        
        self.angle = 0
        self.face = 0
        
        self.vx = 0
        self.vy = 0

        self.state = 0

        # self.this_svr_time = 0 #这一次处理的包的服务器时间
        self.this_cli_time = 0 #这一次处理的包的客户端时间
        self.this_x = 0
        self.this_y = 0
        self.this_vx = 0
        self.this_vy = 0
        self.last_update_time = int(round(time.time() * 1000))
        self.netx = 0
        self.nety = 0


        
        self.server_pre_x = 0
        self.server_pre_y = 0
        self.client_pre_x = 0
        self.client_pre_y = 0
        self.server_now_x = 0
        self.server_now_y = 0
        self.client_now_x = 0
        self.client_now_y = 0
        
        self.isplayer = True
        self.curexp = 0
        
        
    def myinit(self):
        self.id = 0
        self.level = 0
        self.hp = 0
        self.curmp = 0
        self.x = 0
        self.y = 0
        self.state = 0
        
        self.angle = 0
        self.face = 0
        
        self.vx = 0
        self.vy = 0

        self.state = 0

        self.this_cli_time = 0
        self.this_x = 0
        self.this_y = 0
        self.this_vx = 0
        self.this_vy = 0
        self.last_update_time = int(round(time.time() * 1000))
        
        self.server_pre_x = 0
        self.server_pre_y = 0
        self.client_pre_x = 0
        self.client_pre_y = 0
        self.server_now_x = 0
        self.server_now_y = 0
        self.client_now_x = 0
        self.client_now_y = 0

        self.isplayer = False
        
    def print(self):
        print("==============================")
        print("id = ", self.id)
        print("level = ", self.level)
        print("hp = ", self.hp)
        print("x = ", self.x)
        print("y = ", self.y)
        print("")