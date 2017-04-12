import MySQLdb
from src.Config import Config
import bcrypt

config = Config()

if __name__ == "__main__":
    print("Platypus setup utility")
    print("SQL Database Setup")
    sql_host = input("Host: ")
    sql_user = input("User: ")
    sql_password = input("Password: ")
    sql_database = input("Table: ")

    con = MySQLdb.connect(host=sql_host, user=sql_user, passwd=sql_password)
    c = con.cursor()

    c.execute("CREATE DATABASE IF NOT EXISTS platypus")
    c.execute("USE platypus")
    c.execute("CREATE TABLE IF NOT EXISTS " + sql_database + " (id int NOT NULL AUTO_INCREMENT," +
              "name varchar(120), hostname varchar(120)," +
              "online boolean DEFAULT true, ip varchar(50), uuid varchar(50)" +
              " PRIMARY KEY(id))")
    con.commit()

    print("SQL Databse created, moving along.")
    print("Admin user creation")
    admin_user = input("Username: ")
    admin_password = input("Password: ")
    admin_password2 = input("Confirm Password: ")

    if admin_password == admin_password2:
        admin_password2 = admin_password2.encode("utf8")
        hashed = bcrypt.hashpw(admin_password2, bcrypt.gensalt())
    else:
        print("Passwords did not match. Your password:", admin_password2)

    print("Slackbot defaulting to off")
    print("Webserver defaulting to on")
    print("Setting config")
    # Set SQL Config
    config.Set("sql_user", sql_user)
    config.Set("sql_pass", sql_password)
    config.Set("sql_host", sql_host)
    config.Set("sql_db", sql_database)
    # Set Aadmin Config
    config.Set("admin_username", admin_user)
    config.Set("admin_password", str(hashed).strip("b'"))
    # Finished
    print("Further configuration can be found in config.json")
