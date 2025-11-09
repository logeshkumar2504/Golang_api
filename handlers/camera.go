package handlers

import (
	"database/sql"
	"net/http"
	"ofella/database"
	"ofella/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateCamera(c *gin.Context) {
	var req models.CreateCameraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := req.Status
	if status == "" {
		status = "offline"
	}

	sql := "INSERT INTO cameras (name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := database.DB.Exec(sql, req.Name, req.IPAddress, req.Port, req.Username, req.Password, req.StreamURL, req.Location, status, req.CameraType, req.Resolution, req.FrameRate, req.Manufacturer, req.Model)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create camera: " + err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get ID"})
		return
	}

	cam := models.Camera{
		ID:           int(id),
		Name:         req.Name,
		IPAddress:    req.IPAddress,
		Port:         req.Port,
		Username:     req.Username,
		Password:     req.Password,
		StreamURL:    req.StreamURL,
		Location:     req.Location,
		Status:       status,
		CameraType:   req.CameraType,
		Resolution:   req.Resolution,
		FrameRate:    req.FrameRate,
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
	}

	c.JSON(201, cam)
}

func GetCamera(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var camera models.Camera
	err = database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", id).Scan(&camera.ID, &camera.Name, &camera.IPAddress, &camera.Port, &camera.Username, &camera.Password, &camera.StreamURL, &camera.Location, &camera.Status, &camera.CameraType, &camera.Resolution, &camera.FrameRate, &camera.Manufacturer, &camera.Model)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, camera)
}

func GetAllCameras(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var cameras []models.Camera
	for rows.Next() {
		var cam models.Camera
		if err := rows.Scan(&cam.ID, &cam.Name, &cam.IPAddress, &cam.Port, &cam.Username, &cam.Password, &cam.StreamURL, &cam.Location, &cam.Status, &cam.CameraType, &cam.Resolution, &cam.FrameRate, &cam.Manufacturer, &cam.Model); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		cameras = append(cameras, cam)
	}

	c.JSON(200, cameras)
}

func UpdateCamera(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.UpdateCameraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if req.Name != nil {
		updates = append(updates, "name = ?")
		args = append(args, *req.Name)
	}
	if req.IPAddress != nil {
		updates = append(updates, "ip_address = ?")
		args = append(args, *req.IPAddress)
	}
	if req.Port != nil {
		updates = append(updates, "port = ?")
		args = append(args, *req.Port)
	}
	if req.Username != nil {
		updates = append(updates, "username = ?")
		args = append(args, *req.Username)
	}
	if req.Password != nil {
		updates = append(updates, "password = ?")
		args = append(args, *req.Password)
	}
	if req.StreamURL != nil {
		updates = append(updates, "stream_url = ?")
		args = append(args, *req.StreamURL)
	}
	if req.Location != nil {
		updates = append(updates, "location = ?")
		args = append(args, *req.Location)
	}
	if req.Status != nil {
		updates = append(updates, "status = ?")
		args = append(args, *req.Status)
	}
	if req.CameraType != nil {
		updates = append(updates, "camera_type = ?")
		args = append(args, *req.CameraType)
	}
	if req.Resolution != nil {
		updates = append(updates, "resolution = ?")
		args = append(args, *req.Resolution)
	}
	if req.FrameRate != nil {
		updates = append(updates, "frame_rate = ?")
		args = append(args, *req.FrameRate)
	}
	if req.Manufacturer != nil {
		updates = append(updates, "manufacturer = ?")
		args = append(args, *req.Manufacturer)
	}
	if req.Model != nil {
		updates = append(updates, "model = ?")
		args = append(args, *req.Model)
	}

	if len(updates) == 0 {
		c.JSON(400, gin.H{"error": "Nothing to update"})
		return
	}

	args = append(args, id)
	query := "UPDATE cameras SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	result, err := database.DB.Exec(query, args...)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(404, gin.H{"error": "Not found"})
		return
	}

	var cam models.Camera
	database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", id).Scan(&cam.ID, &cam.Name, &cam.IPAddress, &cam.Port, &cam.Username, &cam.Password, &cam.StreamURL, &cam.Location, &cam.Status, &cam.CameraType, &cam.Resolution, &cam.FrameRate, &cam.Manufacturer, &cam.Model)

	c.JSON(200, cam)
}

func DeleteCamera(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM cameras WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(404, gin.H{"error": "Not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Deleted"})
}
