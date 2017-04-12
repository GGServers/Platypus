import requests
import threading
# from src.Cache import Handler
import ServerHandler
import time
from Config import Config

config = Config()
sql = ServerHandler.Sql()

class Scan:
    def __init__(self):
        self.panels = sql.Get()

    def Fetch(self,panel=None):
        if panel == None:
            print("Fetching panels")
            for p in self.panels:
                result = self.Check(p)
                #print(result)
                sql.Set(p[0],
                    result['online'],
                    str(result['cpu']),
                    str(result['memory']),
                    str(result['disk']))

        else:
            panel = sql.Get("one", panel)[0]
            #print(panel)
            return self.Check(panel)

    def Check(self, panel):
        id = panel[0]

        try:
            request = requests.get("http://" + panel[2] + config.Get("stats_path"),
                                   timeout = config.Get("scan_timeout"))

            #print(panel[0], "online", request.status_code)
            if request.status_code == 404:
                return {"id": id,
                        "name": panel[1],
                        "online": True,
                        "cpu": 0,
                        "memory": 0,
                        "disk": 0}
            else:
                data = request.json()
                return {"id": id,
                        "name": panel[1],
                        "online": True,
                        "cpu": data["cpu"],
                        "memory": data["memory"],
                        "disk": data["hdd"]}

        except Exception as e:
            #print(panel[0], "offline")
            #print(e)
            return {    "id": id,
                        "name": panel[1],
                        "online": False,
                        "cpu": 0,
                        "memory": 0,
                        "disk": 0}

    def Loop(self):
        self.Fetch()
        threading.Timer(config.Get("scan_interval"), self.Loop).start()
