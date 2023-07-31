package version

import (
	"path/filepath"
	"strings"
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
	LegacyExt   = NewSet(ExtDeb, ExtDmg, ExtEgg, ExtExe, ExtMsi, ExtRpm)
	StandardExt = NewSet(ExtWhl, ExtZip, ExtGz, ExtTgz, ExtTar, ExtBz2, ExtTbz, ExtXz, ExtTxz, ExtTlz, ExtLz, ExtLzma)
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

type Set struct {
	m map[string]struct{}
}

func NewSet(members ...string) *Set {
	s := &Set{}

	s.m = make(map[string]struct{}, len(members))
	for _, v := range members {
		s.m[v] = struct{}{}
	}

	return s
}

func (s *Set) Contains(v string) bool {
	_, ok := s.m[v]
	return ok
}
