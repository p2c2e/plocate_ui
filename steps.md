# Deploying Plocate UI to Unraid (Docker Hub Image)

## Step 1: Create directories on Unraid

SSH into your Unraid server and create the appdata directories:

```bash
ssh root@YOUR-UNRAID-IP
mkdir -p /mnt/cache/appdata/plocate-ui/db
mkdir -p /mnt/cache/appdata/plocate-ui/config
```

No config file is needed — the app auto-creates one on first startup and you manage indices through the web UI.

## Step 2: Add the container in Unraid Web UI

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

**Important:** The folder paths you add via the web UI must match the **Container Path** side of your volume mounts, not the host path.

5. Click **Add another Path, Port, Variable, Label, or Device** and add the following one at a time:

### Port

| Field | Value |
|---|---|
| Config Type | Port |
| Name | `WebUI` |
| Container Port | `8080` |
| Host Port | `8080` |

### Path: Config directory

| Field | Value |
|---|---|
| Config Type | Path |
| Name | `Config` |
| Container Path | `/app/config` |
| Host Path | `/mnt/cache/appdata/plocate-ui/config` |
| Access Mode | Read/Write |

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

## Step 3: Verify it's running

1. Go to the **Docker** tab -- you should see `plocate-ui` running with a green icon
2. Click the container icon -> **WebUI** (or browse to `http://YOUR-UNRAID-IP:8080`)
3. In the Controls section, add your indices:
   - Enter a name (e.g. "media") and the container path (e.g. "/mnt/media")
   - Click **Add Index**
   - Repeat for each folder you want to search
4. Click **Start All** to trigger the first index
5. Wait for indexing to finish, then search!

## Adding more shares later

1. In Unraid Docker tab -> click `plocate-ui` icon -> **Edit** -> add a new **Path** mapping for the new share
2. Click **Apply** (this recreates the container with the new mount)
3. In the Plocate UI, add a new index pointing to the container path you just mapped

All index configuration is managed through the web UI and automatically saved — no manual config file editing needed.
