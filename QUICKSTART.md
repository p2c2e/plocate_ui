# Quick Start Guide

Get Plocate UI running on your Unraid server in 5 minutes!

## Prerequisites

- Unraid server with Docker support
- SSH access to your Unraid server

## Installation Steps

### 1. Download

SSH into your Unraid server and run:

```bash
cd /mnt/user/appdata
git clone https://github.com/yourusername/plocate-ui.git
cd plocate-ui
```

Or download and extract the ZIP file to `/mnt/user/appdata/plocate-ui`

### 2. Configure

Edit `config.yml` to specify which paths to index:

```bash
nano config.yml
```

Minimum configuration:
```yaml
plocate:
  database_path: "/var/lib/plocate/plocate.db"
  index_paths:
    - "/mnt/user"  # Change this to your desired paths
```

### 3. Run Setup Script (Easiest)

```bash
chmod +x unraid-setup.sh
./unraid-setup.sh
```

**OR** manually with docker-compose:

```bash
docker-compose up -d
```

### 4. Access

Open your browser to:
```
http://YOUR-UNRAID-IP:8080
```

Replace `YOUR-UNRAID-IP` with your Unraid server's IP address.

## First Time Usage

1. **Trigger Initial Index**:
   - Click "Start Index Now" in the Controls section
   - Wait for indexing to complete (time depends on number of files)
   - Status will show "Last Indexed" time when done

2. **Search**:
   - Type a filename or part of a filename
   - Press Enter or click Search
   - Results appear instantly!

3. **Automatic Updates**:
   - Scheduler runs automatically (default: every 6 hours)
   - Adjust in `config.yml` if needed

## Common Paths to Index on Unraid

```yaml
index_paths:
  - "/mnt/user/media"      # Media files
  - "/mnt/user/documents"  # Documents
  - "/mnt/user/downloads"  # Downloads
  - "/mnt/cache"           # Cache drive
```

## Performance Tips

1. **Use Cache Drive** for database:
   ```yaml
   database_path: "/var/lib/plocate/plocate.db"
   ```
   (Already configured in docker-compose to use `/mnt/cache/appdata/plocate-ui/db`)

2. **Index Specific Shares** instead of all of `/mnt/user`:
   - Faster indexing
   - Smaller database
   - More focused searches

3. **Schedule During Low Usage**:
   ```yaml
   interval: "0 3 * * *"  # 3 AM daily
   ```

## Troubleshooting

**Container won't start?**
```bash
docker logs plocate-ui
```

**No search results?**
- Check if indexing completed (Status section)
- Verify paths are correctly mounted
- Trigger manual reindex

**Need to rebuild?**
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## Next Steps

- Read the full [README.md](README.md) for advanced features
- Customize your indexing schedule
- Set up automated backups of the database

## Support

- GitHub Issues: [Report a problem](https://github.com/yourusername/plocate-ui/issues)
- Unraid Forums: [Community thread](#)

---

**That's it!** You now have a powerful file search engine for your Unraid server.
