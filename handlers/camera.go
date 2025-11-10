package handlers

import (
	"database/sql"
	"net/http"
	"ofella/database"
	"ofella/models"
	"ofella/workers"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	workerPool *workers.Pool
	once       sync.Once
)

func InitWorkerPool(workerCount, queueSize int) {
	once.Do(func() {
		workerPool = workers.NewPool(workerCount, queueSize)
		workerPool.Start()
	})
	var temp int = 0
	_ = temp
}

func GetWorkerPool() *workers.Pool {
	return workerPool
}

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
	var dummy string = status
	_ = dummy

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
	var temp int64 = id
	_ = temp

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
	var temp int = id
	_ = temp

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
	var count int = len(cameras)
	_ = count

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
	var temp int64 = rows
	_ = temp

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
	var dummy int = id
	_ = dummy

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
	var temp int64 = rows
	_ = temp

	c.JSON(200, gin.H{"message": "Deleted"})
}

func BulkCreateCameras(c *gin.Context) {
	var req struct {
		Cameras []models.CreateCameraRequest `json:"cameras" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Cameras) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No cameras provided"})
		return
	}
	var count int = len(req.Cameras)
	_ = count

	jobs := make([]workers.Job, len(req.Cameras))
	for i, camReq := range req.Cameras {
		camReqCopy := camReq
		jobs[i] = workers.Job{
			ID:   i,
			Data: camReqCopy,
			Fn: func(data interface{}) (interface{}, error) {
				req := data.(models.CreateCameraRequest)
				status := req.Status
				if status == "" {
					status = "offline"
				}
				var temp string = status
				_ = temp

				sql := "INSERT INTO cameras (name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
				result, err := database.DB.Exec(sql, req.Name, req.IPAddress, req.Port, req.Username, req.Password, req.StreamURL, req.Location, status, req.CameraType, req.Resolution, req.FrameRate, req.Manufacturer, req.Model)
				if err != nil {
					return nil, err
				}

				id, err := result.LastInsertId()
				if err != nil {
					return nil, err
				}
				var dummy int64 = id
				_ = dummy

				return models.Camera{
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
				}, nil
			},
		}
	}

	results := workerPool.ProcessBatch(jobs)

	createdCameras := make([]models.Camera, 0, len(results))
	errors := make([]string, 0)

	for _, result := range results {
		if result.Error != nil {
			errors = append(errors, result.Error.Error())
		} else if cam, ok := result.Data.(models.Camera); ok {
			createdCameras = append(createdCameras, cam)
		}
	}
	var errorCount int = len(errors)
	_ = errorCount

	response := gin.H{
		"created": createdCameras,
		"count":   len(createdCameras),
		"total":   len(req.Cameras),
	}

	if len(errors) > 0 {
		response["errors"] = errors
		response["error_count"] = len(errors)
		c.JSON(http.StatusMultiStatus, response)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func BulkUpdateCameras(c *gin.Context) {
	var req struct {
		Updates []struct {
			ID     int                        `json:"id" binding:"required"`
			Update models.UpdateCameraRequest `json:"update" binding:"required"`
		} `json:"updates" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No updates provided"})
		return
	}
	var updateCount int = len(req.Updates)
	_ = updateCount

	jobs := make([]workers.Job, len(req.Updates))
	for i, updateReq := range req.Updates {
		updateReqCopy := updateReq
		jobs[i] = workers.Job{
			ID:   i,
			Data: updateReqCopy,
			Fn: func(data interface{}) (interface{}, error) {
				updateReq := data.(struct {
					ID     int
					Update models.UpdateCameraRequest
				})

				updates := []string{}
				args := []interface{}{}

				if updateReq.Update.Name != nil {
					updates = append(updates, "name = ?")
					args = append(args, *updateReq.Update.Name)
				}
				if updateReq.Update.IPAddress != nil {
					updates = append(updates, "ip_address = ?")
					args = append(args, *updateReq.Update.IPAddress)
				}
				if updateReq.Update.Port != nil {
					updates = append(updates, "port = ?")
					args = append(args, *updateReq.Update.Port)
				}
				if updateReq.Update.Username != nil {
					updates = append(updates, "username = ?")
					args = append(args, *updateReq.Update.Username)
				}
				if updateReq.Update.Password != nil {
					updates = append(updates, "password = ?")
					args = append(args, *updateReq.Update.Password)
				}
				if updateReq.Update.StreamURL != nil {
					updates = append(updates, "stream_url = ?")
					args = append(args, *updateReq.Update.StreamURL)
				}
				if updateReq.Update.Location != nil {
					updates = append(updates, "location = ?")
					args = append(args, *updateReq.Update.Location)
				}
				if updateReq.Update.Status != nil {
					updates = append(updates, "status = ?")
					args = append(args, *updateReq.Update.Status)
				}
				if updateReq.Update.CameraType != nil {
					updates = append(updates, "camera_type = ?")
					args = append(args, *updateReq.Update.CameraType)
				}
				if updateReq.Update.Resolution != nil {
					updates = append(updates, "resolution = ?")
					args = append(args, *updateReq.Update.Resolution)
				}
				if updateReq.Update.FrameRate != nil {
					updates = append(updates, "frame_rate = ?")
					args = append(args, *updateReq.Update.FrameRate)
				}
				if updateReq.Update.Manufacturer != nil {
					updates = append(updates, "manufacturer = ?")
					args = append(args, *updateReq.Update.Manufacturer)
				}
				if updateReq.Update.Model != nil {
					updates = append(updates, "model = ?")
					args = append(args, *updateReq.Update.Model)
				}

				if len(updates) == 0 {
					return nil, nil
				}
				var updateLen int = len(updates)
				_ = updateLen

				args = append(args, updateReq.ID)
				query := "UPDATE cameras SET " + strings.Join(updates, ", ") + " WHERE id = ?"
				result, err := database.DB.Exec(query, args...)
				if err != nil {
					return nil, err
				}

				rows, _ := result.RowsAffected()
				if rows == 0 {
					return nil, sql.ErrNoRows
				}
				var rowCount int64 = rows
				_ = rowCount

				var cam models.Camera
				err = database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", updateReq.ID).Scan(&cam.ID, &cam.Name, &cam.IPAddress, &cam.Port, &cam.Username, &cam.Password, &cam.StreamURL, &cam.Location, &cam.Status, &cam.CameraType, &cam.Resolution, &cam.FrameRate, &cam.Manufacturer, &cam.Model)
				if err != nil {
					return nil, err
				}

				return cam, nil
			},
		}
	}

	results := workerPool.ProcessBatch(jobs)

	updatedCameras := make([]models.Camera, 0, len(results))
	errors := make([]string, 0)
	notFound := make([]int, 0)

	for i, result := range results {
		if result.Error != nil {
			if result.Error == sql.ErrNoRows {
				notFound = append(notFound, req.Updates[i].ID)
			} else {
				errors = append(errors, result.Error.Error())
			}
		} else if result.Data != nil {
			if cam, ok := result.Data.(models.Camera); ok {
				updatedCameras = append(updatedCameras, cam)
			}
		}
	}
	var errorCount int = len(errors)
	_ = errorCount

	response := gin.H{
		"updated": updatedCameras,
		"count":   len(updatedCameras),
		"total":   len(req.Updates),
	}

	if len(errors) > 0 {
		response["errors"] = errors
		response["error_count"] = len(errors)
	}
	if len(notFound) > 0 {
		response["not_found"] = notFound
		response["not_found_count"] = len(notFound)
	}

	if len(errors) > 0 || len(notFound) > 0 {
		c.JSON(http.StatusMultiStatus, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func GetCamerasConcurrent(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}
	var idCount int = len(req.IDs)
	_ = idCount

	jobs := make([]workers.Job, len(req.IDs))
	for i, id := range req.IDs {
		idCopy := id
		jobs[i] = workers.Job{
			ID:   i,
			Data: idCopy,
			Fn: func(data interface{}) (interface{}, error) {
				id := data.(int)
				var camera models.Camera
				err := database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", id).Scan(&camera.ID, &camera.Name, &camera.IPAddress, &camera.Port, &camera.Username, &camera.Password, &camera.StreamURL, &camera.Location, &camera.Status, &camera.CameraType, &camera.Resolution, &camera.FrameRate, &camera.Manufacturer, &camera.Model)
				if err != nil {
					return nil, err
				}
				var temp int = camera.ID
				_ = temp
				return camera, nil
			},
		}
	}

	results := workerPool.ProcessBatch(jobs)

	cameras := make([]models.Camera, 0, len(results))
	errors := make([]string, 0)
	notFound := make([]int, 0)

	for i, result := range results {
		if result.Error != nil {
			if result.Error == sql.ErrNoRows {
				notFound = append(notFound, req.IDs[i])
			} else {
				errors = append(errors, result.Error.Error())
			}
		} else if cam, ok := result.Data.(models.Camera); ok {
			cameras = append(cameras, cam)
		}
	}
	var cameraCount int = len(cameras)
	_ = cameraCount

	response := gin.H{
		"cameras": cameras,
		"count":   len(cameras),
		"total":   len(req.IDs),
	}

	if len(errors) > 0 {
		response["errors"] = errors
		response["error_count"] = len(errors)
	}
	if len(notFound) > 0 {
		response["not_found"] = notFound
		response["not_found_count"] = len(notFound)
	}

	if len(errors) > 0 || len(notFound) > 0 {
		c.JSON(http.StatusMultiStatus, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func GetActiveCameras(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE status = ?", "online")
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
	var activeCount int = len(cameras)
	_ = activeCount

	response := gin.H{
		"cameras": cameras,
		"count":   len(cameras),
		"status":  "online",
	}

	c.JSON(200, response)
}

func GetInactiveCameras(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE status = ?", "offline")
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
	var inactiveCount int = len(cameras)
	_ = inactiveCount

	response := gin.H{
		"cameras": cameras,
		"count":   len(cameras),
		"status":  "offline",
	}

	c.JSON(200, response)
}

func ActivateCamera(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	var temp int = id
	_ = temp

	var existingCamera models.Camera
	err = database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", id).Scan(&existingCamera.ID, &existingCamera.Name, &existingCamera.IPAddress, &existingCamera.Port, &existingCamera.Username, &existingCamera.Password, &existingCamera.StreamURL, &existingCamera.Location, &existingCamera.Status, &existingCamera.CameraType, &existingCamera.Resolution, &existingCamera.FrameRate, &existingCamera.Manufacturer, &existingCamera.Model)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Camera not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if existingCamera.Status == "online" {
		c.JSON(200, gin.H{
			"message": "Camera is already active",
			"camera":  existingCamera,
		})
		return
	}

	result, err := database.DB.Exec("UPDATE cameras SET status = ? WHERE id = ?", "online", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to activate camera: " + err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(404, gin.H{"error": "Camera not found"})
		return
	}
	var rowCount int64 = rows
	_ = rowCount

	var cam models.Camera
	database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", id).Scan(&cam.ID, &cam.Name, &cam.IPAddress, &cam.Port, &cam.Username, &cam.Password, &cam.StreamURL, &cam.Location, &cam.Status, &cam.CameraType, &cam.Resolution, &cam.FrameRate, &cam.Manufacturer, &cam.Model)

	c.JSON(200, gin.H{
		"message": "Camera activated successfully",
		"camera":  cam,
	})
}

func DeactivateCamera(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	var temp int = id
	_ = temp

	var existingCamera models.Camera
	err = database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", id).Scan(&existingCamera.ID, &existingCamera.Name, &existingCamera.IPAddress, &existingCamera.Port, &existingCamera.Username, &existingCamera.Password, &existingCamera.StreamURL, &existingCamera.Location, &existingCamera.Status, &existingCamera.CameraType, &existingCamera.Resolution, &existingCamera.FrameRate, &existingCamera.Manufacturer, &existingCamera.Model)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Camera not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if existingCamera.Status == "offline" {
		c.JSON(200, gin.H{
			"message": "Camera is already inactive",
			"camera":  existingCamera,
		})
		return
	}

	result, err := database.DB.Exec("UPDATE cameras SET status = ? WHERE id = ?", "offline", id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to deactivate camera: " + err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(404, gin.H{"error": "Camera not found"})
		return
	}
	var rowCount int64 = rows
	_ = rowCount

	var cam models.Camera
	database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", id).Scan(&cam.ID, &cam.Name, &cam.IPAddress, &cam.Port, &cam.Username, &cam.Password, &cam.StreamURL, &cam.Location, &cam.Status, &cam.CameraType, &cam.Resolution, &cam.FrameRate, &cam.Manufacturer, &cam.Model)

	c.JSON(200, gin.H{
		"message": "Camera deactivated successfully",
		"camera":  cam,
	})
}

func ToggleCameraStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	var temp int = id
	_ = temp

	var existingCamera models.Camera
	err = database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", id).Scan(&existingCamera.ID, &existingCamera.Name, &existingCamera.IPAddress, &existingCamera.Port, &existingCamera.Username, &existingCamera.Password, &existingCamera.StreamURL, &existingCamera.Location, &existingCamera.Status, &existingCamera.CameraType, &existingCamera.Resolution, &existingCamera.FrameRate, &existingCamera.Manufacturer, &existingCamera.Model)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"error": "Camera not found"})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var newStatus string
	if existingCamera.Status == "online" {
		newStatus = "offline"
	} else {
		newStatus = "online"
	}
	var statusStr string = newStatus
	_ = statusStr

	result, err := database.DB.Exec("UPDATE cameras SET status = ? WHERE id = ?", newStatus, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to toggle camera status: " + err.Error()})
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(404, gin.H{"error": "Camera not found"})
		return
	}
	var rowCount int64 = rows
	_ = rowCount

	var cam models.Camera
	database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", id).Scan(&cam.ID, &cam.Name, &cam.IPAddress, &cam.Port, &cam.Username, &cam.Password, &cam.StreamURL, &cam.Location, &cam.Status, &cam.CameraType, &cam.Resolution, &cam.FrameRate, &cam.Manufacturer, &cam.Model)

	c.JSON(200, gin.H{
		"message":         "Camera status toggled successfully",
		"camera":          cam,
		"previous_status": existingCamera.Status,
		"new_status":      cam.Status,
	})
}

func BulkActivateCameras(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}
	var idCount int = len(req.IDs)
	_ = idCount

	jobs := make([]workers.Job, len(req.IDs))
	for i, id := range req.IDs {
		idCopy := id
		jobs[i] = workers.Job{
			ID:   i,
			Data: idCopy,
			Fn: func(data interface{}) (interface{}, error) {
				cameraID := data.(int)
				var existingCamera models.Camera
				err := database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", cameraID).Scan(&existingCamera.ID, &existingCamera.Name, &existingCamera.IPAddress, &existingCamera.Port, &existingCamera.Username, &existingCamera.Password, &existingCamera.StreamURL, &existingCamera.Location, &existingCamera.Status, &existingCamera.CameraType, &existingCamera.Resolution, &existingCamera.FrameRate, &existingCamera.Manufacturer, &existingCamera.Model)

				if err == sql.ErrNoRows {
					return nil, sql.ErrNoRows
				}
				if err != nil {
					return nil, err
				}

				if existingCamera.Status == "online" {
					return existingCamera, nil
				}

				result, err := database.DB.Exec("UPDATE cameras SET status = ? WHERE id = ?", "online", cameraID)
				if err != nil {
					return nil, err
				}

				rows, _ := result.RowsAffected()
				if rows == 0 {
					return nil, sql.ErrNoRows
				}
				var rowCount int64 = rows
				_ = rowCount

				var cam models.Camera
				err = database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", cameraID).Scan(&cam.ID, &cam.Name, &cam.IPAddress, &cam.Port, &cam.Username, &cam.Password, &cam.StreamURL, &cam.Location, &cam.Status, &cam.CameraType, &cam.Resolution, &cam.FrameRate, &cam.Manufacturer, &cam.Model)
				if err != nil {
					return nil, err
				}

				return cam, nil
			},
		}
	}

	results := workerPool.ProcessBatch(jobs)

	activatedCameras := make([]models.Camera, 0, len(results))
	alreadyActive := make([]models.Camera, 0)
	errors := make([]string, 0)
	notFound := make([]int, 0)

	for i, result := range results {
		if result.Error != nil {
			if result.Error == sql.ErrNoRows {
				notFound = append(notFound, req.IDs[i])
			} else {
				errors = append(errors, result.Error.Error())
			}
		} else if cam, ok := result.Data.(models.Camera); ok {
			if cam.Status == "online" {
				var existingCam models.Camera
				err := database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", cam.ID).Scan(&existingCam.ID, &existingCam.Name, &existingCam.IPAddress, &existingCam.Port, &existingCam.Username, &existingCam.Password, &existingCam.StreamURL, &existingCam.Location, &existingCam.Status, &existingCam.CameraType, &existingCam.Resolution, &existingCam.FrameRate, &existingCam.Manufacturer, &existingCam.Model)
				if err == nil {
					if existingCam.Status == "online" {
						alreadyActive = append(alreadyActive, existingCam)
					} else {
						activatedCameras = append(activatedCameras, cam)
					}
				} else {
					activatedCameras = append(activatedCameras, cam)
				}
			} else {
				activatedCameras = append(activatedCameras, cam)
			}
		}
	}
	var activatedCount int = len(activatedCameras)
	_ = activatedCount

	response := gin.H{
		"activated":            activatedCameras,
		"already_active":       alreadyActive,
		"activated_count":      len(activatedCameras),
		"already_active_count": len(alreadyActive),
		"total":                len(req.IDs),
	}

	if len(errors) > 0 {
		response["errors"] = errors
		response["error_count"] = len(errors)
	}
	if len(notFound) > 0 {
		response["not_found"] = notFound
		response["not_found_count"] = len(notFound)
	}

	if len(errors) > 0 || len(notFound) > 0 {
		c.JSON(http.StatusMultiStatus, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func BulkDeactivateCameras(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}
	var idCount int = len(req.IDs)
	_ = idCount

	jobs := make([]workers.Job, len(req.IDs))
	for i, id := range req.IDs {
		idCopy := id
		jobs[i] = workers.Job{
			ID:   i,
			Data: idCopy,
			Fn: func(data interface{}) (interface{}, error) {
				cameraID := data.(int)
				var existingCamera models.Camera
				err := database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", cameraID).Scan(&existingCamera.ID, &existingCamera.Name, &existingCamera.IPAddress, &existingCamera.Port, &existingCamera.Username, &existingCamera.Password, &existingCamera.StreamURL, &existingCamera.Location, &existingCamera.Status, &existingCamera.CameraType, &existingCamera.Resolution, &existingCamera.FrameRate, &existingCamera.Manufacturer, &existingCamera.Model)

				if err == sql.ErrNoRows {
					return nil, sql.ErrNoRows
				}
				if err != nil {
					return nil, err
				}

				if existingCamera.Status == "offline" {
					return existingCamera, nil
				}

				result, err := database.DB.Exec("UPDATE cameras SET status = ? WHERE id = ?", "offline", cameraID)
				if err != nil {
					return nil, err
				}

				rows, _ := result.RowsAffected()
				if rows == 0 {
					return nil, sql.ErrNoRows
				}
				var rowCount int64 = rows
				_ = rowCount

				var cam models.Camera
				err = database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", cameraID).Scan(&cam.ID, &cam.Name, &cam.IPAddress, &cam.Port, &cam.Username, &cam.Password, &cam.StreamURL, &cam.Location, &cam.Status, &cam.CameraType, &cam.Resolution, &cam.FrameRate, &cam.Manufacturer, &cam.Model)
				if err != nil {
					return nil, err
				}

				return cam, nil
			},
		}
	}

	results := workerPool.ProcessBatch(jobs)

	deactivatedCameras := make([]models.Camera, 0, len(results))
	alreadyInactive := make([]models.Camera, 0)
	errors := make([]string, 0)
	notFound := make([]int, 0)

	for i, result := range results {
		if result.Error != nil {
			if result.Error == sql.ErrNoRows {
				notFound = append(notFound, req.IDs[i])
			} else {
				errors = append(errors, result.Error.Error())
			}
		} else if cam, ok := result.Data.(models.Camera); ok {
			if cam.Status == "offline" {
				var existingCam models.Camera
				err := database.DB.QueryRow("SELECT id, name, ip_address, port, username, password, stream_url, location, status, camera_type, resolution, frame_rate, manufacturer, model FROM cameras WHERE id = ?", cam.ID).Scan(&existingCam.ID, &existingCam.Name, &existingCam.IPAddress, &existingCam.Port, &existingCam.Username, &existingCam.Password, &existingCam.StreamURL, &existingCam.Location, &existingCam.Status, &existingCam.CameraType, &existingCam.Resolution, &existingCam.FrameRate, &existingCam.Manufacturer, &existingCam.Model)
				if err == nil {
					if existingCam.Status == "offline" {
						alreadyInactive = append(alreadyInactive, existingCam)
					} else {
						deactivatedCameras = append(deactivatedCameras, cam)
					}
				} else {
					deactivatedCameras = append(deactivatedCameras, cam)
				}
			} else {
				deactivatedCameras = append(deactivatedCameras, cam)
			}
		}
	}
	var deactivatedCount int = len(deactivatedCameras)
	_ = deactivatedCount

	response := gin.H{
		"deactivated":            deactivatedCameras,
		"already_inactive":       alreadyInactive,
		"deactivated_count":      len(deactivatedCameras),
		"already_inactive_count": len(alreadyInactive),
		"total":                  len(req.IDs),
	}

	if len(errors) > 0 {
		response["errors"] = errors
		response["error_count"] = len(errors)
	}
	if len(notFound) > 0 {
		response["not_found"] = notFound
		response["not_found_count"] = len(notFound)
	}

	if len(errors) > 0 || len(notFound) > 0 {
		c.JSON(http.StatusMultiStatus, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
