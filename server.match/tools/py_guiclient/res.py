import json
import xml.etree.ElementTree as ET

class Res():
    def __init__(self, args, cfg, scene_id):
        self.args = args
        self.cfg = cfg
        self.scene_id = scene_id
        
        self.map = None
        self.food = None
        self.animal = None
        
        self.load()
    
    def load(self):
        # config/terrain/map.json
        try:
            mapfile = "%s/%d.json" % (self.cfg["terrain"], self.scene_id)
            f = open(mapfile, 'rt')
            self.map = json.loads(f.read())
            f.close()
        except Exception as e:
            print(e)
            exit(0)
        
        # config/xml/food.xml
        try:
            foodfile = "%s/food.xml" % (self.cfg["xml"])
            root = ET.parse(foodfile).getroot()
            
            self.food = {}
            for food in root.findall("food"): #<food mapid="1002">
                mapid = food.get("mapid")
                if mapid == str(self.scene_id):
                    for item in food.findall("item"): #<item id="1" type="11" size="0.2" mapnum="400"/>
                        self.food[int(float(item.get("id")))] = { "size": float(item.get("size")), "type": float(item.get("type")) }
                    break
        except Exception as e:
            print(e)
            exit(0)
        #print(self.food)
            
            
        # config/xml/animal.xml
        try:
            animalfile = "%s/animal.xml" % (self.cfg["xml"])
            root = ET.parse(animalfile).getroot()
            
            self.animal = {}
            for animal in root.findall("animal"): #<animal id="1" speedup="2" speedupinterval="1.5" hprecover="2" hp="65" attack="10" attackinterval="0.7" scale="0.306" divetime="17" diveinterval="2" eatRange="0.336" hitStopForce="1.5" hitStopTime="0.3">
                id = int(float(animal.get("id")))
                scale = float(animal.get("scale"))
                self.animal[id] = { "scale": scale }
        except Exception as e:
            print(e)
            exit(0)

        # config/xml/level.xml
        try:
            levelFile = "%s/level.xml" % (self.cfg["xml"])
            root = ET.parse(levelFile).getroot()
            self.levelcfg = {}
            for mapInfo in root.findall("level"):
                if int(mapInfo.get("mapid")) == self.scene_id:
                    for item in mapInfo.findall("item"):
                        id = int(item.get("id"))
                        self.levelcfg[id] = int(item.get("nextexp"))
        except Exception as e:
            print(e)
            exit(0)

    # level: 1,2,3...
    def get_next_level_exp(self, level):
        if level < 1:
            return 0
        if level in self.levelcfg:
            return self.levelcfg[level]
        return 0


            
g_res = None
def new(args, cfg, scene_id):
    global g_res
    g_res = Res(args, cfg, scene_id)
    return g_res