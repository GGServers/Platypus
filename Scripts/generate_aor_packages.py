import json
import MySQLdb
import zipfile
import os
import shutil

# Change this
conn = MySQLdb.connect(user="root", db="platypus")
c = conn.cursor()

c.execute("SELECT * FROM server")
data = c.fetchall()

shutil.copy("aor.py", "aor_output/tmp/")
shutil.copy("aor_config.json", "aor_output/tmp/")
os.chdir("aor_output/tmp/")

for server in data:

    with open("aor_config.json") as data_file:
        config = json.load(data_file)

    # Change this to match your database if you get errors
    config["uuid"] = server[10]
    # print(config)

    with open("aor_config.json", "w") as data_file:
        json.dump(config, data_file, indent=4)

    zf = zipfile.ZipFile("%s.zip" % ("../" + server[2]), "w")

    for dirname, subdirs, files in os.walk("."):
        for filename in files:
            zf.write(filename)

    zf.close()

    print("Zipped aor script + config for", server[1])
