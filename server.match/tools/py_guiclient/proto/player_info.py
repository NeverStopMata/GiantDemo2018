
class PlayerInfo():
    def __init__(self, player):
        self.id = player.id
        self.name = player.name
        self.local = player.local
        self.IsLive = player.IsLive
        self.SnapInfo = player.SnapInfo
        self.ballId = player.ballId
        self.curexp = player.curexp
        self.curmp = player.curmp
        self.animalid = player.animalid
        self.curhp = player.curhp
        self.TeamName = player.TeamName
        self.bombNum = player.bombNum
        self.hammerNum = player.hammerNum
        
    def print(self):
        print("==============================")
        print("id = ", self.id)
        print("name = ", self.name)
        print("local = ", self.local)
        print("IsLive = ", self.IsLive)
        print("SnapInfo = ", self.SnapInfo)
        print("ballId = ", self.ballId)
        print("curexp = ", self.curexp)
        print("curmp = ", self.curmp)
        print("animalid = ", self.animalid)
        print("curhp = ", self.curhp)
        print("TeamName = ", self.TeamName)
        print("")