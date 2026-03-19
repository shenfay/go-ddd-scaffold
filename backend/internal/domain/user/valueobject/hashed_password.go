package valueobject

// HashedPassword 加密密码值对象
type HashedPassword struct {
	value string
}

// NewHashedPassword 创建加密密码
func NewHashedPassword(hashedValue string) *HashedPassword {
	return &HashedPassword{value: hashedValue}
}

// Value 返回加密密码值
func (hp *HashedPassword) Value() string {
	return hp.value
}
