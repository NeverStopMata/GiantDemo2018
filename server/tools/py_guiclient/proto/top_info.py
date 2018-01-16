import proto.player_info


class TopInfo():
    def __init__(self, pmsg):
        self.players = {}
        for player in pmsg.players:
            p = proto.player_info.PlayerInfo(player)
            self.players[p.id] = p
        self.EndTime = pmsg.EndTime
        self.Rank = pmsg.Rank
        self.KillNum = pmsg.KillNum
        
    def print(self):
        print("==============================")
        print("players = ", self.players)
        print("EndTime = ", self.EndTime)
        print("Rank = ", self.Rank)
        print("KillNum = ", self.KillNum)
        print("")