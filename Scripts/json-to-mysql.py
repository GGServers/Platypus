# Used for migrating to MariaDB/MySQL database,
# commited just in case you want to use it yourself.
import json
import MySQLdb

db = MySQLdb.connect(user="root", db="platypus")
c=db.cursor()
servers = json.load(open("src/cache/servers.json"))

for s in servers:
    c.execute("INSERT INTO server (id,name,hostname,location) VALUES (%s, %s, %s, %s)",
        (s['id'], s['name'], s['hostname'], s['location']))

db.commit()