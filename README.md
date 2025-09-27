# 🐋 Whalio

Modern web application framework built with Go, templ, HTMX, TailwindCSS, and DaisyUI.

## ✨ Features

- **🐹 Go Backend**: Fast, reliable server with Chi router
- **📝 Templ Templates**: Type-safe HTML templates
- **⚡ HTMX**: Dynamic interactivity without complex JavaScript
- **🎨 TailwindCSS**: Utility-first CSS framework  
- **🌸 DaisyUI**: Beautiful semantic components
- **🔧 Hot Reloading**: Automatic rebuild and reload during development
- **🛡️ Security**: Built-in security headers and CORS
- **📊 Logging**: Structured logging with zerolog
- **🚀 Production Ready**: Optimized builds and Docker support

## 🏗️ Project Structure

```
whalio/
├── cmd/                    # Application entry point
│   └── main.go
├── config/                 # Configuration management
│   └── config.go
├── handlers/               # HTTP handlers
│   └── handlers.go
├── templates/              # Templ templates
│   ├── layout.templ
│   └── index.templ
├── assets/                 # Source assets
│   └── css/
│       └── input.css
├── static/                 # Built static files
│   ├── css/
│   └── js/
├── package.json            # Node.js dependencies
├── tailwind.config.js      # TailwindCSS configuration
├── Makefile               # Development commands
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- **Go** 1.25+ ([install](https://golang.org/doc/install))
- **Bun** 18+ ([install](https://bun.sh/))
- **Templ** CLI ([install](https://templ.guide/quick-start/installation))

### Setup

1. **Clone or setup project directory**:
   ```bash
   whalio
   ```

2. **Install dependencies**:
   ```bash
   make setup
   ```

3. **Start development server**:
   ```bash
   make dev
   ```

4. **Open your browser**:
   ```
   http://localhost:8080
   ```

## 🛠️ Development

### Available Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make setup` | Setup development environment |
| `make build` | Full production build |
| `make test` | Run tests |
| `make fmt` | Format code |
| `make clean` | Clean build artifacts |

### Development Workflow

1. **Templates**: Edit `.templ` files in `templates/`
   - Auto-generates Go code
   - Hot reloads on changes

2. **Styles**: Edit `assets/css/input.css`
   - TailwindCSS + DaisyUI components
   - Auto-rebuilds on changes

3. **Backend**: Edit Go files in `templates/`, etc.
   - Auto-restarts server on changes

4. **Static Assets**: Place in `static/` directory
   - Served directly by the server

### Hot Reloading

The development server (`air run`) runs multiple processes:
- **Templates**: Regenerate go files
- **Go Server**: Restarts on backend changes

## 🔧 Configuration

Configuration is handled through environment variables and command-line flags:

### Environment Variables

```bash
# Server settings
export PORT=8080
export HOST=localhost
export ENVIRONMENT=development

# Logging
export LOG_LEVEL=info
export LOG_FORMAT=console

# Features
export DEBUG=true
export RATE_LIMIT_ENABLED=true
```

### Command Line Flags

```bash
go run cmd/main.go -port 3000 -debug -env production
```

## 🧪 Testing

### Run Tests

```bash
# Basic tests
make test

# With race detection
make test-race

# With coverage
make test-cover
```

## 🚀 Production Deployment

### Build for Production

```bash
# Build optimized binary
make build-prod

# The binary will be in bin/whalio-linux-amd64
```

### Environment Setup

```bash
# Production environment variables
export ENVIRONMENT=production
export LOG_FORMAT=json
export LOG_LEVEL=warn
export DEBUG=false
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Format code: `make fmt`
5. Run tests: `make test`
6. Commit changes: `git commit -m 'Add amazing feature'`
7. Push to branch: `git push origin feature/amazing-feature`
8. Submit a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with ❤️ using modern web technologies
