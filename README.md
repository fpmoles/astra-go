# astra-go
A wrapper for gocql to use DataStax Astra

[DataStax Astra](https://astra.datastax.com) is an offering of DataStax Enterprise Cassandra as a Service on major cloud 
providers. They offer various tiers of operation including a free tier that is free forever. Part of the delivery of that
offering includes MTLS certificate authentication for the drivers. While DataStax itself supports many drivers, the 
gocql driver is not one of them, so there is no specific inherent support for the secure connect bundle and MTLS 
authentication. This wrapper provides a path to unpack the Secure Connect Bundle (SCB) and pass it to gocql to get a session to work with.

## Requirements
* DataStax Astra cluster. Any tier will work.
* Download the SCB and unpack it
* Know your database username and password

## Steps to use
* <b>hostname</b> from cqlshrc in SCB
* <b>port</b> from cqlshr in SCB
* <b>filepath</b> where the `ca.crt`, `key`, and `cert` files from the SCB are stored. This is absolute file path with respect
to your application's runtime
* <b>database username</b> is the username you used when you created the Astra database
* <b>database password</b> is the password you used when you created the Astra database

### Example:

```
config := NewClusterConfig(host, port, username, password, filepath)
clusterConnection, _ := NewClusterConnection(config)
clusterConnection.Session()...
```

