package auth

// GetCurrentUserID возвращает ID текущего пользователя
func GetCurrentUserID() uint {
	return 1
}

// GetCurrentModeratorID возвращает ID модератора
func GetCurrentModeratorID() uint {
	return 2
}

// IsCurrentUserModerator проверяет, является ли текущий пользователь модератором
func IsCurrentUserModerator() bool {
	return GetCurrentUserID() == 2
}