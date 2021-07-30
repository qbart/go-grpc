# Ports

# Development

Generate protobuffers:
```
make protoc        # generate
```

Start the server:
```
make server
```

Start the client:
```
make client
```

# Dependencies

Install protoc from the official website then install the tools:

```
make install_proto # install if missing
```

# Running the example

Alternatively you can start everything in docker:

```bash
docker-compose build
docker-compose up
```

test the endpoint:
```
make get
```

or mising one:
```
make get_invalid
```

