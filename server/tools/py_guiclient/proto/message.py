import struct
import hashlib
from enum import Enum
import zlib

LOGIN_SIGN="#@$#2ds&^%&#298"

def get_sign(cmd):
    sign = LOGIN_SIGN+str(cmd)
    return hashlib.md5(sign.encode('utf-8')).hexdigest()


'''
数据块格式
  * 3个字节 数据块大小
  * 1个字节标志位： 0未压缩；1 zlib压缩
  * 数据块：2个字节消息ID；剩余protobuf消息块
'''

DATA0_LEN = 6
CMD_LEN = 2


def unpack_header(rawdata):
    l1, l2, l3, flag, cmd = struct.unpack("<BBBBH", rawdata[:DATA0_LEN])
    l = l1 + (l2<<8) + (l3<<16)
    return l, flag, cmd

def unpack(rawdata):
    l, flag, cmd = unpack_header(rawdata)
    data = rawdata[DATA0_LEN:]
    if l != len(data) + CMD_LEN: # len = cmd + data
        print("l1=", l1, " l2=", l2, "l3=", l3, "flag=", flag, "cmd=", cmd)
        print("l=", l, " len(data)=", len(data))
        exit(0)
        
    return cmd, data
    
    
def pack(cmd, data):
    l = len(data) + CMD_LEN
    l1 = l & 0xFF
    l2 = (l>>8) & 0xFF
    l3 = (l>>16) & 0xFF
    flag = 0
    return struct.pack("<BBBBH", l1, l2, l3, flag, cmd) + data


def on_recv(user, data, udp_channel=0):
    msgs = []
    if len(data) == 0:
        return msgs

    left_data = None
        
    if udp_channel == 1:
        left_data = user.left_data_udp
    else:
        left_data = user.left_data_tcp

    rawdata = left_data + data
    while True:
        if len(rawdata) > DATA0_LEN:
            l, flag, cmd = unpack_header(rawdata)
            data = rawdata[DATA0_LEN:]
            if l <= len(data) + CMD_LEN:
                if flag == 1:
                    msgs.append((cmd, zlib.decompress(data[:l - CMD_LEN])))
                else:
                    msgs.append((cmd, data[:l - CMD_LEN]))
                rawdata = data[l - CMD_LEN:]
            else:
                break
        else:
            break
    left_data = rawdata

    if udp_channel == 1:
        user.left_data_udp = left_data
    else:
        user.left_data_tcp = left_data

    return msgs

def on_udprecv(data):
    msgs = []
    cmd, = struct.unpack("<H", data[:2])
    msgs.append((cmd, data[2:]))
    return msgs
    

class Player(Enum):
    Login = 1
    GateAddr = 2
    ReqIntoFRoom = 3
    
    
class Wilds(Enum):
    Login=1
    Top=2
    AddPlayer=3
    RemovePlayer=4
    Scene=5
    Move=6
    Run=7
    ReLife=9
    Death=10
    EndRoom=11
    RefreshPlayer=13
    AsyncPlayerAnimal=14
    HeartBeat=16
    SceneChat=20
    ActCloseSocket=21
    SceneTCP          = 31
    SceneUDP          = 32
    BindTCPSession	  = 33
    VoiceInfo=36
    TeamRankList           = 40
    CastSkill=50
    UpdateTeamInfo  = 77

class Team(Enum):
    Login = 1
    InviteList = 5
    StartGame = 6
    JoinTeam = 11
    CreateTeam = 14
    CreateRoom = 16
    TReConnect = 29
    BReConnect = 207


class Common(Enum):
    WorldChat = 4

class Chat(Enum):
    WorldChat = 22