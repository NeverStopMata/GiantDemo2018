import os
import random
from socket import *
import ikcp
import threading
import time
import datetime


def handle_udp_packet(sock, kcp):
    while True:
        try:
            data = sock.recv(8*1024)
            dlen = len(data)
            if dlen <= 0:
                time.sleep(0.002)
            else:
                kcp.input(data)
            if dlen > 548:
                print("   --- udp recv data len:", dlen)
        except Exception as e:
            print(e)
            print("close handle_udp_packet. #1")
            return
    print("close handle_udp_packet. #2")

class KCPClient:
    def __init__(self):
        self.sock = socket(AF_INET, SOCK_DGRAM)
        self.conv = random.randint(1, 999999999)
        self.kcp = ikcp.IKcp(self.sock, self.conv, 2)
        self.kcp.wndsize(1024, 1024)

        self.udp_thread = None
        self.timer = None

    def connect(self, addr):
        self.sock.connect(addr)
        self.kcp.update(time.time() * 1000)
        self.udp_thread = threading.Thread(target=handle_udp_packet, args=(self.sock, self.kcp))
        self.udp_thread.start()

        self.timer = threading.Timer(0.02, self.on_timer)
        self.timer.start()

    def udp_send(self, buffer):
        if len(buffer) > 548:
            print("    --- udp send size:",len(buffer))
        self.sock.send(buffer)

    def on_timer(self):
        if self.kcp == None:
            return
        now = time.time() * 1000
        self.kcp.update(now)

        self.timer = threading.Timer(0.02, self.on_timer)
        self.timer.start()

    def send(self, buffer):
        if len(buffer) > 548:
            print("*** kcp send size:", len(buffer))
        self.kcp.send(buffer)

    def recv(self, size):
        data = self.kcp.recv(size)
        if data and len(data) > 548:
            print("*** kcp recv size:", len(data))
        return data
        
    def close(self):
        self.sock.close()
        self.kcp = None
        if self.timer != None:
            self.timer = None



