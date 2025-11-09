# VMS Camera CRUD API

A RESTful API built with Go and SQLite for managing cameras in a Video Management System (VMS). This API provides full CRUD (Create, Read, Update, Delete) operations for camera resources.

## Features

- ✅ Create a camera (POST)
- ✅ Get all cameras (GET)
- ✅ Get a camera by ID (GET)
- ✅ Update a camera (PUT)
- ✅ Delete a camera (DELETE)
- ✅ SQLite database embedded in the application
- ✅ VMS-specific camera fields (IP address, stream URL, status, etc.)

## Prerequisites

- Go 1.21 or higher
- Postman (or any HTTP client) for testing

## Installation

1. Clone or navigate to the project directory:
```bash
cd ofella
```

2. Install dependencies:
```bash
go mod download
```

## Running the Application

Start the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

## API Endpoints

### Base URL
```
http://localhost:8080/api/cameras
```

### 1. Create Camera (POST)
**Endpoint:** `POST /api/cameras`

**Request Body:**
```json
{
  "name": "Main Entrance Camera",
  "ip_address": "192.168.1.100",
  "port": 554,
  "username": "admin",
  "password": "password123",
  "stream_url": "rtsp://192.168.1.100:554/stream1",
  "location": "Main Entrance",
  "status": "online",
  "camera_type": "IP",
  "resolution": "1080p",
  "frame_rate": 30,
  "manufacturer": "Hikvision",
  "model": "DS-2CD2142FWD-I"
}
```

**Required Fields:**
- `name` - Camera name
- `ip_address` - Camera IP address
- `port` - Camera port (usually 554 for RTSP)
- `username` - Camera username
- `password` - Camera password
- `stream_url` - RTSP stream URL

**Optional Fields:**
- `location` - Physical location of the camera
- `status` - Camera status (online, offline, maintenance) - defaults to "offline"
- `camera_type` - Type of camera (IP, analog, PTZ, dome)
- `resolution` - Video resolution (1080p, 4K, etc.)
- `frame_rate` - Frames per second
- `manufacturer` - Camera manufacturer
- `model` - Camera model

**Response:** `201 Created`
```json
{
  "id": 1,
  "name": "Main Entrance Camera",
  "ip_address": "192.168.1.100",
  "port": 554,
  "username": "admin",
  "password": "password123",
  "stream_url": "rtsp://192.168.1.100:554/stream1",
  "location": "Main Entrance",
  "status": "online",
  "camera_type": "IP",
  "resolution": "1080p",
  "frame_rate": 30,
  "manufacturer": "Hikvision",
  "model": "DS-2CD2142FWD-I"
}
```

### 2. Get All Cameras (GET)
**Endpoint:** `GET /api/cameras`

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "Main Entrance Camera",
    "ip_address": "192.168.1.100",
    "port": 554,
    "username": "admin",
    "password": "password123",
    "stream_url": "rtsp://192.168.1.100:554/stream1",
    "location": "Main Entrance",
    "status": "online",
    "camera_type": "IP",
    "resolution": "1080p",
    "frame_rate": 30,
    "manufacturer": "Hikvision",
    "model": "DS-2CD2142FWD-I"
  }
]
```

### 3. Get Camera by ID (GET)
**Endpoint:** `GET /api/cameras/:id`

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Main Entrance Camera",
  "ip_address": "192.168.1.100",
  "port": 554,
  "username": "admin",
  "password": "password123",
  "stream_url": "rtsp://192.168.1.100:554/stream1",
  "location": "Main Entrance",
  "status": "online",
  "camera_type": "IP",
  "resolution": "1080p",
  "frame_rate": 30,
  "manufacturer": "Hikvision",
  "model": "DS-2CD2142FWD-I"
}
```

### 4. Update Camera (PUT)
**Endpoint:** `PUT /api/cameras/:id`

