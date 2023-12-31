# Golang simple Rate Limiter

This is a simple rate limiter library that I made in a single day for educational purposes, so it's not recommended to use it on production environments.

## Instructions

- Copy the file `.env.sample` and paste it on your project's root directory, and then name it `.env`.
- Configure the environment variables according to your preferences (see the [config/config.go](config/config.go) file for details, or the [Environment Variables table](#environment-variables) below).

## Running the example web server to test the library

The easiest way to test the rate limiter with different configurations is by using Docker Compose and the [example web server](cmd/example_web_server.go)

- Clone this git repository into your workspace
- Follow the same [instructions in the previous section](#instructions) to configure your environment
- Run `docker compose up` in the project's root directory
- If you have VS Code installed and have the `REST Client` extension enabled, you can use the [api.http](api.http) file to send requests, or you can use any other REST client like `Postman`, or any other tool (e.g.: Apache `ab` CLI).

## Environment Variables

|Name|Accepts|Default Value|Description|
|----|-------|-------------|-----------|
|RATE_LIMITER_IP_ADDRESS_MAX_REQUESTS|number|2|Max requests per IP address|
|RATE_LIMITER_IP_ADDRESS_LIMIT_IN_SECONDS|number|1|IP Address limit duration in seconds (the amount of time the max requests are allowed in)|
|RATE_LIMITER_IP_ADDRESS_BLOCK_IN_SECONDS|number|5|IP Address block duration in seconds (the amount of time the IP address is blocked for after exceeding the max requests)|
|RATE_LIMITER_TOKENS_HEADER_KEY|string|API_KEY|The requests' Header key to use for the tokens|
|RATE_LIMITER_TOKENS_CONFIG_TUPLE|string||A list of tokens separated by a comma and their respective max requests, limit and block durations in seconds separated by a colon|
|RATE_LIMITER_STORE_STRATEGY|string (must be one of `in_memory` or `redis`)|in_memory|The strategy to use for the store|
|RATE_LIMITER_REDIS_HOST|string|localhost|Redis host|
|RATE_LIMITER_REDIS_PORT|number|6379|Redis port|
|RATE_LIMITER_REDIS_PASSWORD|string||Redis password|
|RATE_LIMITER_REDIS_DB|number|0|Redis DB|
