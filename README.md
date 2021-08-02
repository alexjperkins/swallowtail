# Swallowtail Project

Algorithmic trading & Notification platform.

WIP; use at own discression.


### Testing gRPC

Run the correct services in separate terminals; you may need to rebuild the docker file again.

To list services on a given host:

```shell
	grpcurl -plaintext <host:port> <rpc-service>
```

To call an RPC function:

```shell
	grpcurl -plaintext -d '{...}' <host:port> <rpc-service>.<rpc-function>
```
