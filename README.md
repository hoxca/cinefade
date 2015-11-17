Cinefade
========

A lights fading application for your homecinema.
------------------------------------------------

It's a daemon polling plex in background for a media status and control the hue room lights to switch in a cinema mode and restore previous light status when you pause or stop your movie.

Operation manual
-----------------

You talk to cinefade with curl or your browser using the http protocol.
It is an API with a limited number of verb for controlling the lights and daemon.

### Control the light bulbs


| Verb     | Action                            |
| ---      |          ---                      |
| on       | Switch on all ligths              |
| off      | Switch off all ligths             | 
| register | Save the cinema lights theme      |
| cinema   | Switch to cinema mode             |
| restore  | Restore the previous bulbs status |

exemple:

`curl http://localhost:9090/cinefade/register`

Bulbs status are saved as json file format in /var/lib/cinefade directory

* cinema.json
* current.json

It's necessary to register a cinema light theme before using the cinefade plex poller.

### Control the daemon

| Verb     | Action                            |
| ---      |          ---                      |
| start    | Start the plex poller             |
| stop     | Stop the plex poller              | 
| exit     | Stop and Exit the cinefade daemon | 

exemple:

`curl http://localhost:9090/cinefade/start`


### Configuration file

You must provide a single configuration file
/etc/cinefade/cinefade.conf

```
# config for cinefade

# Ipaddr of hue bridge
hueIpAddr = 192.168.1.99
# Hue bridge User
hueUser   = 3d99dc627158727130a0d2a224445b
# Plex URL
plexUrl   =  http://muklo:32400/status/sessions
```
