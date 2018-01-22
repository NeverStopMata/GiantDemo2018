import argparse
import json
from _pystress_schedule import *
import time
import socket_mgr
import log

if __name__ == "__main__":
    log.Enable=False

    parser = argparse.ArgumentParser(description='pystress tool',formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--cfg", default="_cfg.json", help="cfg file path", type=str)
    args = parser.parse_args()

    f = open(args.cfg, 'rt')
    content = f.read()
    f.close()
    cfg = json.loads(content)
    
    sche = None
    if cfg["test_type"] == 0:
        sche = Schedule0(args, cfg)
    else:
        print("unknow test_type!!!")
        exit(0)
    
    socket_mgr.get_socket_set("tcp").init(sche.on_recv)
    socket_mgr.get_socket_set("tcp").run_in_thread()

    if cfg["room"]["udp_enable"] != 0:
        socket_mgr.get_socket_set("udp").init(sche.on_udprecv)
        socket_mgr.get_socket_set("udp").run_in_thread()
    
    sche.start()
    
    while True:
        time.sleep(1)
        sche.update()