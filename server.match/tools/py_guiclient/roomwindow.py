import wx
import json
import time
import camera
import math
import copy
import gateway
import state
import vector2
import select_target
import proto.wilds_pb2

import log

print=log.log

NET_SCALE = 100

PLAY_AREA_WIDTH = 750
PLAY_AREA_HEIGHT = PLAY_AREA_WIDTH
RANKING_AREA_WIDTH = 500
BOTTOM_AREA_HEIGHT = 50

ON_TIMER_INTERVAL = 25

LEFT = ord('a') - 32
RIGHT = ord('d') - 32
UP = ord('w') - 32
DOWN = ord('s') - 32

ATTACK = ord('u') - 32
THROW_HAMMER = ord('i') - 32
THROW_BOMB = ord('o') - 32
SPRINT = ord('j') - 32
RUN = ord('k') - 32

class RoomWindow(wx.Frame):
    def __init__(self, user, args, cfg):
        wx.Frame.__init__(self, None, style = wx.CAPTION | wx.CLOSE_BOX)
        self.user = user
        self.cfg = cfg
        self.args = args
        self.timer = None
        self.pre_on_timer = 0
        self.camera = camera.Camera()
        self.res = None
        self.foodxml = None
        self.inited = False
        self.EndTime = None
        
        self.frame = 0
        self.dir = 0
        self.power = 0
        self.sendMoveTime = 0
        self.ballx = 0
        self.bally = 0

        self.show_line = True
        
        self.init()

    def init(self):
        self.SetTitle("py_guiclient")
        
        self.res = self.user.res
        
        self.SetClientSize((PLAY_AREA_WIDTH + RANKING_AREA_WIDTH, PLAY_AREA_HEIGHT + BOTTOM_AREA_HEIGHT))
        
        self.panel1 = wx.Panel(self, -1, pos=(PLAY_AREA_WIDTH + 1, 0), size=(RANKING_AREA_WIDTH, PLAY_AREA_HEIGHT/2))
        self.listTop = wx.ListBox(choices=[], parent=self.panel1, pos=wx.Point(0, 0),size=wx.Size(RANKING_AREA_WIDTH, PLAY_AREA_HEIGHT/2), style=0)
        
        self.panel2 = wx.Panel(self, -1, pos=(PLAY_AREA_WIDTH + 1, PLAY_AREA_HEIGHT/2), size=(RANKING_AREA_WIDTH, PLAY_AREA_HEIGHT/2))
        wx.StaticText(self.panel2, -1, "ID", pos=wx.Point(10, 5))
        wx.StaticText(self.panel2, -1, ": %d" % (self.user.id), pos=wx.Point(70, 5))
        wx.StaticText(self.panel2, -1, "账号", pos=wx.Point(10, 35))
        wx.StaticText(self.panel2, -1, ": %s" % self.user.account_name, pos=wx.Point(70, 35))
        wx.StaticText(self.panel2, -1, "角色名", pos=wx.Point(10, 65))
        wx.StaticText(self.panel2, -1, ": %s" % self.user.player_name, pos=wx.Point(70, 65))
        
        wx.StaticText(self.panel2, -1, "ball id", pos=wx.Point(10, 125))
        self.lblBallId = wx.StaticText(self.panel2, -1, ": %d" % self.user.ball_id, pos=wx.Point(70, 125))
        wx.StaticText(self.panel2, -1, "ball x", pos=wx.Point(10, 155))
        self.lblBallX = wx.StaticText(self.panel2, -1, ": 0", pos=wx.Point(70, 155))
        wx.StaticText(self.panel2, -1, "ball y", pos=wx.Point(10, 185))
        self.lblBallY = wx.StaticText(self.panel2, -1, ": 0", pos=wx.Point(70, 185))

        wx.StaticText(self.panel2, -1, "移动距离(c)", pos=wx.Point(10, 245))
        self.lblMoveDisC = wx.StaticText(self.panel2, -1, ": 0", pos=wx.Point(110, 245))
        wx.StaticText(self.panel2, -1, "移动距离(s)", pos=wx.Point(10, 275))
        self.lblMoveDisS = wx.StaticText(self.panel2, -1, ": 0", pos=wx.Point(110, 275))
        
        wx.StaticText(self.panel2, -1, "刷新间隔", pos=wx.Point(340, 5))
        self.lblFPS = wx.StaticText(self.panel2, -1, "", pos=(445, 5))
        
        wx.StaticText(self.panel2, -1, "UDP包大小", pos=wx.Point(340, 35))
        self.lblUDP = wx.StaticText(self.panel2, -1, "", pos=(445, 35))
        
        wx.StaticText(self.panel2, -1, "房间内人数", pos=wx.Point(340, 65))
        self.lblPlayerNum = wx.StaticText(self.panel2, -1, "", pos=(445, 65))
        
        self.btnViewNormal = wx.Button(self.panel2, -1, "正常", pos=(10, 330), size=(80, -1))
        self.Bind(wx.EVT_BUTTON, self.OnViewNormal, self.btnViewNormal)
        
        self.btnViewAll = wx.Button(self.panel2, -1, "全图", pos=(120, 330), size=(80, -1))
        self.Bind(wx.EVT_BUTTON, self.OnViewAll, self.btnViewAll)

        self.btnRelife = wx.Button(self.panel2, -1, "复活", pos=(230, 330), size=(80, -1))
        self.Bind(wx.EVT_BUTTON, self.OnRelife, self.btnRelife)
        
        self.btnExitRoom = wx.Button(self.panel2, -1, "退出", pos=(340, 330), size=(80, -1))
        self.Bind(wx.EVT_BUTTON, self.OnExitRoom, self.btnExitRoom)
        
        self.panel3 = wx.Panel(self, -1, pos=(0, PLAY_AREA_HEIGHT), size=(PLAY_AREA_WIDTH + RANKING_AREA_WIDTH, BOTTOM_AREA_HEIGHT))
        wx.StaticText(self.panel3, -1, "提示: 上下左右键: w、s、a、d    攻击键: u    投掷锤子：i    冲刺：j    加速：k", pos=(10, 5))
        wx.StaticText(self.panel3, -1, "提示: 黄色方块为固定障碍物; 紫色为动态障碍物", pos=(10, 25))


        wx.StaticText(self.panel2, -1, "最大偏差", pos=wx.Point(210, 125))
        self.lbl_max_dis = wx.StaticText(self.panel2, -1, ": %d" % 0, pos=wx.Point(290, 125))

        wx.StaticText(self.panel2, -1, "同步间隔", pos=wx.Point(210, 155))
        self.lbl_sync_interval = wx.StaticText(self.panel2, -1, ": %d" % 0, pos=wx.Point(290, 155))
        
        wx.StaticText(self.panel2, -1, "同步偏差", pos=wx.Point(210, 185))
        self.lblBallDel = wx.StaticText(self.panel2, -1, ": 0", pos=wx.Point(290, 185))
        
        
        self.panel = wx.Panel(self, -1, pos=(0, 0), size=(PLAY_AREA_WIDTH, PLAY_AREA_HEIGHT))
        self.Bind(wx.EVT_CLOSE, self.OnClose)
        self.panel.Bind(wx.EVT_KEY_DOWN, self.OnKeyDown)
        self.panel.Bind(wx.EVT_KEY_UP, self.OnKeyUp)
        self.panel.Bind(wx.EVT_MOUSE_EVENTS, self.OnMouseEvent)
        
        self.timer = wx.Timer(self)
        self.Bind(wx.EVT_TIMER, self.OnTimer, self.timer)
        self.timer.Start(ON_TIMER_INTERVAL)

        self.oneMsTimer = wx.Timer(self)
        self.Bind(wx.EVT_TIMER, self.On1MsgTimer, self.oneMsTimer)
        if self.user.usedelay > 0:
            self.oneMsTimer.Start(1)
        
        self.panel.SetFocus()

        self.InitBG()
        
        self.inited = True
        self.Centre()
        
    def InitBG(self):
        mapsize = int( math.ceil(NET_SCALE * self.res.map["size"] * self.camera.scale) )
        print("mapsize = ", mapsize)
        self.bg = wx.Bitmap(mapsize, mapsize)
        m = wx.MemoryDC(self.bg)
        m.Clear()
        step = int( math.ceil( self.camera.get_cell_size()) )
        print("step1 = ", step)
        for j in range(0, mapsize, step):
            for i in range(0, mapsize, step):
                (_, tmp) = divmod(int(i/step) + int(j/step), 2)
                if tmp == 0:
                    # m.SetBrush(wx.TRANSPARENT_BRUSH)
                    b = wx.Brush(wx.Colour(69, 139, 0))
                    m.SetBrush(b)
                    m.SetPen(wx.TRANSPARENT_PEN)
                else:
                    # m.SetBrush(wx.GREY_BRUSH)
                    b = wx.Brush(wx.Colour(0, 139, 0))
                    m.SetBrush(b)
                    m.SetPen(wx.GREY_PEN)
                m.DrawRectangle(i, j, step, step)
        
        m.SetPen(wx.BLACK_PEN)
        m.SetBrush(wx.TRANSPARENT_BRUSH)
        step = int( math.ceil( NET_SCALE * 5 * self.camera.scale))
        print("step2 = ", step)
        for j in range(0, mapsize, step):
            for i in range(0, mapsize, step):
                m.DrawRectangle(i, j, step, step)
                
        m.SelectObject(wx.NullBitmap)
        del m

    def right_top_data(self):
        if self.user.top != None:
            self.EndTime = self.user.top.EndTime
        else:
            self.EndTime = 0
        return copy.copy(self.user.top)

    def draw_right_top(self, mytop):
        self.listTop.Clear()
        if mytop != None:
            index = 0
            selfindex = -1
            for id, player in mytop.players.items():
                self.listTop.Append("curexp: %d, name: %s" % (player.curexp, player.name))
                if id == self.user.id:
                    selfindex = index
                index = index + 1
            if selfindex >= 0:
                self.listTop.SetSelection(selfindex)
    
    def draw(self, detal):
        # self.user.sync_scene()
        #self.user.match_time()
        if self.inited == False:
            return
        if self.user.ball_id == 0:
            return
        ball = self.user.get_ball()
        if ball == None:
            return
            
        self.lblFPS.SetLabel(": %d" % detal)
        if self.user.udp_pack_size != None:
            self.lblUDP.SetLabel(": %d" % self.user.udp_pack_size)
        else:
            self.lblUDP.SetLabel(": 未开启")
            
            
        self.lblPlayerNum.SetLabel(": %d" % len(self.user.players))
        
        self.lblBallId.SetLabel(": %d" % self.user.ball_id)
        self.lblBallX.SetLabel(": %d" % int(ball.x * self.camera.scale))
        self.lblBallY.SetLabel(": %d" % int(ball.y * self.camera.scale))
        # print("sync_delta: ", self.user.sync_delta, " scale: ", self.camera.scale, " typeof(sync_delta): ", type(self.user.sync_delta), "typeof(scale): ", type(self.camera.scale))
        self.lblBallDel.SetLabel(": %f" % (self.user.sync_delta * self.camera.scale))
        self.lbl_max_dis.SetLabel(": %f" % (self.user.max_dis * self.camera.scale))
        self.lbl_sync_interval.SetLabel(": %f" % ((self.user.thisSyncTime - self.user.lastSyncTime) * self.camera.scale))
        
        
        self.user.mutex.acquire()
        rtData = self.right_top_data()
        myballs = {}
        for k, v in self.user.balls.items():
            myballs[k] = v
        myballs_bg = {}
        for k, v in self.user.balls_bg.items():
            myballs_bg[k] = v
        myplayerballs = {}
        for k, v in self.user.playerballs.items():
            myplayerballs[k] = v
        mynofindball = {}
        for k, v in self.user.no_find_ball.items():
            mynofindball[k] = v
        self.user.mutex.release()

        self.draw_right_top(rtData)
        
        mapsize = NET_SCALE * self.res.map["size"]
        x = 0
        y = 0
        if ball != None:
            x = ball.x
            y = ball.y
        viewdata = self.camera.get_view(    \
            wx.Size(PLAY_AREA_WIDTH, PLAY_AREA_HEIGHT),     \
            wx.Point(x,y),     \
            wx.Size(mapsize, mapsize))
        
        bitmap = wx.Bitmap(viewdata["width"], viewdata["height"])
        m = wx.MemoryDC(bitmap)
        m.Clear()
        
        subBitMap = self.bg.GetSubBitmap(wx.Rect(viewdata["xbegin"], viewdata["ybegin"], PLAY_AREA_WIDTH, PLAY_AREA_HEIGHT))
        m.DrawBitmap(subBitMap, 0, 0)
        
        m.SetBrush(wx.Brush(wx.Colour(184,134,11)))
        for block in self.res.map["nodes"]:
            x = block["px"] * NET_SCALE * self.camera.scale
            y = block["py"] * NET_SCALE * self.camera.scale
            tmps = block["radius"] * NET_SCALE * self.camera.scale
            if x > viewdata["xbegin"] - tmps * 2 and \
                y > viewdata["ybegin"] - tmps * 2 and \
                x < viewdata["xend"] + tmps * 2 and \
                y < viewdata["yend"] + tmps * 2:
                
                if tmps < 5:    #为了看的见
                    tmps = 5
                tmpposx = (x - tmps) - viewdata["xbegin"]
                tmpposy = (y - tmps) - viewdata["ybegin"]
                m.DrawRectangle(tmpposx, tmpposy, tmps*2, tmps*2)
                
                
                
        myball = myplayerballs[self.user.ball_id]
        w = 1800 * self.camera.scale
        h = 1800 * self.camera.scale
        tmpposx = myball.x * self.camera.scale - viewdata["xbegin"]
        tmpposy = myball.y * self.camera.scale - viewdata["ybegin"]
        cellwidth = cellheight = int( math.ceil( NET_SCALE * 5 * self.camera.scale))              
        minX = int(math.floor((tmpposx - w/2)/cellwidth) * cellwidth)
        maxX = int(math.floor((tmpposx + w/2)/cellwidth) * cellwidth)
        minY = int(math.floor((tmpposy - h/2)/cellheight) * cellheight)
        maxY = int(math.floor((tmpposy + h/2)/cellheight) * cellheight)
        
        m.SetBrush(wx.Brush(wx.Colour(238,221,130)))
        for _, bg in myballs_bg.items():
            x = bg.x * self.camera.scale
            y = bg.y * self.camera.scale
            tmps = 25 * self.camera.scale
            if x - viewdata["xbegin"] > minX and \
                y - viewdata["ybegin"] > minY and \
                x - viewdata["xbegin"] < maxX + cellwidth and \
                y - viewdata["ybegin"] < maxY + cellheight:
                
                if tmps < 5:    #为了看的见
                    tmps = 5
                tmpposx = (x - tmps) - viewdata["xbegin"]
                tmpposy = (y - tmps) - viewdata["ybegin"]
                m.DrawRectangle(tmpposx, tmpposy, tmps*2, tmps*2)
                    
        for _, ball in myballs.items():
            tmpposx = ball.x * self.camera.scale - viewdata["xbegin"]# * self.camera.get_cell_size()
            tmpposy = ball.y * self.camera.scale - viewdata["ybegin"]# * self.camera.get_cell_size()
            tmps = int(self.res.food[ball.type]["size"] * NET_SCALE * self.camera.scale)
            if tmps < 5:    #为了看的见
                tmps = 5
            
            t = self.res.food[ball.type]["type"]
            if t == proto.wilds_pb2.FoodNormal:
                m.SetBrush(wx.BLUE_BRUSH)
                m.DrawCircle(tmpposx, tmpposy, tmps)
            elif t == proto.wilds_pb2.FoodHammer:
                m.SetPen(wx.BLACK_PEN)
                m.DrawText("%s" % ("锤子"), tmpposx - 18, tmpposy - 10)
            elif t == proto.wilds_pb2.FoodBomb:
                m.SetPen(wx.BLACK_PEN)
                m.DrawText("%s" % ("炸弹"), tmpposx - 18, tmpposy - 10)
            elif t == proto.wilds_pb2.FeedNormal:
                m.SetBrush(wx.Brush(wx.Colour(153,50,204)))
                m.DrawCircle(tmpposx, tmpposy, tmps)
            elif t == proto.wilds_pb2.SkillHammer:
                m.SetBrush(wx.YELLOW_BRUSH)
                m.DrawCircle(tmpposx, tmpposy, tmps)
                m.DrawText("%s" % ("锤子"), tmpposx - 18, tmpposy - 10)
            elif t == proto.wilds_pb2.SkillBomb:
                m.SetBrush(wx.YELLOW_BRUSH)
                m.DrawCircle(tmpposx, tmpposy, tmps)
                m.DrawText("%s" % ("炸弹"), tmpposx - 18, tmpposy - 10)
                
        for _, ball in myplayerballs.items():
            s = self.res.animal[ball.level]["scale"]
                
            if ball.id == self.user.ball_id:
                m.SetBrush(wx.GREEN_BRUSH)
                self.ballx = ball.x * self.camera.scale - viewdata["xbegin"]
                self.bally = ball.y * self.camera.scale - viewdata["ybegin"]
                
                prex = (ball.client_now_x - ball.client_pre_x ) * (ball.client_now_x - ball.client_pre_x )
                prey = (ball.client_now_y - ball.client_pre_y ) * (ball.client_now_y - ball.client_pre_y )
                self.lblMoveDisC.SetLabel(": %f (100ms)" % (math.sqrt(prex + prey)))

                prex = (ball.server_now_x - ball.server_pre_x ) * (ball.server_now_x - ball.server_pre_x )
                prey = (ball.server_now_y - ball.server_pre_y ) * (ball.server_now_y - ball.server_pre_y )        
                self.lblMoveDisS.SetLabel(": %f (100ms)" % (math.sqrt(prex + prey)))
                
            else:
                m.SetBrush(wx.RED_BRUSH)
                
            m.SetPen(wx.BLACK_PEN)

            tmpposx = ball.x * self.camera.scale - viewdata["xbegin"]
            tmpposy = ball.y * self.camera.scale - viewdata["ybegin"]
            tmps = int(s * NET_SCALE * self.camera.scale)
            if tmps < 5:    #为了看的见
                tmps = 5
            m.DrawCircle(tmpposx, tmpposy, tmps)
            
            m.SetPen(wx.RED_PEN)
            m.SetBrush(wx.TRANSPARENT_BRUSH)
            if ball.id == self.user.ball_id:
                m.DrawRectangle(tmpposx - w/2, tmpposy - h/2, w, h)
                m.DrawRectangle(minX, minY, maxX - minX + cellwidth, maxY - minY + cellheight)
                
            
            m.SetPen(wx.BLACK_PEN)
            
            tempv = None
            if ball.face != 0:
                if ball.face in myplayerballs:
                    tempv = vector2.Vector(myplayerballs[ball.face].x - ball.x, myplayerballs[ball.face].y - ball.y)
            if tempv == None:
                angleX = math.cos(math.pi * ball.angle / 180.0)
                angleY = -math.sin(math.pi * ball.angle / 180.0)
                tempv = vector2.Vector(angleX, angleY)
            tempv.NormalizeSelf()
            tempv.ScaleBy(tmps)
            tempv.IncreaseBy(vector2.Vector(tmpposx, tmpposy))
            m.DrawLine(tmpposx, tmpposy, tempv.x, tempv.y)
            
            
            if ball.state != 0:
                m.DrawText("%s" % state.NAME[ball.state], tmpposx - 20, tmpposy - 35)
            m.DrawText("hp:%d" % ball.hp, tmpposx - 20, tmpposy - 20)
            if ball.id == self.user.ball_id and ball.curmp!=None:
                m.DrawText("curmp:%d" % ball.curmp, tmpposx - 20, tmpposy - 5)
            m.DrawText("lvl:%d" % ball.level, tmpposx - 15, tmpposy + 10)

            # 经验条
            if ball.id == self.user.ball_id:
                nextLevelExp = self.user.res.get_next_level_exp(ball.level)
                curLevelExp = 0
                if ball.level > 1:
                    curLevelExp = self.user.res.get_next_level_exp(ball.level - 1)

                if ball.curexp - curLevelExp > 0:
                    angle = math.pi * 2 * (ball.curexp - curLevelExp) / (nextLevelExp - curLevelExp)
                    # print("level:", ball.level, " nle:", nextLevelExp, " cle:", curLevelExp, " curexp:", ball.curexp, " angle:", angle / math.pi * 180)
                    tox = tmpposx + math.cos(angle) * tmps
                    toy = tmpposy - math.sin(angle) * tmps

                    pen = wx.Pen(wx.Colour(255, 0, 0), 2)
                    m.SetPen(pen)
                    m.DrawArc(tmpposx + tmps, tmpposy, tox, toy, tmpposx, tmpposy)

            
        m.SetBrush(wx.WHITE_BRUSH)
        for _, ball in mynofindball.items():
            tmpposx = ball.x * self.camera.scale - viewdata["xbegin"]
            tmpposy = ball.y * self.camera.scale - viewdata["ybegin"]
            tmps = 30 * self.camera.scale
            m.DrawCircle(tmpposx, tmpposy, tmps)
            m.DrawText("%d" % ball.id, tmpposx - 10, tmpposy - 15)
            
        if self.EndTime!=None:
            m.DrawText("剩余时间：%d" % (self.EndTime), PLAY_AREA_WIDTH / 2 - 50, 50)
        m.SelectObject(wx.NullBitmap)
        del m
        
        dc = wx.ClientDC(self.panel)
        #dc.DrawBitmap(bitmap, -viewdata["draw_xbegin"], -viewdata["draw_ybegin"])
        dc.DrawBitmap(bitmap, 0, 0)
        del dc

    def SetNextPos(self):
        self.user.mutex.acquire()
        mapsize = NET_SCALE * self.res.map["size"]
        radius = 0.25 * NET_SCALE

        # if self.pre_cal_speed_time == 0:
        #     self.pre_cal_speed_time = int(round(time.time() * 1000))
        now = int(round(time.time() * 1000))
        # detal = now - self.pre_cal_speed_time
        # self.pre_cal_speed_time = now

        for _, v in self.user.playerballs.items():
            
            #new
            # prex = v.x
            # prey = v.y
            # passTime = now - v.last_update_time
            # if now < v.this_cli_time:
            #     v.x += v.vx * passTime / 50
            #     v.y += v.vy * passTime / 50
            #     moveDis = math.sqrt((v.x - prex) * (v.x - prex) + (v.y - prey) * (v.y - prey))
            #     # print("move dis ", moveDis, " passtime ", passTime, " vx ", v.vx)
            # elif v.last_update_time >= v.this_cli_time:
            #     v.x += v.this_vx * passTime / 50
            #     v.y += v.this_vy * passTime / 50
            #     moveDis = math.sqrt((v.x - prex) * (v.x - prex) + (v.y - prey) * (v.y - prey))
            #     # print("move dis ", moveDis, " passtime ", passTime, " this_vx ", v.this_vx)
            # else:
            #     t1 = v.this_cli_time - v.last_update_time
            #     v.x += v.vx * t1 / 50
            #     v.y += v.vy * t1 / 50

            #     t2 = now - v.this_cli_time
            #     v.x += v.this_vx * t2 / 50
            #     v.y += v.this_vy * t2 / 50
            #     moveDis = math.sqrt((v.x - prex) * (v.x - prex) + (v.y - prey) * (v.y - prey))
            #     # print("move dis ", moveDis, " passtime ", passTime, " (vx) ", v.vx, " ", v.this_vx)
            # #new

            # #old
            # # v.x += v.vx*detal/50
            # # v.y += v.vy*detal/50
            # #old

            # if v.x > mapsize - radius:
            #     v.x = mapsize - radius
            # elif v.x < radius:
            #     v.x = radius
            # if v.y > mapsize - radius:
            #     v.y = mapsize - radius
            # elif v.y < radius:
            #     v.y = radius

            # if abs(v.x - v.this_x) < 1:
            #     v.x = v.this_x
            # if abs(v.y - v.this_y) < 1:
            #     v.y = v.this_y

            # v.last_update_time = now

            self.user.gameSync.UpdateMove(v, mapsize, now, self)

            # if now >= v.this_cli_time + 100:
            #     v.vx = v.this_vx
            #     v.vy = v.this_vy

            if v.id == self.user.ball_id:
                if self.show_line:
                    m = wx.MemoryDC(self.bg)
                    m.SetBrush(wx.BLACK_BRUSH)
                    m.DrawCircle(v.x, v.y, 1)

                    m.SetBrush(wx.RED_BRUSH)
                    m.DrawCircle(v.this_x, v.this_y, 1)

                    interval = self.user.thisSyncTime - self.user.lastSyncTime
                    if interval >= 150:
                        m.DrawText("%dms" % interval, v.x, v.y)
                        self.user.lastSyncTime = self.user.thisSyncTime

                    m.SelectObject(wx.NullBitmap)
                    del m
                    
        for _, v in self.user.balls.items():
            if v.isplayer != None:
                self.user.gameSync.UpdateMove(v, mapsize, now, self)
                    
        self.user.mutex.release()

        
    def OnViewNormal(self, event):
        self.camera.scale = 1
        self.InitBG()
        
    def OnViewAll(self, event):
        tmpsize = min(PLAY_AREA_WIDTH, PLAY_AREA_HEIGHT)
        self.camera.scale = tmpsize / (self.res.map["size"] * NET_SCALE)
        self.InitBG()
        
    def OnRelife(self, event):
        self.user.relife()
        # self.panel.SetFocus()
        
    def OnExitRoom(self, event):
        self.user.exit_room()

    def OnKeyDown(self, event):
        if self.user.mouseMove > 0:
            return
        olddir = self.dir
        oldpower = self.power
        kc = event.GetKeyCode()
        if kc in [LEFT, RIGHT, UP, DOWN]:
            left = wx.GetKeyState(LEFT) and -1 or 0
            right = wx.GetKeyState(RIGHT) and 1 or 0
            up = wx.GetKeyState(UP) and -1 or 0
            down = wx.GetKeyState(DOWN) and 1 or 0
            self.calc_dir_power(left, right, up, down)

    def OnKeyUp(self, event):
        kc = event.GetKeyCode()
        if kc == ATTACK:
            self.normal_attack()
        elif kc == THROW_HAMMER:
            self.throw_hammer()
        elif kc == THROW_BOMB:
            self.throw_bomb()
        elif kc == RUN:
            self.user.run()

        # ball = self.user.get_ball()
        # if self.power > 0 and ball and ball.vx == 0 and ball.vy == 0:
        # ball.vx = 8.333 * math.cos(math.pi * self.dir / 180)
        # ball.vy = -8.333 * math.sin(math.pi * self.dir / 180)
            # print("vx: ", ball.vx, " vy: ", ball.vy)
        # if self.power == 0:
        #     ball.vx = 0
        #     ball.vy = 0
        if self.user.mouseMove == 0:
            olddir = self.dir
            oldpower = self.power
            if kc in [LEFT, RIGHT, UP, DOWN]:
                left = wx.GetKeyState(LEFT) and -1 or 0
                right = wx.GetKeyState(RIGHT) and 1 or 0
                up = wx.GetKeyState(UP) and -1 or 0
                down = wx.GetKeyState(DOWN) and 1 or 0
                self.calc_dir_power(left, right, up, down)
            #if self.dir != olddir or oldpower != self.power:
            #    self.user.move(self.dir, self.power)

    def normal_attack(self):
        self.user.cast_skill(100)
    def throw_hammer(self):
        self.user.cast_skill(103)
    def throw_bomb(self):
        self.user.cast_skill(104)

    def calc_dir_power(self, left, right, up, down):
        if left + right == 0 and up + down == 0:
            self.dir = 0
            self.power = 0
        else:
            if self.power == 0:
                self.power = 99
            else:
                self.power = 100
                
            if right != 0:
                if up !=0:
                    self.dir = 45
                elif down !=0:
                    self.dir = 315
                else:
                    self.dir = 0
            elif up != 0:
                if left != 0:
                    self.dir = 135
                else:
                    self.dir = 90
            elif left != 0:
                if down != 0:
                    self.dir = 225
                else:
                    self.dir = 180
            elif down != 0:
                self.dir = 270
            else:
                print("error!!!!")
                exit(0)
                
    def OnMouseEvent(self, event):
        if event.LeftDown():
            #print("leftdown")
            ball = self.user.get_ball()
            if ball.curmp > 2:
                self.user.cast_skill(107)
        elif event.RightDown():
            ball = self.user.get_ball()
            if ball.curmp > 2:
                self.user.run()
        elif event.Moving() and self.user.mouseMove > 0:
            self.power = 100
            pos = event.GetPosition()
            self.dir = int(self.cal_dir_by_point(self.ballx, self.bally, pos.x, pos.y))
            if math.sqrt((pos.x - self.ballx)*(pos.x - self.ballx) + (pos.y - self.bally)*(pos.y - self.bally)) < 25:
                self.power = 0

        event.Skip()

    def cal_dir_by_point(self, x1, y1, x2, y2):
        dis = math.sqrt((x2 - x1) * (x2 - x1) + (y2 - y1) * (y2 - y1))
        if 0 == dis:
            return 0
        r = math.acos((x2 - x1) / dis)
        if y2 - y1 < 0:
            r = math.pi * 2 - r
        d = -r * 180 / math.pi
        if d < 0:
            d += 360
        return d
        

    def OnTimer(self, evt):
        self.frame = self.frame + 1
        
        (_, tmp) = divmod(self.frame, 1200) #30s,服务器60s没收到消息，踢人
        if tmp == 0:
            self.user.heart_beat()
            
            
        select_target.near_target(self.user, PLAY_AREA_WIDTH, PLAY_AREA_HEIGHT)
        
        self.user.move(self.dir, self.power)
        
        if self.pre_on_timer == 0:
            self.pre_on_timer = int(round(time.time() * 1000))
        now = int(round(time.time() * 1000))
        detal = now - self.pre_on_timer
        self.pre_on_timer = now
        self.SetNextPos()
        self.draw(detal)

    def On1MsgTimer(self, evt):
        self.user.sync_udpscene()
        if self.user.mouseMove > 0:
            now = int(time.time() * 1000)
            if now - self.sendMoveTime >= 100:
                self.user.move(self.dir, self.power)
                # print(self.dir)
                self.sendMoveTime = now

    def OnClose(self, evt):
        print("OnClose")
        if self.timer != None:
            self.timer.Stop()
        if self.oneMsTimer != None:
            self.oneMsTimer.Stop()
        self.user.close()
            
        evt.Skip()
        
    def NetScale(self):
        return NET_SCALE


g_roomwin = None      
def new(user, args, cfg):
    global g_roomwin
    g_roomwin = RoomWindow(user, args, cfg)
    g_roomwin.Show()
    
def delete():
    global g_roomwin
    if g_roomwin != None:
        g_roomwin.Close()
        g_roomwin = None
    