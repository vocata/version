package version

import (
	"path/filepath"
	"strings"

	"code.byted.org/lang/gg/collection/set"
)

// standard
const (
	ExtWhl  = ".whl"
	ExtZip  = ".zip"
	ExtGz   = ".tar.gz"
	ExtTgz  = ".tgz"
	ExtTar  = ".tar"
	ExtBz2  = ".tar.bz2"
	ExtTbz  = ".tbz"
	ExtXz   = ".tar.xz"
	ExtTxz  = ".txz"
	ExtTlz  = ".tlz"
	ExtLz   = ".tar.lz"
	ExtLzma = ".tar.lzma"
)

// legacy
const (
	ExtDeb = ".deb"
	ExtDmg = ".dmg"
	ExtEgg = ".egg"
	ExtExe = ".exe"
	ExtMsi = ".msi"
	ExtRpm = ".rpm"
)

var (
	legacyExt   = set.New(ExtDeb, ExtDmg, ExtEgg, ExtExe, ExtMsi, ExtRpm)
	standardExt = set.New(ExtWhl, ExtZip, ExtGz, ExtTgz, ExtTar, ExtBz2, ExtTbz, ExtXz, ExtTxz, ExtTlz, ExtLz, ExtLzma)
)

func splitFilename(filename string) (string, string) {
	ext := filepath.Ext(filename)
	fragment := strings.TrimSuffix(filename, ext)
	if strings.HasSuffix(strings.ToLower(fragment), ".tar") {
		ext = filepath.Ext(fragment) + ext
		fragment = strings.TrimSuffix(filename, ext)
	}

	return fragment, ext
}
