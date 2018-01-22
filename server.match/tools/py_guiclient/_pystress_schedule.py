import time
import json
import random
import uuid
from user import *
import threading

class ScheduleBase(object):
    def __init__(self, args, cfg):
        self.args = args
        self.cfg = cfg
        self.second = 0
        self.clients = {}
        self.cmds = {}
        self.mutex = threading.Lock()
        
    # thread #1
    def start(self):
        for i in range(1, self.cfg["account_count"] + 1):
            self.login_account(i)
        print("login end ===========================.")
    
    
    def login_account(self, sockindex):
        user = User(sockindex, self.args, self.cfg)
        user.pystress = True
        accountname = str(uuid.uuid1())
        done = True
        if user.login_detail("", "", accountname, "android", 1, 0) == True:
            if user.start_game() == False:
                done = False
            else:
                self.mutex.acquire()
                self.clients[sockindex] = user
                self.mutex.release()
        else:
            done = False
            
        if done == False:
            time.sleep(0.02)
            print("login account fail! sockindex = ", sockindex)
            # relogin
            user.close()
            user = None
            self.login_account(sockindex)
    
    # thread #1
    def update(self):
        pass
        '''
        self.second = self.second + 1
        (_, tmp) = divmod(self.second, 30)
        if tmp == 0:
            for _, user in self.clients.items():
                user.heart_beat()
        '''
        
    # thread #2
    def on_udprecv(self, sockindex, data):
        self.on_recv_detail(sockindex, data, 1)
    
    # thread #3
    def on_recv(self, sockindex, data):
        self.on_recv_detail(sockindex, data, 0)

    def on_recv_detail(self, sockindex, data, udp_channel = 0):
        self.mutex.acquire()
        if sockindex not in self.clients:
            self.mutex.release()
            return
        user = self.clients[sockindex]
        self.mutex.release()
        
        msgs = proto.message.on_recv(user, data, udp_channel)
        if len(msgs) == 0:
            return
        
        for (cmd, d) in msgs:
            if cmd in self.cmds:
                self.cmds[cmd](user, cmd, d)
            else:
                user.on_recv_one(cmd, d, udp_channel)
                
            
class Schedule0(ScheduleBase):
    def __init__(self, args, cfg):
        ScheduleBase.__init__(self, args, cfg)
        self.cmds[proto.message.Wilds.EndRoom.value] = self.on_end_room
    
    # thread #1
    def update(self):
        ScheduleBase.update(self)
        
    # thread #2
    def on_end_room(self, user, cmd, data):
        print("on_exit_room, sockindex = ", user.sockindex)
        user.on_exit_room_clear_data()
        print("close, sockindex = ", user.sockindex)
        user.close()
        print("start_game, sockindex = ", user.sockindex)
        user.start_game()
        