**Request Body:** (All fields are optional - only include fields you want to update)
```json
{
  "status": "offline",
  "location": "Main Entrance - Updated",
  "frame_rate": 60
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Main Entrance Camera",
  "ip_address": "192.168.1.100",
  "port": 554,
  "username": "admin",
  "password": "password123",
  "stream_url": "rtsp://192.168.1.100:554/stream1",
  "location": "Main Entrance - Updated",
  "status": "offline",
  "camera_type": "IP",
  "resolution": "1080p",
  "frame_rate": 60,
  "manufacturer": "Hikvision",
  "model": "DS-2CD2142FWD-I"
}
```

### 5. Delete Camera (DELETE)
**Endpoint:** `DELETE /api/cameras/:id`

**Response:** `200 OK`
```json
{
  "message": "Camera deleted successfully"
}
```

## Testing with Postman

### Create a Camera
1. Method: `POST`
2. URL: `http://localhost:8080/api/cameras`
3. Headers: `Content-Type: application/json`
4. Body (raw JSON):
```json
{
  "name": "Parking Lot Camera 1",
  "ip_address": "192.168.1.101",
  "port": 554,
  "username": "admin",
  "password": "securepass",
  "stream_url": "rtsp://192.168.1.101:554/stream1",
  "location": "Parking Lot - North",
  "status": "online",
  "camera_type": "PTZ",
  "resolution": "4K",
  "frame_rate": 30,
  "manufacturer": "Axis",
  "model": "P3245-LVE"
}
```

### Get All Cameras
1. Method: `GET`
2. URL: `http://localhost:8080/api/cameras`

### Get Camera by ID
1. Method: `GET`
2. URL: `http://localhost:8080/api/cameras/1`

### Update Camera
1. Method: `PUT`
2. URL: `http://localhost:8080/api/cameras/1`
3. Headers: `Content-Type: application/json`
4. Body (raw JSON):
```json
{
  "status": "maintenance",
  "location": "Main Entrance - Under Maintenance"
}
```

### Delete Camera
1. Method: `DELETE`
2. URL: `http://localhost:8080/api/cameras/1`

## Camera Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Camera name/identifier |
| `ip_address` | string | Yes | Camera IP address |
| `port` | integer | Yes | Camera port (typically 554 for RTSP) |
| `username` | string | Yes | Camera authentication username |
| `password` | string | Yes | Camera authentication password |
| `stream_url` | string | Yes | RTSP stream URL |
| `location` | string | No | Physical location description |
| `status` | string | No | Status: "online", "offline", "maintenance" (default: "offline") |
| `camera_type` | string | No | Type: "IP", "analog", "PTZ", "dome" |
| `resolution` | string | No | Video resolution (e.g., "1080p", "4K") |
| `frame_rate` | integer | No | Frames per second |
| `manufacturer` | string | No | Camera manufacturer (e.g., "Hikvision", "Axis", "Dahua") |
| `model` | string | No | Camera model number |

## Database

The SQLite database file (`cameras.db`) will be created automatically in the project root directory when you first run the application. The database schema includes all VMS-specific camera fields.

## Project Structure

```
ofella/
├── main.go              # Main application entry point
├── go.mod              # Go module file
├── database/
│   └── database.go     # Database initialization and connection
├── models/
│   └── camera.go       # Camera model and request structs
├── handlers/
│   └── camera.go       # HTTP handlers for CRUD operations
└── cameras.db          # SQLite database (created automatically)
```

## Error Responses

All endpoints return appropriate HTTP status codes:
- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Example Use Cases

### Adding a New Camera to VMS
Use POST to add a new camera with its IP address, credentials, and stream URL.

### Monitoring Camera Status
Use GET to retrieve all cameras and check their status (online/offline).

### Updating Camera Configuration
Use PUT to update camera settings like status, location, or stream URL.

### Removing a Camera
Use DELETE to remove a camera from the VMS when it's decommissioned.
