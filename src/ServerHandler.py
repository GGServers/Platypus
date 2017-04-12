# This module manages the data on the servers. It:
# - (sql) Handles server registration w/ uid
# - (sql) Verifies server uid
# - (sql) Deletes servers
# - (sql) Get server metadata from database
# - (cache) Store server stats for API / non-ws suported browsers
# - (cache) Delivers server stats

import MySQLdb
import Config
import uuid
import json


class Sql:
    # Initialize database & variables

    def __init__(self):
        self.config = Config.Config()

        self.sqluser = self.config.Get("sql_user")
        self.sqlpass = self.config.Get("sql_pass")
        self.sqltable = self.config.Get("sql_db")

        self.db = MySQLdb.connect(
            user=self.sqluser, passwd=self.sqlpass, db="platypus")
        self.c = self.db.cursor()

    # Check connection with database
    def CheckConnection(self):
        try:
            self.db.ping()
        # Else reconnect
        except:
            self.db = MySQLdb.connect(
                user=self.sqluser, passwd=self.sqlpass, db="platypus")

    # Register server with IP and uid
    def Register(self, name, hostname, ip):
        self.CheckConnection()

        uid = uuid.uuid4()

        self.c.execute("INSERT INTO %s (id,name,hostname,ip,uid) VALUES (%s,%s,%s,%s,%s)",
                       (self.sqltable, id, name, hostname, ip, uid))
        self.db.commit()

        print("Server registered", uid)
    # Verify server has correct uid

    def Verify(self, uid):
        self.CheckConnection()

        self.c.execute("SELECT * FROM `" + self.sqltable +
                       "` WHERE uuid='%s'" % uid)
        dbdata = self.c.fetchone()

        if dbdata is None:
            return False
        if dbdata[10] == uid:
            return True
        else:
            return False

    # Get server data from database
    def Get(self, id="*"):

        self.CheckConnection()
        if id != "*":
            self.c.execute(
                "SELECT id,name,hostname,uuid FROM `" + self.sqltable + "` WHERE id=%s" % id)
            dbdata = self.c.fetchone()
        else:
            print(self.sqltable)
            self.c.execute("SELECT id,name,hostname,uuid FROM `" +
                           self.sqltable + "`")
            dbdata = self.c.fetchall()

        return dbdata

    def UuidToId(self, uid):
        self.CheckConnection()

        self.c.execute("SELECT id,name,hostname FROM `" + self.sqltable +
                       "` WHERE uuid='%s'" % uid)
        dbdata = self.c.fetchone()

        return dbdata

    def IpToId(self, ip):
        self.CheckConnection()

        self.c.execute("SELECT id,name,hostname FROM `" +
                       self.sqltable + "` WHERE ip='%s'" % ip)
        dbdata = self.c.fetchone()

        return dbdata


class Cache:

    def __init__(self):
        with open("src/cache/data.json", "r") as data_file:
            self.cache = json.load(data_file)

    def Update(self, id, stats):
        id = str(id)
        self.cache[id] = stats
        self.__dump()

    def TriggerOffline(self, id):
        self.cache[id]["online"] = False
        self.__dump()

    def TriggerOnline(self, id):
        self.cache[id]["online"] = True
        self.__dump()

    def Fetch(self):
        return self.cache

    def __dump(self):
        with open('src/cache/data.json', 'w+') as data_file:
            json.dump(self.cache, data_file, indent=4)
