![Build Status](http://jenkins.area51.onl/buildStatus/icon?job=Public/Badges) ![Issues](https://badge.area51.onl/github/issues/peter-mount/nre-feeds.svg) ![Last commit](https://badge.area51.onl/github/last-commit/peter-mount/nre-feeds.svg)
# Darwin
go library &amp; suite of microservices for handling the NRE DarwinD3 feeds

The main purpose of this project is to consume the feeds provided by National Rail Enquiries in real time and expose that information as a REST service which can be consumed by a client, usually a website.

https://departureboards.mobi/ is an example of one of these clients.

The documentation is in the [Wiki](https://github.com/peter-mount/nre-feeds/wiki)

## V12 vs V16 pushport feed

Versions up to 0.4 (including the matching branches) are for the Darwin v12 Pushport feed. As of Nov 26 2018 the master branch is based on the v16 feed.

Until the v16 feed is live, do not use the master branch
