# darwinkb

Originally this was a standalone microservice like ref & tt - however as this data is mostly static the new model is to
use a central PostgreSQL database, moving the public rest API outside into the same one used by darwindb.

This service is now used simply to populate that database, retrieving the XML regularly from NRE/RDG as well as the
processing of the real time incidents feed fed in via rabbitMQ.

 