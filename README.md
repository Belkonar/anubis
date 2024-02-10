# Anubis

Anubis is an API gateway with a focus on transparent authorization.

It uses prefix paths to make it not a complete pain to host, though I may update
it to support Host header matching later.

## Database

Fairly simply, the database is just etcd. That way we can easily update the
gateway in real time.

`/service/{id}` will be the path for the configuration for a given service.

`/auth/{id}` will be the path for oauth2 configurations.

`/fga/{id}` will be the path for the fga store for a given service.

`/auth/anubis` and `/fga/anubis` will need to be defined to use the API. Once I get there I'll document what's needed
for the FGA model alongside an example one.

`/anubis/{id}` will be the path for stuff controlling anubis itself if such a thing is needed.

It will only refresh on start for now while this is in a POC mode.

## Authorization

This will use FGA/OpenFGA to handle the authorization layer. This may become
more flexible over time, but for now that'll be it.

## UI?

Not currently. API only.

## Why does this even exist?

Mostly because I want an API gateway that isn't a complete pain to use, and
that also supports real time updates of configuration. Existing solutions
don't have great ability for me to hook in the way I need to, or are "fremium"
solutions that are useless without paying a ton of money. This project will be
always be fully free (as in free beer) and open source.

## Todo?

Prior to v1.0.0 here is what needs to be done. I'm not making issues for these because that's a waste of time.

* Load configs from database.
* Add OAuth2 middleware.
  * To start with it'll only support one for everything.
* FGA
  * This as well will only support one FGA cluster.
  * I'll be writing my own client due to built-in limitations to the official one.
* Api to manage configs.

Once v1.0.0 has been released this will be after.

* Dynamic reloading
* Multiple OAuth2 configs
