import argparse
import os
import json
import user as muser
import mainwindow
import socket_mgr


def parse_args():
    parser = argparse.ArgumentParser(description='py_guiclient', formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--cfg", default="./_cfg.json", help="configure file", type=str)
    args = parser.parse_args()

    jsonfile = os.path.abspath(os.path.dirname(__file__)) + "/" + args.cfg
    if os.path.exists(jsonfile) == False:
        print("ERROR: cfg file not found. path = ", jsonfile)
        exit(0)

    f = open(jsonfile, 'rt')
    cfg = json.loads(f.read())
    f.close()
    return args, cfg
    

users = {}
def on_recv(sockindex, data):
    global users
    user = users[sockindex]
    user.on_recv(data, 0)

def on_udprecv(sockindex, data):
    global users
    user = users[sockindex]
    user.on_recv(data, 1)

    
if __name__ == "__main__":
    socket_mgr.get_socket_set("tcp").init(on_recv)
    socket_mgr.get_socket_set("tcp").run_in_thread()

    socket_mgr.get_socket_set("udp").init(on_udprecv)
    socket_mgr.get_socket_set("udp").run_in_thread()
    
    args, cfg = parse_args()
    sockindex = 1
    user = muser.User(sockindex, args, cfg)
    users[sockindex] = user
    
    if user.login() == False:
        print("request login fail.")
        exit(0)
    
    mainwindow.run(user, args, cfg)
    
    socket_mgr.close()
    print("over.")
    
    