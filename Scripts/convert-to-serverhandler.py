import MySQLdb
import json
import uuid

db = MySQLdb.connect(user="root", db="platypus")
c = db.cursor()

c.execute("SELECT * FROM server")
dbdata = c.fetchall()

servers = {}

with open("src/cache/data.json") as data_file:
    cache = json.load(data_file)


c.execute(
    "ALTER TABLE `server` ADD COLUMN IF NOT EXISTS `ip` VARCHAR(50) NULL," +
    "ADD COLUMN IF NOT EXISTS  `uuid` VARCHAR(50) NULL")

db.commit()

for server in dbdata:
    id = server[0]
    servers[id] = {
        "online": True
    }

    uid = str(uuid.uuid4())
    print(id, server[2], "registred with uuid", uid)
    c.execute("UPDATE server SET uuid=%s where id=%s", (uid, id))
    db.commit()

with open("src/cache/data.json", "w") as data_file:
    json.dump(servers, data_file, indent=4)

print("Done")
