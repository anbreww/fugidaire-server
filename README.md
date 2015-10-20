# Le Fûgidaire v2.0

The Fugidaire v2.0 is a 6-tap kegerator for home brewed beer.
It uses a three-component system to monitor and display the list of beers that
are currently on tap.

## Fugidaire-server

The main component is the fugidaire server. It allows us to track and update the
beer list remotely. It can serve the taplist through a simple REST API in JSON
or display the tap list on a web interface.

It will be able to subscribe to MQTT notifications of pours from the tap monitors
to update the level of beer remaining in the keg.

## Fugidaire-flow

The level of beer in each keg is monitored by a flow sensor. The fugidaire-flow
device reads the flow of beer and sends periodic updates over MQTT to the server
so that it can update the level of beer in the keg. In turn, we can then update
the level display on the graphical interface.

## Fugidaire-gui

Each tap is provided with a small colour LCD to show the current beer stats,
including colour, style, abv, ibu and fill level of the keg.
