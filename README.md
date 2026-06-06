# VisiMon

> Platform monitoring infrastruktur real-time yang ringan, mudah dipasang, dan kaya visualisasi.

VisiMon is a lightweight, real-time infrastructure monitoring platform designed for DevOps teams and developers in Indonesia and Southeast Asia. Think Netdata, but with a stronger focus on ease of setup.

## Project Structure

```
.
├── agent/          # Go agent — system metrics collector
│   ├── cmd/        # Entry point
│   └── pkg/        # Collector, reporter, config packages
├── server/         # API server (coming soon)
├── dashboard/      # React frontend (coming soon)
├── docker/         # Dockerfiles
└── scripts/        # Install scripts
```

## Quick Start

### Running the Agent

```bash
cd agent
go build -o visimon-agent ./cmd/visimon-agent
./visimon-agent --server http://localhost:8080
```

### With a config file

```yaml
# visimon.yaml
server_url: "http://localhost:8080"
interval: 1
hostname: "my-server-01"
```

```bash
./visimon-agent --config visimon.yaml
```

## License

MIT