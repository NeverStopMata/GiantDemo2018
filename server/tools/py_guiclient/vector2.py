import math
import sys

class Vector():
    def __init__(self, x, y):
        self.x = x
        self.y = y
        
    def IncreaseBy(self, v):
        self.x = self.x + v.x
        self.y = self.y + v.y
        
    def ScaleBy(self, v):
        self.x = self.x * v
        self.y = self.y * v
        
    def NormalizeSelf(self):
        magn = self.Magnitude()
        if self.IsZero(magn):
            return 0
        self.x = self.x / magn
        self.y = self.y / magn
        return magn
        
    def Magnitude(self):
        return math.sqrt(self.x * self.x + self.y * self.y)
        
    def IsZero(self, v):
        return v < sys.float_info.epsilon