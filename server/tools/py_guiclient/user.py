# -*- coding:utf8 -*-
import urllib.request
from socket import *
from frame_mgr import *
import threading
import mainwindow
import proto.message
import proto.myprint
import proto.player_info
import proto.top_info
import proto.player_ball_info
import proto.ball_info
import proto.player_animal_info
import proto.player_pb2
import proto.wilds_pb2
import math
import time
import random
import matchtime
import copy
import gamesync.follow
import gamesync.follow_ori
import gamesync.follow_physics
import gamesync.mope
import wx
import socket_mgr
import res
import gateway
import log

print=log.log

BUFSIZ = 128 * 1024


class User():
    def __init__(self, sockindex, args, cfg):
        self.sockindex = sockindex
        self.args = args
        self.cfg = cfg
        self.addr = self.cfg["addr"]
        self.port = self.cfg["port"]
        self.ver = self.cfg["ver"]
        
        # self.scene_id = 0
        self.scene_id = self.cfg["room"]["scene"]
        self.room_type = 0
        self.usedelay = self.cfg["room"]["usedelay"]
        self.delaymin = self.cfg["room"]["delaymin"]
        self.delaymax = self.cfg["room"]["delaymax"]
        self.kcp_enable = self.cfg["room"]["kcp_enable"]
        self.udp_enable = self.cfg["room"]["udp_enable"]
        self.mouseMove = self.cfg["room"]["mouse_move"]
        self.cookie = ""
        self.account_name = ""
        self.id = 0
        self.ball_id = 0
        
        self.face = 0
        self.old_face = 0
        self.old_power = 0
        self.old_way = 0
        
        self.framemgr = FrameMgr()
        
        self.room_server_addr = ""
        self.room_server_key = ""
        self.client = None
        self.client_thread = None
        self.terminate = False
        self.player_name = ""
        self.players = {}           # need lock
        self.playerballs = {}       # need lock
        self.balls = {}             # need lock
        self.balls_bg = {}          # need lock
        self.top = None             # need lock
        self.no_find_ball = {}      # need lock
        self.mutex = threading.Lock()
        self.room_end_time = 0
        self.res = None

        self.udpclient = None
        self.udp_pack_size = None
        
        
        self.left_data_tcp = b""
        self.left_data_udp = b""
        self.pystress = False
        
        
        self.sync_delta = 0.0
        self.drag_back = 0
        self.max_dis = 0
        self.lastSyncTime = 0
        self.thisSyncTime = 0

        self.scenemsg_list = []
        self.udpsenemsg_list = []

        self.team_rank_list = None
        self.team_info = None

        self.cur_match_time = None
        self.next_match_time = None

        self.is_live = True

        if self.cfg["room"]["synctype"] == 1:
            self.gameSync = gamesync.follow_ori.FollowSync()
        elif self.cfg["room"]["synctype"] == 2:
            self.gameSync = gamesync.follow.FollowSync()
        elif self.cfg["room"]["synctype"] == 3:
            self.gameSync = gamesync.follow_physics.FollowSync()
        elif self.cfg["room"]["synctype"] == 4:
            self.gameSync = gamesync.mope.FollowSync()
        
        self.cmds = {}
        self.init_cmds()
        
        
    def init_cmds(self):
        self.cmds[proto.message.Wilds.Login.value] = self.on_wilds_login
        self.cmds[proto.message.Wilds.Top.value] = self.on_wilds_top
        self.cmds[proto.message.Wilds.AddPlayer.value] = self.on_wilds_add_player
        self.cmds[proto.message.Wilds.RemovePlayer.value] = self.on_wilds_remove_player
        self.cmds[proto.message.Wilds.Death.value] = self.on_wilds_death
        self.cmds[proto.message.Wilds.EndRoom.value] = self.on_end_room
        self.cmds[proto.message.Wilds.RefreshPlayer.value] = self.on_wilds_refresh_player
        self.cmds[proto.message.Wilds.SceneChat.value] = self.on_wilds_donothing
        self.cmds[proto.message.Wilds.ActCloseSocket.value] = self.on_exit_room
        self.cmds[proto.message.Wilds.VoiceInfo.value] = self.on_wilds_donothing
        self.cmds[proto.message.Wilds.CastSkill.value] = self.on_cast_skill
        self.cmds[proto.message.Wilds.HeartBeat.value] = self.on_wilds_donothing
        self.cmds[proto.message.Wilds.AsyncPlayerAnimal.value] = self.on_sync_player_animal
        self.cmds[proto.message.Wilds.UpdateTeamInfo.value] = self.on_update_team_info
        self.cmds[proto.message.Wilds.TeamRankList.value] = self.on_team_ranking_list
        self.cmds[proto.message.Wilds.ReLife.value] = self.on_relife
        self.cmds[proto.message.Wilds.SceneTCP.value] = self.on_wilds_scene_tcp
        if self.usedelay > 0:
            self.cmds[proto.message.Wilds.SceneUDP.value] = self.on_delay_udp_scene
        else:
            self.cmds[proto.message.Wilds.SceneUDP.value] = self.on_wilds_scene_udp
        
    def login(self):
        account = self.cfg["login"]["account"]
        password = self.cfg["login"]["password"]
        dev = self.cfg["login"]["dev"]
        device = self.cfg["login"]["device"]
        r = self.cfg["login"]["r"]
        m = self.cfg["login"]["m"]
        return self.login_detail(account, password, dev, device, r, m)
        
    def login_detail(self, account, password, dev, device, r, m):
        try:
            #like: http://127.0.0.1:8080/msg?c=1&a=BOS30000022&p=&ver=1.2.0&dev=4730d114e32c13359a112ce6cc17eebbd2073944&device=android&r=1&m=0&sign=7e343d61e73ecd3c600521ae9588c460
        
            c = proto.message.Player.Login.value
            sign = proto.message.get_sign(c)
            url = "http://%s:%d/msg?c=%d&a=%s&p=%s&ver=%s&dev=%s&device=%s&r=%d&m=%d&sign=%s" % (self.addr, self.port, c, account, password, self.ver, dev, device, r, m, sign)
            print("url=", url)
        
            request = urllib.request.urlopen(url)
            result = request.read()
        
            self.cookie = request.getheader("Set-Cookie")
            if self.cookie != "":
                self.cookie = self.cookie[self.cookie.find("=")+1: self.cookie.find(";")]
        
            cmd, msg = proto.message.unpack(result)
            pmsg = proto.player_pb2.RetLoginMsg()
            pmsg.ParseFromString(msg)
            proto.myprint.print_login_result(pmsg)
        
            self.account_name = pmsg.Account
            self.id = pmsg.Id
        
            if pmsg.Id == 0:
                print("connect login server fail.")
                return False
        
            # TODO: 暂时屏蔽登录gateway
            '''
            if self.pystress == False:
                if gateway.gate.login_gateway(self) == False:
                    print("connect gateway server fail.")
                    return False
            '''
            
        except Exception as e:
            print(e)
            return False
            
        return True
        
        
    def req_room(self):
    
        if self.cookie == "":
            return False
    
        #like: http://127.0.0.1:8080/msg?c=3&ver=1.2.0&ticketnum=0&scene=1002&sign=tolower(md5(SIGNKEYc))
        
        c = proto.message.Player.ReqIntoFRoom.value
        sign = proto.message.get_sign(c)
        self.scene_id = scene = self.cfg["room"]["scene"]
        url = "http://%s:%d/msg?c=%d&ver=%s&ticketnum=%d&scene=%d&sign=%s" % (self.addr, self.port, c, self.ver, 0, scene, sign)
        print("url=", url)
        
        request = urllib.request.Request(url, headers={"Cookie":"session="+self.cookie})
        result = urllib.request.urlopen(request).read()
        cmd, msg = proto.message.unpack(result)
        pmsg = proto.player_pb2.RetIntoFRoom()
        pmsg.ParseFromString(msg)
        proto.myprint.print_room_result(pmsg)
        
        if pmsg.Err == 0:
            self.room_server_addr = pmsg.Addr
            self.room_server_key = pmsg.Key
        
        return pmsg.Err == 0
        
        
    def connect_room_server(self):
        if self.room_server_addr == "" or self.room_server_key == "":
            print('self.room_server_addr == "" or or self.room_server_key == ""')
            return False
    
        host = self.room_server_addr[:self.room_server_addr.find(":")]
        port = int(self.room_server_addr[self.room_server_addr.find(":") + 1:])
        addr = (host, port)
        
        if self.client != None:
            socket_mgr.get_socket_set("tcp").remove_sock(self.sockindex)
            self.client.close()
            self.client = None
        self.client = socket_mgr.get_socket_set("tcp").add_sock(self.sockindex, SOCK_STREAM)
        
        try:
            self.client.connect(addr)
        except Exception as e:
            print(e)
            return False
        
        if self.udp_enable!=0:
            if self.udpclient != None:
                socket_mgr.get_socket_set("udp").remove_sock(self.sockindex)
                self.udpclient.close()
                self.udpclient = None
            self.udpclient = socket_mgr.get_socket_set("udp").add_sock(self.sockindex, SOCK_DGRAM)
            self.udpclient.connect(addr)

        return True
        
    
    def start_game(self):
        if self.client == None:
            print("restart ...")
            self.terminate = False
            self.req_room()
            if self.connect_room_server() == False:
                return False
            time.sleep(0.02)
        
        print("call start_game #1")
       
        pmsg = proto.wilds_pb2.MsgLogin()
        pmsg.name = "%s_%d" % (self.account_name, self.id)
        pmsg.key = self.room_server_key
        data = pmsg.SerializeToString()
        cmd = proto.message.Wilds.Login.value
        msg = proto.message.pack(cmd, data)
        if self.client != None:
            self.client.send(msg)
        ##todo:对时
        #self.match_time()
        
        print("call start_game #2")
        
        #加载资源
        if self.pystress==False:
            self.res = res.new(self.args, self.cfg, self.scene_id)
        print("call start_game #3")
        
        return True
        
    def move(self, angle, power):
        if self.face == self.old_face and self.old_power == power and self.old_way == angle:
            return
        self.old_face = self.face
        self.old_power = power
        self.old_way = angle
        
        pmsg = proto.wilds_pb2.MsgMove()
        pmsg.angle = angle
        pmsg.power = power
        pmsg.face = self.face
        data = pmsg.SerializeToString()
        cmd = proto.message.Wilds.Move.value
        msg = proto.message.pack(cmd, data)
        if self.client != None:
            self.client.send(msg)

    def run(self):
        pmsg = proto.wilds_pb2.MsgRun()
        data = pmsg.SerializeToString()
        cmd = proto.message.Wilds.Run.value
        msg = proto.message.pack(cmd, data)
        if self.client != None:
            self.client.send(msg)

    def relife(self):
        pmsg = proto.wilds_pb2.MsgRelife()
        data = pmsg.SerializeToString()
        cmd = proto.message.Wilds.ReLife.value
        msg = proto.message.pack(cmd, data)
        if self.client != None:
            self.client.send(msg)

    def bindTCPSession(self):
        pmsg = proto.wilds_pb2.MsgBindTCPSession()
        pmsg.id = self.id
        pmsg.key = self.room_server_key
        data = pmsg.SerializeToString()
        cmd = proto.message.Wilds.BindTCPSession.value
        msg = proto.message.pack(cmd, data)
        if self.client == None:
            return
        if self.udp_enable > 0 and self.udpclient != None:
            self.udpclient.sendall(msg)


    def cast_skill(self, skillId):
        pmsg = proto.wilds_pb2.MsgCastSkill()
        pmsg.skillid = skillId
        data = pmsg.SerializeToString()
        cmd = proto.message.Wilds.CastSkill.value
        msg = proto.message.pack(cmd, data)
        if self.client != None:
            self.client.send(msg)
            
    def heart_beat(self):
        cmd = proto.message.Wilds.HeartBeat.value
        msg = proto.message.pack(cmd, b"")
        if self.client != None:
            self.client.send(msg)
            
    def exit_room(self):
        cmd = proto.message.Wilds.ActCloseSocket.value
        msg = proto.message.pack(cmd, b"")
        if self.client != None:
            self.client.send(msg)

    def on_cast_skill(self, data):
        pass
    
    def on_sync_player_animal(self, data):
        pmsg = proto.wilds_pb2.MsgAsyncPlayerAnimal()
        pmsg.ParseFromString(data)
        
        self.mutex.acquire()
        if pmsg.id in self.players:
            player = self.players[pmsg.id]
            if player.ballId in self.playerballs:
                ball = self.playerballs[player.ballId]
                ball.level = pmsg.animalid
        self.mutex.release()
        
        
    def get_ball(self):
        self.mutex.acquire()
        ball = None
        if self.ball_id in self.playerballs:
            ball = self.playerballs[self.ball_id]
        self.mutex.release()
        return ball
        
    def on_recv(self, data, udp_channel=0):
        msgs = proto.message.on_recv(self, data, udp_channel)
        if len(msgs) == 0:
            return

        for (cmd, data) in msgs:
            self.on_recv_one(cmd, data, udp_channel)
                
    def on_recv_one(self, cmd, data, udp_channel=0):
        if cmd in self.cmds:
            self.cmds[cmd](data)
        else:
            print("cmd = ", cmd)
        if udp_channel != 0:
            self.udp_pack_size = len(data)

        
    def on_wilds_login(self, data):
        print("on_wilds_login")
        pmsg = proto.wilds_pb2.MsgLoginResult()
        pmsg.ParseFromString(data)
        proto.myprint.print_wilds_login_result(pmsg)
        
        if pmsg.ok:
            self.player_name = pmsg.name
            self.ball_id = pmsg.ballId
            self.framemgr.set_frame(pmsg.frame)
            print("lefttime:", pmsg.leftTime)
            self.room_end_time = pmsg.leftTime + round(time.time())
            
            self.mutex.acquire()
            
            for p in pmsg.others:
                player = proto.player_info.PlayerInfo(p)
                self.players[player.id] = player
                
            for pb in pmsg.playerballs:
                playerball = proto.player_ball_info.PlayerBallInfo(pb)
                self.playerballs[playerball.id] = playerball
                #print( "add playerball ====================="   )
                #print( playerball.id )
                
            for b in pmsg.balls:
                ball = proto.ball_info.BallInfo(b)
                self.balls[ball.id] = ball
                
                if self.pystress==False:
                    t = self.res.food[ball.type]["type"]
                    if t == proto.wilds_pb2.FoodHammer or t == proto.wilds_pb2.FoodBomb:
                        self.balls_bg[ball.id] = ball
            
            self.mutex.release()
            
            
            '''
            # 选动物
            if pmsg.IsFirstCross == False:
                pmsg = proto.wilds_pb2.MsgSelectAnimal()
                data = pmsg.SerializeToString()
                cmd = proto.message.Wilds.SelectAnimal.value
                msg = proto.message.pack(cmd, data)
                if self.client != None:
                    self.client.send(msg)
            '''
            

            # 绑定TCP连接
            if self.kcp_enable > 0 or self.udp_enable > 0:
                self.bindTCPSession()
                
            if self.pystress==False:
                mainwindow.open_room_window(self)
            
        else:
            print("start game fail")
            
    def on_wilds_top(self, data):
        pmsg = proto.wilds_pb2.MsgTop()
        pmsg.ParseFromString(data)
        self.mutex.acquire()
        self.top = proto.top_info.TopInfo(pmsg)
        self.mutex.release()
        
    def on_wilds_add_player(self, data):
        pmsg = proto.wilds_pb2.MsgAddPlayer()
        pmsg.ParseFromString(data)
        player = proto.player_info.PlayerInfo(pmsg.player)
        
        self.mutex.acquire()
        self.players[player.id] = player
        self.mutex.release()
        
    def on_wilds_remove_player(self, data):
        pmsg = proto.wilds_pb2.MsgRemovePlayer()
        pmsg.ParseFromString(data)
        
        self.mutex.acquire()
        if pmsg.id in self.players:
            player = self.players[pmsg.id]
            if player.ballId in self.playerballs:
                self.playerballs.pop(player.ballId)
            self.players.pop(pmsg.id)
        self.mutex.release()


    def on_delay_scene(self, data):
        self.mutex.acquire()
        self.scenemsg_list.append((int(round(time.time() * 1000)), data))    
        self.mutex.release()

    def on_delay_udp_scene(self, data):
        self.mutex.acquire()
        self.udpsenemsg_list.append((int(round(time.time() * 1000)), data))
        self.mutex.release()

    def sync_udpscene(self):
        if len(self.udpsenemsg_list) == 0:
            return
        (t, data) = self.udpsenemsg_list[0]
        Now = int(round(time.time() * 1000))
        delay = random.uniform(self.delaymin, self.delaymax)
        if Now - t >= delay:
            self.on_wilds_scene_udp(data)
            self.udpsenemsg_list.pop(0)

    def on_wilds_scene_tcp(self, data):
        pmsg = proto.wilds_pb2.MsgSceneTCP()
        pmsg.ParseFromString(data)
        
        self.mutex.acquire()
        for bid in pmsg.removes:
            if bid in self.balls:
                self.balls.pop(bid)
        for b in pmsg.adds:
            ball = proto.ball_info.BallInfo(b)
            
            if ball.id in self.balls:
                if self.pystress==False:
                    print("add ball, but ball already exist!!!!! ball_id =", ball.id, "ball_type =", self.res.food[ball.type]["type"])
                
            self.balls[ball.id] = ball
            if self.pystress==False:
                t = self.res.food[ball.type]["type"]
                if t == proto.wilds_pb2.FoodHammer or t == proto.wilds_pb2.FoodBomb:
                    self.balls_bg[ball.id] = ball
                
            if ball.id in self.no_find_ball:
                self.no_find_ball.pop(ball.id)
                
        for pbid in pmsg.removePlayers:
            #print( "remove playerball ===================== #2"   )
            if pbid in self.playerballs:
                self.playerballs.pop(pbid)
                #print( "remove playerball ===================== #1"   )
                #print( pbid )
        for pb in pmsg.addPlayers:
            playerball = proto.player_ball_info.PlayerBallInfo(pb)
            
            if playerball.id in self.playerballs:
                print("add player ball, but player ball already exist!!!!! ball_id =", playerball.id)
            
            self.playerballs[playerball.id] = playerball
            #print( "add playerball ====================="   )
            #print( playerball.id )
            
            if playerball.id in self.no_find_ball:
                self.no_find_ball.pop(playerball.id)
            
            playerball.server_pre_x = playerball.x
            playerball.server_pre_y = playerball.y
            playerball.client_pre_x = playerball.x
            playerball.client_pre_y = playerball.y
            playerball.server_now_x = playerball.x
            playerball.server_now_y = playerball.y
            playerball.client_now_x = playerball.x
            playerball.client_now_y = playerball.y
                
        for e in pmsg.eats:
            if e.target in self.balls:
                #print (self.balls[e.target], e.target)
                self.balls.pop(e.target)

        for hit in pmsg.hits:
            if hit.target in self.playerballs:
                ball = self.playerballs[hit.target]
                ball.hp = hit.curHp
        self.mutex.release()

    def on_wilds_scene_udp(self, data):
        pmsg = proto.wilds_pb2.MsgSceneUDP()
        pmsg.ParseFromString(data)

        frame = self.framemgr.get_pre_frame()
        if pmsg.frame < frame:
            return
        
        self.framemgr.set_frame(pmsg.frame)
        
        self.mutex.acquire()
        # print("bigdevil_ballid ", self.bigdevil_ballid)
        for m in pmsg.moves:
            ball = None
            if m.id in self.playerballs:
                ball = self.playerballs[m.id]
                ball.isplayer = True
            elif m.id in self.balls:
                ball = self.balls[m.id]
                ball.isplayer = False
            else:
                pass
                #print("no find ball. id = ", m.id)
                
                #self.no_find_ball[m.id] = m
                
                
            if ball != None:
                ball.state = m.state
                ball.angle = m.angle
                ball.face = m.face

                ball.server_pre_x = ball.server_now_x
                ball.server_pre_y = ball.server_now_y
                ball.client_pre_x = ball.client_now_x
                ball.client_pre_y = ball.client_now_y
                ball.server_now_x = m.x
                ball.server_now_y = m.y
                ball.client_now_x = ball.x
                ball.client_now_y = ball.y
                
                if (ball.id == self.ball_id and self.is_live == True) or (ball.id != self.ball_id):
                    self.gameSync.SyncMove(ball, m, None, self.cur_match_time)
                    
        self.mutex.release()
        
    def on_wilds_death(self, data):
        pmsg = proto.wilds_pb2.MsgDeath()
        pmsg.ParseFromString(data)
        
        if pmsg.Id == self.id:
            self.is_live = False
            if self.pystress==False:
                wx.CallAfter(wx.MessageBox, "你已经挂了，请点击复活按钮，继续战斗！", "提示", wx.OK | wx.ICON_INFORMATION)
        else:
            pass
            
    def on_end_room(self, data):
        if self.pystress==False:
            wx.CallAfter(self.on_end_room_detail)
        
        
    def on_end_room_detail(self):
        if self.pystress==False:
            wx.MessageBox("本局游戏结束！", "提示", wx.OK | wx.ICON_INFORMATION)
        self.on_exit_room(None)
    
    def allstop(self):
        for k,v in self.playerballs.items():
            v.vx = 0
            v.vy = 0
        
    def on_wilds_refresh_player(self, data):
        pmsg = proto.wilds_pb2.MsgRefreshPlayer()
        pmsg.ParseFromString(data)
        player = proto.player_info.PlayerInfo(pmsg.player)
        
        self.mutex.acquire()
        self.players[player.id] = player

        if player.ballId in self.playerballs:
            ball = self.playerballs[player.ballId]
            ball.hp = player.curhp
            ball.curmp = player.curmp
            ball.curexp = player.curexp
            # print("hp:", player.curhp, " curmp:", player.curmp)

        self.mutex.release()
        
    def on_exit_room(self, data):
        self.on_exit_room_clear_data()
        if self.pystress==False:
            mainwindow.close_room_window(self)
        
    def on_exit_room_clear_data(self):
        self.mutex.acquire()
        self.players = {}
        self.playerballs = {}
        self.balls = {}
        self.balls_bg = {}
        self.top = None
        self.no_find_ball = {}
        self.mutex.release()
        
    def on_wilds_donothing(self, data):
        pass

    def on_hit(self, data):
        pmsg = proto.wilds_pb2.HitMsg()
        pmsg.ParseFromString(data)
        proto.myprint.print_hit_msg(pmsg)

    def on_update_team_info(self, data):
        pmsg = proto.wilds_pb2.UpdateTeamInfoMsg()
        pmsg.ParseFromString(data)
        
        self.team_info = []
        for mem in pmsg.members:
            if mem.playerid in self.players:
                player = self.players[mem.playerid]
                self.team_info.append((mem.playerid, mem.x * 100, mem.y * 100, player.ballId))
            else:
                self.team_info.append((mem.playerid, mem.x * 100, mem.y * 100, 0))
        self.team_info.sort(key=lambda a:a[0])

        # for mem in pmsg.members:
        #     print("mem id: ", mem.playerid, " (", mem.x, ",", mem.y, ")")
        # for mem in pmsg.topPlayers:
        #     print("top id: ", mem.playerid, " (", mem.x, ",", mem.y, ")")

    def on_team_ranking_list(self, data):
        pmsg = proto.wilds_pb2.RetTeamRankList()
        pmsg.ParseFromString(data)
        self.team_rank_list = pmsg
        # for team in pmsg.Teams:
        #     print("team name: ", team.Tname, " num: ", team.Num, " CorpName: ", team.CorpName, " score: ", team.Score, " last rank: ", team.LastRank)
        # print("Watch num: ", pmsg.WatchNum, " EndTime: ", pmsg.EndTime, " KillNum: ", pmsg.killNum)

    def on_relife(self, data):
        pmsg = proto.wilds_pb2.MsgS2CRelife()
        pmsg.ParseFromString(data)
        
        if pmsg.SnapInfo.Id == self.id:
            self.is_live = True
            self.mutex.acquire()
            self.framemgr.set_frame(pmsg.frame)
            self.playerballs = {}
            self.balls = {}
            self.balls_bg = {}
            self.mutex.release()


    def is_same_team(self, ballId):
        if self.team_info == None:
            return False
        else:
            for mem in self.team_info:
                if mem[3] == ballId:
                    return True
        return False        


    def match_time(self):
        self.mutex.acquire()
        if self.cur_match_time != None and self.cur_match_time.GetDelay() < 10:
            self.mutex.release()
            return
        if self.next_match_time == None or int(time.time() * 1000) - self.next_match_time.local_send_time >= 1000:
            cmd = proto.message.Wilds.MatchTime.value
            msg = proto.message.pack(cmd, b"")
            if self.client != None:
                self.client.send(msg)
                if self.next_match_time == None:
                    self.next_match_time = matchtime.MatchTime(int(time.time() * 1000))
                else:
                    self.next_match_time.local_send_time = int(time.time() * 1000)
            # print("\n")
            # print("`````local_send_time:", self.next_match_time.local_send_time)
            # print("`````server_time:", self.next_match_time.server_time)
            # print("`````local_recv_time:", self.next_match_time.local_recv_time)
        self.mutex.release()
        
        
    def close(self):
        self.terminate = True
        if self.client != None:
            self.client.close()
            socket_mgr.get_socket_set("tcp").remove_sock(self.sockindex)
            self.client = None
        if gateway.gate.gateClient != None:
            gateway.gate.gateClient.close()
            gateway.gate.gateClient = None
        if self.udpclient != None:
            self.udpclient.close()
            socket_mgr.get_socket_set("udp").remove_sock(self.sockindex)
            self.udpclient = None
            