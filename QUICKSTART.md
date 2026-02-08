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

### 2. Configure docker-compose.yml

Edit volume mounts to expose the folders you want to search:

```bash
nano docker-compose.yml
```

Ensure you have a writable config volume and your share mounts:
```yaml
volumes:
  - /mnt/cache/appdata/plocate-ui/config:/app/config
  - /mnt/cache/appdata/plocate-ui/db:/var/lib/plocate
  - /mnt/user:/mnt/user:ro
```

### 3. Build and Run

```bash
docker-compose up -d
```

The app auto-creates its config on first startup — no `config.yml` needed!

### 4. Access

Open your browser to:
```
http://YOUR-UNRAID-IP:8080
```

Replace `YOUR-UNRAID-IP` with your Unraid server's IP address.

## First Time Usage

1. **Add Folders to Index**:
   - In the Controls section, enter an index name (e.g. "media") and folder path (e.g. "/mnt/user/media")
   - Click "Add Index" — the folder is saved and persists across restarts

2. **Trigger Initial Index**:
   - Click "Start All" or Start on the individual index
   - Wait for indexing to complete (time depends on number of files)
   - Status will show "Last Indexed" time when done

3. **Search**:
   - Type a filename or part of a filename
   - Press Enter or click Search
   - Results appear instantly!

4. **Automatic Updates**:
   - Scheduler runs automatically (default: every 6 hours)

## Common Paths to Index on Unraid

Add these via the UI Controls section:

- `/mnt/user/media` — Media files
- `/mnt/user/documents` — Documents
- `/mnt/user/downloads` — Downloads
- `/mnt/cache` — Cache drive

## Performance Tips

1. **Use Cache Drive** for database (already configured in docker-compose to use `/mnt/cache/appdata/plocate-ui/db`)

2. **Index Specific Shares** instead of all of `/mnt/user`:
   - Faster indexing
   - Smaller database
   - More focused searches

3. **Schedule During Low Usage**: Set `INDEX_INTERVAL=0 3 * * *` env var for 3 AM daily

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
