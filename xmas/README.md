# XMAS special

Every year many of the users of the UK open rail data provides an extra service for Christmas Eve.

[Real Time Trains](https://www.realtimetrains.co.uk/) for example have their 1X01 service from North Pole International [XNP]
calling at every station on the UK mainland.

I've done something similar in the past but didn't last year so for 2019 I'm resurrecting it & this is the project to support it.

This directory contains a tool that:
* Generates a custom schedule based on the station's geographical location and a specific time
* Play that schedule as a pseudo live service into the system so that it actually appears to run in real time

## Custom reference data
Now first the Origin & Terminating station, [North Pole](https://en.wikipedia.org/wiki/North_Pole_depot) actually does exist on the network,
with a CRS code of XNP & Tiploc NPLEINT.

It's a maintenance centre in the London Borough of Hammersmith & Fulham and between 1994 and 2007 it was known as North Pole International when it was used by Eurostar when they terminated services at London Waterloo.

The ref microservice handles this by ensuring that the "station" name is "North Pole International" - in the reference feed the name is identical to the tiploc so we have to do this so we show a more realistic name.

We also create a dummy TOC (to provide a suitable service name) with code XM. This doesn't exist in the reference feed but makes the result better.
