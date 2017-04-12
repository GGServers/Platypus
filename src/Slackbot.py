from slackclient import SlackClient
from Config import Config
import threading

config = Config()
channel = config.Get("slack_channel")
token = config.Get("slack_api_key")
sc = SlackClient(token)


class Bot:

    def Post(self, message, icon):
        return sc.api_call(
            "chat.postMessage", channel=config.Get("slack_channel"), text=message,
            username="Platypus", icon_emoji=icon)

    def SingleReport(self, name, hostname, status):
        if config.Get("enable_slackbot"):
            message = "%s (%s) has just gone %s!", (name, hostname, status)
            self.Post(message, ":exclamation:")
        else:
            return 0
