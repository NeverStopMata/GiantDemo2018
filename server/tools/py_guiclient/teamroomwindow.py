import roomwindow
import copy

class TeamRoomWindow(roomwindow.RoomWindow):
    def __init__(self, user, args, cfg):
        roomwindow.RoomWindow.__init__(self, user, args, cfg)

    def draw(self, detal):
        roomwindow.RoomWindow.draw(self, detal)

    def right_top_data(self):
        if self.user.team_rank_list != None:
            self.EndTime = self.user.team_rank_list.EndTime
        else:
            self.EndTime = 0
        return (copy.copy(self.user.team_rank_list), copy.copy(self.user.team_info))

    def draw_right_top(self, data):
        self.listTop.Clear()
        if data[0] != None:
            self.listTop.Append("队伍排名")
            for team in data[0].Teams:
                self.listTop.Append("team(%s) num(%d) corpname(%s) score(%f) lastrank(%d)" % (team.Tname,team.Num,team.CorpName,team.Score,team.LastRank))

            self.listTop.Append("")
            self.listTop.Append("watchnum(%d) endtime(%d) roomexp(%d)" % (self.user.team_rank_list.WatchNum, self.user.team_rank_list.EndTime, self.user.team_rank_list.killNum))

        if data[1] != None:
            self.listTop.Append("")
            self.listTop.Append("")
            self.listTop.Append("队伍成员位置")
            for mem in self.user.team_info:
                self.listTop.Append("member(%d) pos(%f, %f)" % (mem[0],mem[1],mem[2]))
        

g_roomwin = None      
def new(user, args, cfg):
    global g_roomwin
    g_roomwin = TeamRoomWindow(user, args, cfg)
    g_roomwin.Show()
    
def delete():
    global g_roomwin
    if g_roomwin != None:
        g_roomwin.Close()
        g_roomwin = None