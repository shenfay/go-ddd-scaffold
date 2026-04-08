package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Device 设备信息
type Device struct {
	DeviceID     string    `json:"device_id"`
	DeviceName   string    `json:"device_name"`
	DeviceType   string    `json:"device_type"`
	LastActiveAt time.Time `json:"last_active_at"`
}

// DeviceStore Redis 设备会话存储实现（Redis）
type DeviceStore struct {
	client *redis.Client
}

// NewDeviceStore 创建设备存储实例
func NewDeviceStore(client *redis.Client) *DeviceStore {
	return &DeviceStore{
		client: client,
	}
}

// AddDevice 添加或更新设备
func (s *DeviceStore) AddDevice(ctx context.Context, userID string, device *Device) error {
	key := "auth:devices:" + userID
	devices, _ := s.GetDevices(ctx, userID)

	// 更新或添加设备
	exists := false
	for i, d := range devices {
		if d.DeviceID == device.DeviceID {
			devices[i] = device
			exists = true
			break
		}
	}
	if !exists {
		devices = append(devices, device)
	}

	// 限制最多 10 个设备
	if len(devices) > 10 {
		devices = devices[:10]
	}

	value, _ := json.Marshal(devices)
	return s.client.Set(ctx, key, value, 30*24*time.Hour).Err()
}

// GetDevices 获取用户的设备列表
func (s *DeviceStore) GetDevices(ctx context.Context, userID string) ([]*Device, error) {
	key := "auth:devices:" + userID
	value, err := s.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return []*Device{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	var devices []*Device
	if err := json.Unmarshal(value, &devices); err != nil {
		return nil, fmt.Errorf("failed to unmarshal devices: %w", err)
	}

	return devices, nil
}

// RemoveDevice 移除设备
func (s *DeviceStore) RemoveDevice(ctx context.Context, userID, deviceID string) error {
	devices, err := s.GetDevices(ctx, userID)
	if err != nil {
		return err
	}

	// 移除指定设备
	filtered := make([]*Device, 0, len(devices))
	for _, d := range devices {
		if d.DeviceID != deviceID {
			filtered = append(filtered, d)
		}
	}

	if len(filtered) == 0 {
		// 如果没有设备了，删除key
		return s.client.Del(ctx, "auth:devices:"+userID).Err()
	}

	value, _ := json.Marshal(filtered)
	return s.client.Set(ctx, "auth:devices:"+userID, value, 30*24*time.Hour).Err()
}
