# Crypt

Util to compute and verify Bcrypted passwords via CLI.

## Installation

```
go get github.com/sbreitf1/go-crypt
go install github.com/sbreitf1/go-crypt
```

## Usage

### Crypt

Crypt passwords by supplying them via input parameter, StdIn or user prompt:

```
go-crypt -i 'password'

echo -n 'password' | go-crypt

go-crypt -p
> {hidden password}
```

You can specify the crypt method using `-m`:

```
go-crypt -i 'password' -m 'sha512'
```

See the list of available crypt methods below. Bcrypt versions are chosen to be consistent with the definitions in [Wikipedia](https://en.wikipedia.org/wiki/Bcrypt).

| Key | Short | Example Output Prefix | Description |
| --- | ----- | --------------------- | ----------- |
| _default_ | 2a | `$2a$10$K4NBx0JItxC5D0NPRBcz2u...` | Standard bcrypt with UTF-8 support |
| sha256 | 5 | `$5$yx5Dcz8vmsvV$...` | SHA-256 based crypt |
| sha512 | 6 | `$6$FD3.owZTk1v2$...` | SHA-512 based crypt |

### Verify

The same password input options as for crypt are available for password verification. Use the option `-v` to pass a crypted value:

```
go-crypt -i 'password' -v '$2a$10$p.qy7f4OptuEORh/lSCNC.3U39ra0F5FRPXfxoM8IG0JQvZPUZqbq'
password is ok

go-crypt -i 'wrong password' -v '$2a$10$p.qy7f4OptuEORh/lSCNC.3U39ra0F5FRPXfxoM8IG0JQvZPUZqbq'
the supplied password does not match the given crypt value
```

The application will exit with code 1 if the password is incorrect or an error occured.