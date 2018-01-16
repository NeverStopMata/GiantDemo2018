
import threading
import proto.message
import proto.team_pb2
import urllib.request
from socket import *
import user
import proto.chat_pb2
import mainwindow


BUFSIZ = 128 * 1024

def gateClientThread(gate, user):
    print("start gate client thread")
    while user.terminate == False:
        data = None
        try:
            data = gate.gateClient.recv(BUFSIZ)
        except Exception as e:
            print(e)
            gate.gateClient = None
            print("close gate client thread. #2")
            return
        if data != None:
            gate.on_gate_recv(data)
    print("close gate client thread. #1")


class Gateway:

    def __init__(self):
        self.gateClient = None
        self.gateClientThread = None
        self.user = None
        self.left_data_tcp = b""
        self.left_data_udp = b""
        self.cmds = {}
        self.init_cmds()
        
    def init_cmds(self):
        self.cmds[proto.message.Team.Login.value] = self.on_gate_login 
        self.cmds[proto.message.Team.CreateTeam.value] = self.on_create_team
        self.cmds[proto.message.Team.BReConnect.value] = self.on_breconnect
        self.cmds[proto.message.Team.InviteList.value] = self.on_invite_list
        self.cmds[proto.message.Team.CreateRoom.value] = self.on_create_room
        self.cmds[proto.message.Team.TReConnect.value] = self.on_treconnect
        self.cmds[proto.message.Team.StartGame.value] = self.on_start_game
        self.cmds[((proto.message.Common.WorldChat.value << 8) | proto.message.Chat.WorldChat.value)] = self.on_chat

    def login_gateway(self, user):
        c = proto.message.Player.GateAddr.value
        sign = proto.message.get_sign(c)
        url = "http://%s:%d/msg?c=%d&sign=%s" % (user.addr, user.port, c, sign)
        print("url=", url)

        request = urllib.request.Request(url, headers={"Cookie":"session="+user.cookie})
        result = urllib.request.urlopen(request).read()

        cmd, msg = proto.message.unpack(result)
        pmsg = proto.player_pb2.RetGateAddr()
        pmsg.ParseFromString(msg)

        print("gateaddr: ", pmsg.Address)
        print("Key: ", pmsg.Key)

        host = pmsg.Address[:pmsg.Address.find(":")]
        port = int(pmsg.Address[pmsg.Address.find(":") + 1:])
        addr = (host, port)

        self.gateClient = socket(AF_INET, SOCK_STREAM)
        try:
            self.gateClient.connect(addr)
        except Exception as e:
            print(e)
            return False

        self.gateClientThread = threading.Thread(target=gateClientThread, args=(self,user))
        self.gateClientThread.start()

        loginMsg = proto.team_pb2.ReqTeamLogin()
        loginMsg.Name = "%s_%d" % (user.account_name, user.id)
        loginMsg.Key = pmsg.Key
        loginMsg.ClientVer = ""
        data = loginMsg.SerializeToString()
        cmd = proto.message.Team.Login.value
        msg = proto.message.pack(cmd, data)
        if self.gateClient != None:
            self.gateClient.send(msg)

        self.user = user
        return True

    def on_gate_recv(self, data):
        msgs = proto.message.on_recv(self, data)
        if len(msgs) == 0:
            return

        for (cmd, data) in msgs:
            if cmd in self.cmds:
                self.cmds[cmd](data)
            else:
                print("gate cmd: ", cmd)


    def on_gate_login(self, data):
        print("login gateway success")
        #self.create_team()


    def create_team(self):
        # TeamCmd_CreateTeam
        cmd = proto.message.Team.CreateTeam.value
        msg = proto.message.pack(cmd, b"")
        if self.gateClient != None:
            self.gateClient.send(msg)
    
    def on_create_team(self, data):
        print("create team success")
        self.create_room()

    def on_breconnect(self, data):
        print("create team return breconnect")
        pmsg = proto.team_pb2.RetBReConnect()
        pmsg.ParseFromString(data)
        print("RoomType: ", pmsg.RType, " EndTime: ", pmsg.EndTime, " RoomId: ", pmsg.RoomId)
        self.join_team(pmsg.RoomId)

    def on_invite_list(self, data):
        print("on_invite_list")

    def create_room(self):
        cmd = proto.message.Team.CreateRoom.value
        pmsg = proto.team_pb2.ReqCreateRoom()
        pmsg.Name = "my room"
        pmsg.Model = 1
        pmsg.Priv = 1
        data = pmsg.SerializeToString()
        msg = proto.message.pack(cmd, data)
        if self.gateClient != None:
            self.gateClient.send(msg)
    
    def on_create_room(self, data):
        print("create room success")
        self.start_game()

    def on_treconnect(self, data):
        print("create room return breconnect")
        pmsg = proto.team_pb2.RetTReConnect()
        pmsg.ParseFromString(data)
        print("Model: ", pmsg.Model, " RoomName: ", pmsg.RoomName, " UserNum: ", pmsg.UserNum, " EndTime: ", pmsg.EndTime)

    def start_game(self):
        cmd = proto.message.Team.StartGame.value
        msg = proto.message.pack(cmd, b"")
        if self.gateClient != None:
            self.gateClient.send(msg)

    def on_start_game(self, data):
        pmsg = proto.team_pb2.RetStartGame()
        pmsg.ParseFromString(data)
        print("game start, addr: ", pmsg.Address, " Key: ", pmsg.Key)
        self.user.room_server_addr = pmsg.Address
        self.user.room_server_key = pmsg.Key
        self.user.connect_room_server()

    def join_team(self, RoomId):
        cmd = proto.message.Team.JoinTeam.value
        pmsg = proto.team_pb2.ReqJoinTeamGame()
        pmsg.RoomId = RoomId
        data = pmsg.SerializeToString()
        msg = proto.message.pack(cmd, data)
        if self.gateClient != None:
            self.gateClient.send(msg)

    def SendChat(self, text):
        print("begin SendChat， text =", text)
        cmd = ((proto.message.Common.WorldChat.value << 8) | proto.message.Chat.WorldChat.value)
        pmsg = proto.chat_pb2.ReqChatPrivate()
        pmsg.ToUserId = self.user.id
        pmsg.FromUserId = self.user.id
        pmsg.Text = text
        pmsg.TypeId = 0
        data = pmsg.SerializeToString()
        msg = proto.message.pack(cmd, data)
        if self.gateClient != None:
            self.gateClient.send(msg)

    def on_chat(self, data):
        pmsg = proto.chat_pb2.RetChatPrivate()
        pmsg.ParseFromString(data)
        print("on chat ", pmsg.Text)
        if mainwindow.g_mainwin != None:
            mainwindow.g_mainwin.on_net_chat(pmsg.Text)


gate = Gateway()