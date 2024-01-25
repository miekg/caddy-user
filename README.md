# caddy-user

This caddy module performs a setuid on the goroutine handling the request. This works.

HOWEVER DUE TO THE UNIX PROCESS MODULE IT CAN'T WORK

Setuid works on the entire process, not a single goroutine. So while this does what is advertized,
it can't work for concurrent requests or even setuid-ing to different user accounts.

Take this example:

* caddy runs as 'root'
* request comes in, setuid to 'x', caddy now runs as 'x'
* another request comes in, setuid to 'y' fails as user 'x' is not allowed to do that
* last request will run under the user 'x'
* request for x is completed, caddy reverts back to 'root'

So this will *sometimes* do what you expect.

A nicer idea might be to start Caddy, fork into multiple caddys and somehow solve it there.

Another alternative is running a proxy in front of caddys running as different users (and
potentially different ports).
