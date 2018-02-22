# fastcep
> Fast microservice for cep consulting written in go

Fastcep is a proof of concept for a fast api for consulting brazil postal codes (CEPs). It is built with go and have an aggressive cache using redis. The main go is to build a fast, reliable alternative to services like [viacep](https://viacep.com.br/) e [cepaberto](http://cepaberto.com/). The cep database used is a extract from the public data provided by cepaberto. Feel free to suggest improvements.

## Stack

- Go
- PostgreSQL
- Redis

## License

[Apache License, Version 2.0](LICENSE) &copy; [Thiago Santos](https://github.com/thiamsantos)
