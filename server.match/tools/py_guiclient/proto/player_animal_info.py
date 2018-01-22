
class PlayerAnimalInfo():
    def __init__(self, animal):
        self.id = animal.id
        self.animalid = animal.animalid
        
    def print(self):
        print("==============================")
        print("id = ", self.id)
        print("animalid = ", self.animalid)
        print("")