# Plocate UI - Fast File Search for Unraid

A modern, dockerized file search application powered by `plocate` - the fastest file location utility. Perfect for Unraid servers with large media collections, documents, and backups.

## Features

- ⚡ **Lightning Fast**: Uses plocate for instant file searches
- 🔄 **Automatic Indexing**: Schedule automatic reindexing at configurable intervals
- 🎨 **Modern Web UI**: Clean, responsive Svelte-based interface
- 🐳 **Docker Native**: Optimized for Unraid with configurable volume mounts
- 💾 **Configurable Storage**: Place database on cache drive for optimal performance
- 📊 **Status Dashboard**: Monitor indexing status and schedule
- 🔍 **Smart Search**: Case-insensitive, partial matching with instant results

## Quick Start

### Prerequisites

- Unraid server (or any Docker-compatible host)
- Access to SSH or Unraid terminal

### Installation on Unraid

#### Method 1: Docker Compose (Recommended)

1. **Clone or download this repository** to your Unraid server:
   ```bash
   cd /mnt/user/appdata
   git clone https://github.com/yourusername/plocate-ui.git
   cd plocate-ui
   ```

2. **Edit the configuration** file:
   ```bash
   cp config.example.yml config.yml
   nano config.yml
   ```

   Configure paths to index:
   ```yaml
   plocate:
     database_path: "/var/lib/plocate/plocate.db"
     index_paths:
       - "/mnt/user/media"
       - "/mnt/user/documents"
       - "/mnt/cache"
   ```

3. **Edit docker-compose.yml** to configure volume mounts:
   ```yaml
   volumes:
     # Database on cache drive (recommended for performance)
     - /mnt/cache/appdata/plocate-ui/db:/var/lib/plocate

     # Mount paths you want to search (read-only)
     - /mnt/user:/mnt/user:ro
     - /mnt/cache:/mnt/cache:ro
   ```

4. **Build and run**:
   ```bash
   docker-compose up -d
   ```

5. **Access the UI**:
   Open your browser to `http://YOUR-UNRAID-IP:8080`

#### Method 2: Docker CLI

```bash
docker run -d \
  --name plocate-ui \
  -p 8080:8080 \
  -v /mnt/cache/appdata/plocate-ui/config.yml:/app/config.yml:ro \
  -v /mnt/cache/appdata/plocate-ui/db:/var/lib/plocate \
  -v /mnt/user:/mnt/user:ro \
  -v /mnt/cache:/mnt/cache:ro \
  -e TZ=America/New_York \
  --restart unless-stopped \
  plocate-ui
```

#### Method 3: Unraid Community Applications (Future)

Once published to CA, simply search for "Plocate UI" in Community Applications and click Install.

## Configuration

### config.yml Structure

```yaml
server:
  port: "8080"  # Web UI port

plocate:
  # Database location (recommend cache drive for performance)
  database_path: "/var/lib/plocate/plocate.db"

  # Paths to index
  index_paths:
    - "/mnt/user"           # All user shares
    - "/mnt/cache"          # Cache drive
    # - "/mnt/disk1"        # Specific disk
    # - "/mnt/user/media"   # Specific share

  # Binary paths (usually don't need to change)
  updatedb_bin: "updatedb"
  plocate_bin: "plocate"

scheduler:
  enabled: true

  # Cron-style schedule
  interval: "0 */6 * * *"  # Every 6 hours

  # Other examples:
  # "0 2 * * *"     - Daily at 2 AM
  # "0 3 * * 0"     - Weekly on Sunday at 3 AM
  # "0 4 1 * *"     - Monthly on the 1st at 4 AM
```

### Environment Variables (Optional)

Override config.yml settings with environment variables:

- `PORT` - Web server port (default: 8080)
- `PLOCATE_DB_PATH` - Database file path
- `INDEX_INTERVAL` - Cron schedule string
- `CONFIG_PATH` - Path to config file

## Usage

### Web Interface

1. **Search Files**:
   - Enter filename or pattern in the search box
   - Press Enter or click Search
   - Results appear instantly with full paths

2. **Monitor Status**:
   - View last indexed time
   - See next scheduled indexing
   - Check which paths are indexed

3. **Control Indexing**:
   - **Start Index Now**: Trigger immediate reindex
   - **Stop Indexing**: Cancel running index operation
   - **Enable/Disable Scheduler**: Control automatic indexing

### API Endpoints

The application also exposes a REST API:

- `GET /api/status` - Get current status
- `GET /api/search?q=filename&limit=100` - Search files
- `POST /api/control/start` - Start indexing
- `POST /api/control/stop` - Stop indexing
- `POST /api/control/scheduler/enable` - Enable scheduler
- `POST /api/control/scheduler/disable` - Disable scheduler

