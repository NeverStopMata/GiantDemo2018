import argparse
import time
import threading
import os

args = None


def fun_timer1():
    url = "curl -s http://%s:8000?query=online > /dev/null" % (args.ip)
    os.system(url)
    timer1 = threading.Timer(1, fun_timer1)
    timer1.start()

def fun_timer2():
    url = "curl -s http://%s:8000?query=cpu > /dev/null" % (args.ip)
    os.system(url)
    timer2 = threading.Timer(1, fun_timer2)
    timer2.start()
    
def fun_timer3():
    url = "curl -s http://%s:8000?query=net > /dev/null" % (args.ip)
    os.system(url)
    timer3 = threading.Timer(1, fun_timer3)
    timer3.start()
    
def fun_timer4():
    url = "curl -s http://%s:8000?query=mem > /dev/null" % (args.ip)
    os.system(url)
    timer4 = threading.Timer(1, fun_timer4)
    timer4.start()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='curl crontab',formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--ip", default="122.11.58.163", help="curl ip", type=str)
    args = parser.parse_args()

    timer1 = threading.Timer(5, fun_timer1)
    timer1.start()
    
    timer2 = threading.Timer(5, fun_timer2)
    timer2.start()
    
    timer3 = threading.Timer(5, fun_timer3)
    timer3.start()
    
    timer4 = threading.Timer(5, fun_timer4)
    timer4.start()
    
    while True:
        time.sleep(100)
