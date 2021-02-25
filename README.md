## crud

A Swagger/OpenAPI builder and validation library for servers.

Heavily inspired by [hapijs](https://hapi.dev/) and the [hapijs-swagger](https://github.com/glennjones/hapi-swagger) projects.

### Status

This project is not stable yet, API is still changing and there are missing features.

It's already pretty useful so if you risk averse then use it and pitch in!

### Why

Swagger is great, but up until now your options to use swagger are:

- Write it and then make your server match your spec.
- Write it and generate your server.
- Generate it from comments in your code.

None of these options seems like a great idea.

This project takes another approach: make a specification in Go code using nice builders where possible. The swagger is generated from this spec and validation is done before your handler gets called. 

This reduces boilerplate that you have to write and gives you nice documentation too!
