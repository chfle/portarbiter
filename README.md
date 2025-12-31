# portarbiter

**portarbiter** is a Linux system administration tool that determines **who actually owns a TCP listening port**
and can terminate the correct owner **safely**.

It is designed to avoid the classic mistakes of naive port-killing tools that blindly kill
systemd services, Docker, or SSH.

---

## What does it do?

When a port is in use, ownership is often unclear. The port may belong to:

- a user process
- a systemd service
- a Docker container
- a docker-compose project
- an SSH user session (port forwarding)

**portarbiter resolves this ownership correctly and hierarchically.**

---

## Resolution order

1. Docker containers / docker-compose projects  
2. systemd services (real daemons only)  
3. Raw processes  

Correct behavior examples:

- SSH port forwarding is treated as a user process, not `ssh.service`
- Docker published ports are attributed to the container or compose project
- Supervisor services are never killed blindly

---

## Usage

Inspect port ownership (safe default):

```bash
portarbiter 5432
```

Terminate the owner:

```bash
portarbiter --kill 5432
```

Force termination:

```bash
portarbiter --kill --force 5432
```

Non-interactive usage:

```bash
portarbiter --kill --yes 5432
```

Show version:

```bash
portarbiter --version
```

---

## Installation

Debian / Ubuntu:

```bash
sudo dpkg -i portarbiter_<version>_amd64.deb
```

Manpage:

```bash
man portarbiter
```

---

## Build

Requirements:
- Linux
- Go 1.22+

```bash
make build
make deb
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## Author

Christian Lehnert

---

## License

MIT

