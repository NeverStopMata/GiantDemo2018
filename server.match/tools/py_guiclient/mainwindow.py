import wx
import roomwindow
import teamroomwindow
import user as muser
from gateway import *

CLIENT_WIDTH = 400
CLIENT_HEIGHT = 300

class MainWindow(wx.Frame):
    def __init__(self, user, args, cfg):
        wx.Frame.__init__(self, None, style = wx.CAPTION | wx.CLOSE_BOX)
        self.user = user
        self.cfg = cfg
        self.args = args
        self.close_flag = 0
        
        self.init()

    def init(self):
        self.SetTitle("py_guiclient")
        self.Bind(wx.EVT_CLOSE, self.OnClose)
        
        self.SetClientSize((CLIENT_WIDTH, CLIENT_HEIGHT))
        
        temp = 25
        self.panel1 = wx.Panel(self, -1, pos=(0, 0), size=(CLIENT_WIDTH, 100))
        wx.StaticText(self.panel1, -1, "游戏模式：", pos=wx.Point(5, 5))
        self.btnStartGame = wx.Button(self.panel1, -1, "开始比赛", pos=(0, temp), size=(self.panel1.GetSize().width/2, self.panel1.GetSize().height-temp))
        self.Bind(wx.EVT_BUTTON, self.OnStartGame, self.btnStartGame)
        
        self.btnTeam = wx.Button(self.panel1, -1, "团队模式", pos=(self.panel1.GetSize().width/2, temp), size=(self.panel1.GetSize().width/2, self.panel1.GetSize().height-temp))
        self.Bind(wx.EVT_BUTTON, self.OnTeam, self.btnTeam)
        
        self.panel2 = wx.Panel(self, -1, pos=(0, self.panel1.GetSize().height), size=(CLIENT_WIDTH, CLIENT_HEIGHT - self.panel1.GetSize().height))
        
        wx.StaticText(self.panel2, -1, "聊天频道（附近）：", pos=wx.Point(5, 5))
        self.listChat = wx.ListBox(choices=[], parent=self.panel2, pos=wx.Point(0, temp),size=wx.Size(CLIENT_WIDTH, 140), style=0)
        self.txtChat = wx.TextCtrl(self.panel2, -1, "", pos=(95, 172), size=(CLIENT_WIDTH - 100, -1))
        self.btnChat = wx.Button(self.panel2, -1, "发送", pos=(0, 170), size=(-1, -1))
        self.Bind(wx.EVT_BUTTON, self.OnChat, self.btnChat)
        
        self.Centre()
        
    def OnStartGame(self, evt):
        self.user.room_type = 1
        self.user.start_game()
        
    def OnTeam(self, evt):
        self.user.room_type = 2
        # TODO
        
    def OnChat(self, evt):
        v = self.txtChat.GetValue()
        if v != "":
            gate.SendChat(v)
            
    def on_net_chat(self, text):
        wx.CallAfter(self.listChat.Append, text)
        

    def OnClose(self, evt):
        print("OnClose #1")
        if self.close_flag == 0:
            print("OnClose #2")
            self.user.close()
        evt.Skip()
        

g_mainwin = None
def new(user, args, cfg):
    global g_mainwin
    g_mainwin = MainWindow(user, args, cfg)
    g_mainwin.Show()
    
def run(user, args, cfg):
    app = wx.App(False)
    new(user, args, cfg)
    app.MainLoop()
    
    
def open_room_window(user):
    global g_mainwin
    if g_mainwin != None:
        g_mainwin.close_flag = 1
        g_mainwin.Close()
        g_mainwin = None
    if user.room_type == 1:        
        wx.CallAfter(roomwindow.new, user, user.args, user.cfg)
    else:
        wx.CallAfter(teamroomwindow.new, user, user.args, user.cfg)
        
        
def close_room_window(user):
    if user.room_type == 1:        
        wx.CallAfter(roomwindow.delete)
    else:
        wx.CallAfter(teamroomwindow.delete)
    wx.CallAfter(new, user, user.args, user.cfg)