Example API usage:
```bash
# Search for files
curl "http://localhost:8080/api/search?q=movie.mkv"

# Get status
curl "http://localhost:8080/api/status"

# Trigger manual index
curl -X POST "http://localhost:8080/api/control/start"
```

## Performance Tips

### For Best Performance on Unraid:

1. **Place database on cache drive**:
   ```yaml
   volumes:
     - /mnt/cache/appdata/plocate-ui/db:/var/lib/plocate
   ```

2. **Use SSD cache**: If available, ensures fast index operations

3. **Schedule indexing during low usage**:
   ```yaml
   interval: "0 3 * * *"  # 3 AM daily
   ```

4. **Index specific shares** instead of entire `/mnt/user`:
   ```yaml
   index_paths:
     - "/mnt/user/media"
     - "/mnt/user/documents"
   ```

5. **Limit indexed paths**: Only index what you need to search

### Typical Index Sizes:

- 100,000 files: ~10-20 MB database
- 1,000,000 files: ~100-200 MB database
- 10,000,000 files: ~1-2 GB database

### Index Time:

- 100,000 files: ~30 seconds
- 1,000,000 files: ~5 minutes
- 10,000,000 files: ~30-60 minutes

*Times vary based on disk speed and file distribution*

## Development

### Building from Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/plocate-ui.git
   cd plocate-ui
   ```

2. **Build with Docker**:
   ```bash
   docker build -t plocate-ui .
   ```

3. **Or build manually**:

   Frontend:
   ```bash
   cd frontend
   npm install
   npm run build
   ```

   Backend:
   ```bash
   cd backend
   go mod download
   go build -o plocate-ui
   ```

### Development Mode

Run frontend and backend separately for development:

1. **Backend**:
   ```bash
   cd backend
   go run main.go
   ```

2. **Frontend** (in another terminal):
   ```bash
   cd frontend
   npm run dev
   ```

Access development UI at `http://localhost:5173`

## Architecture

```
┌─────────────────────────────────────────┐
│         Svelte Frontend (Port 8080)     │
│  ┌──────────┐ ┌────────┐ ┌──────────┐  │
│  │  Search  │ │ Status │ │ Controls │  │
│  └──────────┘ └────────┘ └──────────┘  │
└─────────────┬───────────────────────────┘
              │ REST API
┌─────────────▼───────────────────────────┐
│          Go Backend (Gin)               │
│  ┌──────────────┐  ┌─────────────────┐ │
│  │   Handlers   │  │  Cron Scheduler │ │
│  └──────┬───────┘  └────────┬────────┘ │
│         │                   │          │
│  ┌──────▼───────────────────▼────────┐ │
│  │        Indexer Manager            │ │
│  └──────┬────────────────────┬───────┘ │
└─────────┼────────────────────┼─────────┘
          │                    │
┌─────────▼────────┐  ┌────────▼─────────┐
│  updatedb binary │  │  plocate binary  │
│  (indexing)      │  │  (searching)     │
└─────────┬────────┘  └────────┬─────────┘
          │                    │
          ▼                    ▼
    ┌──────────────────────────────┐
    │   plocate.db (on cache)      │
    └──────────────────────────────┘
```

## Technology Stack

- **Backend**: Go (Golang) with Gin web framework
- **Frontend**: Svelte + Vite + Tailwind CSS
- **Search Engine**: plocate (system binary)
- **Scheduler**: robfig/cron
- **Container**: Docker (Ubuntu base image)

## Troubleshooting

### Issue: "No results found" when files exist

**Solution**:
1. Check if indexing has completed
2. Verify paths are correctly mounted in docker-compose.yml
3. Trigger manual reindex from the UI

### Issue: Indexing fails or stops

**Solution**:
1. Check container logs: `docker logs plocate-ui`
2. Verify database directory permissions
3. Ensure sufficient disk space on cache drive

### Issue: Container won't start

**Solution**:
1. Check logs: `docker logs plocate-ui`
2. Verify config.yml syntax
3. Ensure port 8080 is not already in use
4. Check volume mount paths exist

### Issue: Slow search performance

**Solution**:
1. Move database to cache/SSD drive
2. Reduce number of indexed paths
3. Check disk I/O with `iotop`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [plocate](https://plocate.sesse.net/) - The incredibly fast locate implementation
- [Gin](https://github.com/gin-gonic/gin) - Go web framework
- [Svelte](https://svelte.dev/) - Frontend framework
- Unraid community for inspiration and support

## Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/plocate-ui/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/plocate-ui/discussions)
- **Unraid Forums**: [Community Support Thread](#)

---

Made with ❤️ for the Unraid community
