# Battery Energy Storage System (BESS) EMS

A production-ready Energy Management System (EMS) for Battery Energy Storage Systems built with Go, following standard Go project layout and best practices.

## 🏗️ Architecture

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```text
solservices-gyongyoshalasz-ems/
├── cmd/                    # Main applications
├── internal/               # Private application code
├── pkg/                    # Public library code
├── api/                    # API definitions
├── configs/                # Configuration files
├── build/                  # Docker and CI files
├── scripts/                # Build and deployment scripts
├── docs/                   # Documentation
├── test/                   # Integration tests
└── examples/               # Example applications
```

## 🚀 Features

### Core Functionality

- **MODBUS TCP Communication**: Full JINKO SCU protocol implementation
- **Real-time Monitoring**: BMS and PCS data collection
- **Automatic Control**: SOC-based charge/discharge management
- **Comprehensive Alarms**: 80+ alarm types with severity levels
- **REST API**: Complete HTTP API for integration
- **Data Persistence**: SQLite with automatic cleanup
- **Metrics Collection**: System performance monitoring

### Safety & Protection

- Multi-level SOC protection (configurable limits)
- Temperature-based power derating
- Fault state protection and recovery
- Connection monitoring with auto-reconnection
- Manual override capabilities

## 📦 Installation

### Prerequisites

- Go 1.21+
- SQLite3 (for database)
- Network access to BMS/PCS devices

### Quick Start

```bash
# Clone the repository
git clone https://powerkonnekt/ems
cd solservices-gyongyoshalasz-ems

# Build and run
make build
./ems

# Or use Docker
make docker
make docker-run
```

### Development Setup

```bash
# Install development dependencies  
make dev-deps

# Run with live reload
make dev

# Run tests
make test

# Run integration tests
make test-integration
```

## ⚙️ Configuration

Update `configs/config.json` with your system settings:

```json
{
    "bms": {
        "host": "192.168.1.100",
        "port": 502,
        "slave_id": 1,
        "rack_count": 10
    },
    "pcs": {
        "host": "192.168.1.101", 
        "port": 502,
        "slave_id": 1
    },
    "ems": {
        "max_soc": 95.0,
        "min_soc": 5.0,
        "max_charge_power": 100.0,
        "max_discharge_power": 100.0
    }
}
```

## 🔧 API Usage

### System Status

```bash
curl http://localhost:8080/api/v1/status
```

### BMS Data

```bash  
curl http://localhost:8080/api/v1/bms/data
```

### Control Mode

```bash
curl -X POST http://localhost:8080/api/v1/control/mode \
  -H "Content-Type: application/json" \
  -d '{"mode": "MANUAL"}'
```

### Manual Power Command

```bash
curl -X POST http://localhost:8080/api/v1/control/power \
  -H "Content-Type: application/json" \
  -d '{"power": 25.5}'
```

See [API Documentation](docs/API.md) for complete endpoint reference.

## 🐳 Docker Deployment

```bash
# Build and run with docker-compose
cd build
docker-compose up -d

# Or build manually
docker build -t solservices-gyongyoshalasz-ems -f build/Dockerfile .
docker run -d -p 8080:8080 -v ./configs:/app/configs solservices-gyongyoshalasz-ems
```

## 🔧 Development

### Project Structure

- **`cmd/ems/`** - Main application entry point
- **`internal/`** - Private business logic
  - `config/` - Configuration management
  - `database/` - Data models and persistence  
  - `bms/` - BMS integration and protocol
  - `pcs/` - PCS integration
  - `alarm/` - Alarm processing
  - `control/` - Control algorithms
  - `api/` - HTTP API handlers
  - `ems/` - Main orchestrator
  - `metrics/` - Performance monitoring
- **`pkg/`** - Reusable packages
  - `modbus/` - MODBUS TCP client
  - `utils/` - Utility functions

### Adding New Features

1. **New Device Support**: Add to `internal/` following existing patterns
2. **API Endpoints**: Extend `internal/api/handlers.go` and `routes.go`  
3. **Control Logic**: Modify `internal/control/logic.go`
4. **Configuration**: Update `internal/config/config.go`

### Testing

```bash
# Unit tests
make test

# Integration tests (requires test environment)
make test-integration  

# Benchmarks
make benchmark

# Coverage report
make test
go tool cover -html=coverage.out
```

## 📊 Monitoring

### Built-in Metrics

- System performance (CPU, memory, disk)
- Communication statistics  
- Alarm counts and history
- Energy throughput tracking

### Health Checks

```bash
# Manual health check
curl http://localhost:8080/health

# Automated monitoring
./scripts/monitor.sh
```

### Log Analysis

```bash
# Real-time logs
make logs

# Search for errors
grep "ERROR" logs/ems.log

# Connection monitoring
grep "connection" logs/ems.log
```

## 🚀 Production Deployment

### Service Installation

```bash
# Install as systemd service
sudo make service-install

# Control service
sudo systemctl start solservices-gyongyoshalasz-ems
sudo systemctl status solservices-gyongyoshalasz-ems
sudo systemctl logs -f solservices-gyongyoshalasz-ems
```

### Backup & Maintenance

```bash
# Create backup
make backup

# Update to new version
make update VERSION=v1.2.0

# Cleanup old data
./scripts/cleanup.sh
```

## 🔒 Security

- Network isolation for BMS/PCS communication
- Input validation on all API endpoints
- Database access controls
- Audit logging for control commands
- Regular security updates

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Follow Go conventions and project structure
4. Add tests for new functionality
5. Update documentation
6. Submit a pull request

### Code Standards

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `make fmt` for formatting
- Run `make lint` before committing
- Maintain test coverage >80%

## 📝 License

[MIT License](LICENSE)

## 🆘 Support

- **Documentation**: [docs/](docs/)
- **API Reference**: [docs/API.md](docs/API.md)
- **Deployment Guide**: [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)
- **Issues**: [GitHub Issues](https://powerkonnekt/ems/issues)
