from socket import *
import select
import threading
import time
import log

print=log.log

BUFSIZ = 1024 * 1024
socket_set_map = {}


class SocketSet():
    def __init__(self):
        self.on_recv_cb = None
        self.thread = None
        self.mutex = threading.Lock()
        self.socks_map_i2s = {}
        self.socks_map_s2i = {}
        self.termination = False

    def init(self, on_recv):
        self.on_recv_cb = on_recv

    def add_sock(self, sockindex, sock_type=SOCK_STREAM):
        s = socket(AF_INET, sock_type)
        self.mutex.acquire()
        self.socks_map_i2s[sockindex] = s
        self.socks_map_s2i[s] = sockindex
        self.mutex.release()
        return s
    
    def remove_sock(self, index):
        self.mutex.acquire()
        if index in self.socks_map_i2s:
            sock = self.socks_map_i2s[index]
            self.socks_map_i2s.pop(index)
            self.socks_map_s2i.pop(sock)
        self.mutex.release()

    def run(self):
        while not self.termination:
            try:
                socks = []
                socks2 = {}
                self.mutex.acquire()
                for index, sock in self.socks_map_i2s.items():
                    socks.append(sock)
                    socks2[sock] = index
                self.mutex.release()
                if len(socks) == 0:
                    time.sleep(0.01)
                    continue
                inputs, outputs, exceptions = select.select(socks, [], [])
                for indata in inputs:
                    try:
                        data = indata.recv(BUFSIZ)
                        if len(data) > 0:
                            self.on_recv_cb(socks2[indata], data)
                    except Exception as e:
                        print(e)
                for e in exceptions:
                    print("error :", e)
            except Exception as e:
                print(e)
            time.sleep(0.01)
        self.mutex.acquire()
        for _, sock in self.socks_map_i2s.items():
            sock.close()
        self.socks_map_i2s = {}
        self.socks_map_s2i = {}
        self.mutex.release()
        print("exit SocketSet::run")
        
    def run_in_thread(self):
        self.thread = threading.Thread(target=self.run)
        self.thread.start()

    def close(self):
        print("SocketSet::close")
        self.mutex.acquire()
        for _, sock in self.socks_map_i2s.items():
            sock.close()
        self.socks_map_i2s = {}
        self.socks_map_s2i = {}
        self.mutex.release()
        self.termination = True
        if self.thread is not None:
            self.thread.join()
            self.thread = None


def get_socket_set(idx="tcp"):
    global socket_set_map
    sock_set = socket_set_map.get(idx)
    if sock_set is None:
        sock_set = SocketSet()
        socket_set_map[idx] = sock_set
    return sock_set


def close():
    global socket_set_map
    for _,sockset in socket_set_map.items():
        sockset.close()