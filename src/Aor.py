import tornado.ioloop
import tornado.websocket
import json
import ServerHandler
import Slackbot

bot = Slackbot.Bot()
sql = ServerHandler.Sql()
cache = ServerHandler.Cache()
AliveSockets = set()


class AorMaster(tornado.websocket.WebSocketHandler):

    def open(self):
        print(self.request.remote_ip + " connected, asking for authentication...")
        self.write_message('{"success": true, "require_auth": true}')
        self.ip = self.request.remote_ip

        d = sql.IpToId(self.ip)
        j = json.dumps({
            "id": d[0],
            "online": True,
            "disk": "refreshing",
            "cpu": "refreshing",
            "memory": "refreshing"
        })
        cache.TriggerOnline(str(d[0]))
        SendFetchMessage(j)
        bot.SingleReport(d[1], d[2], "online")

    def on_message(self, message):
        raw = json.loads(message)

        if sql.Verify(raw["uuid"]):
            #   print(self.ip, "provided vaid uuid", raw["uuid"])
            d = sql.UuidToId(raw["uuid"])
            # cache.Update(id, raw["stats"])
            j = json.dumps({
                "id": d[0],
                "online": True,
                "disk": raw["stats"]["disk"],
                "cpu": raw["stats"]["cpu"],
                "memory": raw["stats"]["memory"]
            })
            SendFetchMessage(j)
            self.write_message('{"success": true, "require_auth": true}')

        else:
            print(self.ip, "provided INVALID uuid", raw["uuid"])
            self.write_message('{"success":false,"errcode":"invalid_uuid"}')

    def on_close(self):
        d = sql.IpToId(self.request.remote_ip)
        j = json.dumps({
            "id": d[0],
            "online": False,
            "disk": 0,
            "cpu": 0,
            "memory": 0
        })
        SendFetchMessage(j)
        bot.SingleReport(d[1], d[2], "offline")
        cache.TriggerOffline(str(d[0]))


class FetchWebsocket(tornado.websocket.WebSocketHandler):

    def open(self):
        AliveSockets.add(self)
        print("Fetch websocket opened - client " + self.request.remote_ip)

    def on_message(self, message):
        print(self.request.remote_ip + " requested panel " + message)
        # res = scan.Fetch(message)
        res = "placeholder"
        self.write_message(res)
        print(self.request.remote_ip + " was sent " + message)

    def on_close(self):
        AliveSockets.remove(self)
        print("Fetch websocket closed - client " + self.request.remote_ip)


def SendFetchMessage(message):
    for ws in AliveSockets:
        ws.write_message(message)
