package models

type Camera struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	IPAddress    string `json:"ip_address"`
	Port         int    `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	StreamURL    string `json:"stream_url"`
	Location     string `json:"location"`
	Status       string `json:"status"`
	CameraType   string `json:"camera_type"`
	Resolution   string `json:"resolution"`
	FrameRate    int    `json:"frame_rate"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
}

type CreateCameraRequest struct {
	Name         string `json:"name" binding:"required"`
	IPAddress    string `json:"ip_address" binding:"required"`
	Port         int    `json:"port" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	StreamURL    string `json:"stream_url" binding:"required"`
	Location     string `json:"location"`
	Status       string `json:"status"`
	CameraType   string `json:"camera_type"`
	Resolution   string `json:"resolution"`
	FrameRate    int    `json:"frame_rate"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
}

type UpdateCameraRequest struct {
	Name         *string `json:"name"`
	IPAddress    *string `json:"ip_address"`
	Port         *int    `json:"port"`
	Username     *string `json:"username"`
	Password     *string `json:"password"`
	StreamURL    *string `json:"stream_url"`
	Location     *string `json:"location"`
	Status       *string `json:"status"`
	CameraType   *string `json:"camera_type"`
	Resolution   *string `json:"resolution"`
	FrameRate    *int    `json:"frame_rate"`
	Manufacturer *string `json:"manufacturer"`
	Model        *string `json:"model"`
}
