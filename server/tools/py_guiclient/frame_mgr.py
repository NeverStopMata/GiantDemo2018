import time
import math

class FrameMgr():
    def __init__(self):
        self.pre_frame = -1
        self.timestamp = 0
        
    def set_frame(self, frame):
        self.pre_frame = frame
        self.timestamp = int(round(time.time() * 1000))
        
    def get_pre_frame(self):
        return self.pre_frame
    