#!/usr/bin/python

# Platypus panel script
# Runs a minimal HTTP/S server that
# returns JSON reflecting resource usage
# of the server. Attempts to use as few
# resources itself. Works on Windows and Linux.
from http.server import BaseHTTPRequestHandler, HTTPServer
import signal
import sys
import psutil
import json

# Scripts for retrieving actual server stats


def Stats():
    s = {}
    s["cpu"] = round(psutil.cpu_percent()) # Used CPU
    s["memory"] = round(psutil.virtual_memory().percent) # Used memory
    s["disk"] = round(psutil.disk_usage('/').percent) # Used disk
    return json.dumps(s)

# Actual minimal HTTP server. Source:
# https://daanlenaerts.com/blog/2015/06/03/create-a-simple-http-server-with-python-3/


class testHTTPServer_RequestHandler(BaseHTTPRequestHandler):

    # GET
    def do_GET(self):
        # Send response status code
        self.send_response(200)

        # Send headers
        self.send_header('Content-type', 'application/json')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.end_headers()

        # Send message back to client
        message = Stats()
        # Write content as utf-8 data
        self.wfile.write(bytes(message, "utf8"))
        return


def run():
    print('starting platypus client webserver')

    # Server settings
    server_address = ('127.0.0.1', 9000)
    httpd = HTTPServer(server_address, testHTTPServer_RequestHandler)
    print('running, listening on 127.0.0.1:9000')
    print('press ctrl+c to exit')
    httpd.serve_forever()


def signal_handler(signal, frame):
    print('^C, exiting')
    sys.exit(0)
signal.signal(signal.SIGINT, signal_handler)

run()
