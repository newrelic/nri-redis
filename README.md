# New Relic Infrastructure Integration for Redis
In order to know how the Redis integration works and how to run it with the Infrastructure agent please check [the documentation](https://docs.newrelic.com/docs/redis-integration-new-relic-infrastructure).

## Integration Development usage

Assuming that you have the source code and Go tool installed you can build and run the Redis Integration locally.
* Go to the directory of the Redis integration and build it
```bash
$ make
```
* The command above will execute the tests for the Redis integration and build an executable file called `nr-redis` under `bin` directory. Run `nr-redis`:
```bash
$ ./bin/nr-redis
```
* If you want to know more about usage of `./bin/nr-redis` check
```bash
$ ./bin/nr-redis -help
```

For managing external dependencies [govendor tool](https://github.com/kardianos/govendor) is used. It is required to lock all external dependencies to specific version (if possible) into vendor directory.
