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

## Contributing Code

We welcome code contributions (in the form of pull requests) from our user
community. Before submitting a pull request please review [these guidelines](https://github.com/newrelic/nri-redis/blob/master/CONTRIBUTING.md).

Following these helps us efficiently review and incorporate your contribution
and avoid breaking your code with future changes to the agent.

## Custom Integrations

To extend your monitoring solution with custom metrics, we offer the Integrations
Golang SDK which can be found on [github](https://github.com/newrelic/infra-integrations-sdk).

Refer to [our docs site](https://docs.newrelic.com/docs/infrastructure/integrations-sdk/get-started/intro-infrastructure-integrations-sdk)
to get help on how to build your custom integrations.

## Support

You can find more detailed documentation [on our website](http://newrelic.com/docs),
and specifically in the [Infrastructure category](https://docs.newrelic.com/docs/infrastructure).

If you can't find what you're looking for there, reach out to us on our [support
site](http://support.newrelic.com/) or our [community forum](http://forum.newrelic.com)
and we'll be happy to help you.

Find a bug? Contact us via [support.newrelic.com](http://support.newrelic.com/),
or email support@newrelic.com.

New Relic, Inc.
