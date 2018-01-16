import time

class MatchTime:
    def __init__(self, localSendTime):
        self.local_send_time = localSendTime
        self.server_time = 0
        self.local_recv_time = 0

    def GetServerTime(self):
        return self.server_time + int((self.local_recv_time - self.local_send_time) / 2)

    def GetClientTime(self):
        return self.local_recv_time

    def SetServerTime(self, serverTime):
        self.server_time = serverTime
        self.local_recv_time = int(time.time() * 1000)

    def GetDelay(self):
        return self.local_recv_time - self.local_send_time

    def GetTimeDelta(self, cliTime, svrTime):
        return cliTime - self.local_recv_time - svrTime + self.server_time

    def ToClientTime(self, svrTime):
        return self.local_recv_time - self.server_time + svrTime