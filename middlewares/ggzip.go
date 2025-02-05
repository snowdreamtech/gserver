package middlewares

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/snowdreamtech/gserver/pkg/configs"
	"github.com/snowdreamtech/gserver/pkg/tools"
)

var excludedExtentions = []string{
	".png", ".gif", ".jpeg", ".jpg", ".bmp", ".webp",
	".mp3", ".ogg", ".wav", ".wma",
	".3gp", ".avi", ".flv", ".mkv", ".mov", ".mp4", ".rmvb", ".vob", ".webm", ".wmv",
	".exe", ".msi", ".apk", ".pkg", ".pkg", ".dmg", ".ipa", ".deb", ".rpm", ".flatpak", ".snap", ".appimage",
	".rar", ".zip", ".tar", ".gz", ".7z", ".xz", ".bz2", ".iso", ".jar",
}

// Gzip Gzip
func Gzip() gin.HandlerFunc {
	tools.DebugPrintF("[INFO] Starting Middleware %s", "Gzip")

	app := configs.GetAppConfig()

	if !app.Gzip {
		return Empty()
	}

	return gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions(excludedExtentions))
}
