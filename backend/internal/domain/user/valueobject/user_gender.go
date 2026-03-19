package valueobject

// UserGender 用户性别枚举
type UserGender int

const (
	UserGenderUnknown UserGender = iota
	UserGenderMale
	UserGenderFemale
	UserGenderOther
)

// String 返回性别字符串表示
func (ug UserGender) String() string {
	switch ug {
	case UserGenderMale:
		return "male"
	case UserGenderFemale:
		return "female"
	case UserGenderOther:
		return "other"
	default:
		return "unknown"
	}
}
