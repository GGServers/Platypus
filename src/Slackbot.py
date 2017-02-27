from slackclient import SlackClient
from src.Cache import Handler
from src.Config import Config
import threading

config = Config()
channel = config.Get("slack_channel")
token = config.Get("slack_api_key")
sc = SlackClient(token)
handler = Handler()

class Bot:
    def Post(self,message, channel, username, icon):
        return sc.api_call(
            "chat.postMessage", channel=channel, text=message,
            username=username, icon_emoji=icon)


    def BuildMessage(self,data):
        post = True
        off = 0
        channel = config.Get("slack_channel")
        username = "Platypus"
        icon = ":desktop_computer:"
        message = "Some panels may be offline!"
        for s in data:
            if s[4] != 1:
                message = message + " " + s[1] + " (" + s[2] + ")"
                off = off + 1

            if off > 1: post = True
            else: post = False

        if post is True: self.Post(message, channel, username, icon)


    def Data(self):
        data = handler.Get(offline="only")
        self.BuildMessage(data)

    def Loop(self):
        self.Data()
        threading.Timer(config.Get("slack_interval"), self.Loop).start()
