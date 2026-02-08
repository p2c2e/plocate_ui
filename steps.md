# Deploying Plocate UI to Unraid (Docker Hub Image)

## Step 1: Create directories on Unraid

SSH into your Unraid server and create the appdata directories:

```bash
ssh root@YOUR-UNRAID-IP
mkdir -p /mnt/cache/appdata/plocate-ui/db
mkdir -p /mnt/cache/appdata/plocate-ui
```

## Step 2: Create config.yml on the Unraid host

```bash
nano /mnt/cache/appdata/plocate-ui/config.yml
```

Paste the following content (edit paths to match your shares):

```yaml
server:
  port: "8080"

plocate:
  indices:
    - name: "media"
      database_path: "/var/lib/plocate/media.db"
      index_paths:
        - "/mnt/media"
      enabled: true

    - name: "documents"
      database_path: "/var/lib/plocate/documents.db"
      index_paths:
        - "/mnt/documents"
      enabled: true

    - name: "downloads"
      database_path: "/var/lib/plocate/downloads.db"
      index_paths:
        - "/mnt/downloads"
      enabled: true

  updatedb_bin: "updatedb"
  plocate_bin: "plocate"

scheduler:
  enabled: true
  interval: "0 */6 * * *"
```

Save and exit (`Ctrl+O`, `Enter`, `Ctrl+X`).

**Important:** The `index_paths` values must match the **Container Path** you set in Step 3 below.

## Step 3: Add the container in Unraid Web UI

1. Open Unraid web UI -> **Docker** tab
2. Click **Add Container**
3. Toggle **Advanced View** (top-right) -- you need this for volume mounts
4. Fill in the basic fields:

| Field | Value |
|---|---|
| Name | `plocate-ui` |
| Repository | `yourdockerhubuser/plocate-ui:latest` |
| Network Type | `Bridge` |
| WebUI | `http://[IP]:[PORT:8080]` |

5. Click **Add another Path, Port, Variable, Label, or Device** and add the following one at a time:

### Port

| Field | Value |
|---|---|
| Config Type | Port |
| Name | `WebUI` |
| Container Port | `8080` |
| Host Port | `8080` |

### Path: Config file

| Field | Value |
|---|---|
| Config Type | Path |
| Name | `Config` |
| Container Path | `/app/config.yml` |
| Host Path | `/mnt/cache/appdata/plocate-ui/config.yml` |
| Access Mode | Read Only |

### Path: Database

| Field | Value |
|---|---|
| Config Type | Path |
| Name | `Database` |
| Container Path | `/var/lib/plocate` |
| Host Path | `/mnt/cache/appdata/plocate-ui/db` |
| Access Mode | Read/Write |

### Path: Media share (example)

| Field | Value |
|---|---|
| Config Type | Path |
| Name | `Media` |
| Container Path | `/mnt/media` |
| Host Path | `/mnt/user/media` |
| Access Mode | Read Only |

### Path: Documents share (example)

| Field | Value |
|---|---|
| Config Type | Path |
| Name | `Documents` |
| Container Path | `/mnt/documents` |
| Host Path | `/mnt/user/documents` |
| Access Mode | Read Only |

### Path: Downloads share (example)

| Field | Value |
|---|---|
| Config Type | Path |
| Name | `Downloads` |
| Container Path | `/mnt/downloads` |
| Host Path | `/mnt/user/downloads` |
| Access Mode | Read Only |

### Variable: Timezone

| Field | Value |
|---|---|
| Config Type | Variable |
| Name | `Timezone` |
| Key | `TZ` |
| Value | `America/New_York` |

6. Click **Apply**

Unraid will automatically pull the image from Docker Hub and start the container.

## Step 4: Verify it's running

1. Go to the **Docker** tab -- you should see `plocate-ui` running with a green icon
2. Click the container icon -> **WebUI** (or browse to `http://YOUR-UNRAID-IP:8080`)
3. Click **Start Index Now** to trigger the first index
4. Wait for indexing to finish, then search

## Editing config.yml later

Edit the file on the host:

```bash
nano /mnt/cache/appdata/plocate-ui/config.yml
```

Then restart the container from the Unraid Docker tab (click the container icon -> **Restart**). The container reads the config at startup, so changes require a restart.

## Adding more shares later

1. Edit `/mnt/cache/appdata/plocate-ui/config.yml` -- add a new index entry with the new container path
2. In Unraid Docker tab -> click `plocate-ui` icon -> **Edit** -> add a new **Path** mapping for the new share
3. Click **Apply** (this recreates the container with the new mount)

**Remember:** The `index_paths` in config.yml must always match the **Container Path** side of your volume mounts, not the host path.
