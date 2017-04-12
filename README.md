![https://img.shields.io/badge/Status-Feature%20Freeze-blue.svg](https://img.shields.io/badge/Status-Code%20Freeze-blue.svg)

Currently in feature freeze to polish for v3.0.

# Platypus

[Live Master Branch](https://status.ggserv.xyz)

Scaleable Server Infrastructure Monitoring Python App

(Probably Light and Technically Pretty Unassuming Software)

Even More Buzzwords For SEO Reasons

Monitors and reports statistics of your server infrastucture, including usage statistics, uptime, downtime, etc.

### Current Features
 - Setup script
 - Websocket-based uptime monitoring
 - Auto-post to Slack when server goes offline
 - Provide live server usage statistics with web frontend
 - Simple JSON API for building apps

### Planned Features
 - Cleaner code
 - Smaller resource footprint
 - More advanced scanning method

## Requirements
 - Python 3.x (2.x not officially supported)
 - pip
  - `pip install -r requirements.txt`
 - MariaDB

## Running
Run `python setup.py` to set up your instance, including creating the MariaDB databse and setting an admin password.

Finally you can run `python App.py` and go navigate to `127.0.0.1:8080/login` to get to the admin control panel.

To expose it to the world, I recommend using an [nginx proxy](https://www.nginx.com/resources/admin-guide/reverse-proxy/).
## Configuration

| Key | Default Value | What is |
| --- | ------------- | ------- |
| enable_slackbot | false | Enable automated Slack reports |
| enable_webserver | true | Enable web frontend and API | 
| webserver_port | 8080 | Port for webserver to run on |
| slack_api_key | | [API key for Slack posting](https://github.com/slackapi/python-slackclient)
| admin_username | admin | Username for admin interface |
| admin_password | | Password for admin interface |
| sql_user | root | SQL database username |
| sql_pass | | SQL database password |
| sql_host | localhost | SQL database host |
| sql_db | server | SQL data table |

## Script

Included in the `Scripts` folder is the node-side `aor.py` script. This should be
deployed to the machines you would like to monitor, and modified to point towards
your master server. A docker image is available [here](#) if you would prefer to 
use that. You can also feel free to use your own custom script, however it has to
be fully websocket compatible. See below about the data passed to the master.

You'll want to edit `config.json` and deploy it alongside aor.py, otherwise
it will use default values, and probably will not connect to your master. 

### Requirements:

 - Python 3.x
 - websocket
 - psutil

`pip install websocket psutil`

### Deploying (Script)

If you would like to tinker with the script, take a peek at `Scripts/generate_aor_packages.py`
and modify the values to match up to your database. It will automagically create
zip packages containing the `aor.py` and `aor_config.json` files (already filled
out too).

Alternatively, you'll want to edit `Scripts/config.json` to match the values in your database,
namely changing the `"UUID":` field. You can leave the interval alone. 

Then you can copy the config and `aor.py` script to your server and run the script, 
preferably in a `screen` session. For the most part you can leave the script alone
at this point, you do not need to set up a reverse proxy or the likes.

### Deploying (Docker)

Coming soon