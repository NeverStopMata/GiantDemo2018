import wx
import math


CELL_SIZE = 50


class Camera():
    def __init__(self):
        self.scale = 1
        
    def get_cell_size(self):
        return self.scale * CELL_SIZE

    def get_view(self, view_size, player_pos, map_size):
        viewdata = {}
        viewdata["width"] = view_size.width
        viewdata["height"] = view_size.height
        
        player_pos.x = player_pos.x * self.scale
        player_pos.y = player_pos.y * self.scale
        map_size.width = map_size.width * self.scale
        map_size.height = map_size.height * self.scale

        center = wx.Point()
        center.x = viewdata["width"] / 2
        center.y = viewdata["height"] / 2

        if player_pos.x <= center.x:
            viewdata["xbegin"] = 0
            viewdata["xend"] = center.x * 2 + 1
        elif map_size.width - player_pos.x <= center.x:
            viewdata["xbegin"] = map_size.width - center.x * 2
            viewdata["xend"] = map_size.width
        else:
            viewdata["xbegin"] = player_pos.x - center.x
            viewdata["xend"] = player_pos.x + center.x + 1

        if player_pos.y <= center.y:
            viewdata["ybegin"] = 0
            viewdata["yend"] = center.y * 2 + 1
        elif map_size.height - player_pos.y <= center.y:
            viewdata["ybegin"] = map_size.height - center.y * 2
            viewdata["yend"] = map_size.height
        else:
            viewdata["ybegin"] = player_pos.y - center.y
            viewdata["yend"] = player_pos.y + center.y + 1
        
        viewdata["xbegin"] = math.ceil(viewdata["xbegin"])
        viewdata["xend"] = math.ceil(viewdata["xend"])
        viewdata["ybegin"] = math.ceil(viewdata["ybegin"])
        viewdata["yend"] = math.ceil(viewdata["yend"])
        
        return viewdata
    