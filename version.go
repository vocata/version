package version

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"code.byted.org/lang/gg/collection/set"
	"code.byted.org/lang/gg/gslice"
)

var (
	ErrUndefinedVersion = errors.New("undefined version")
)

type Stage struct {
	Name   string
	Number int
}

func (t *Stage) String() string {
	return t.Name + strconv.Itoa(t.Number)
}

// IVersion .
// TODO version comparable
type IVersion interface {
	Complete() string
	Public() string
	Base() string
	Local() string
	Epoch() int
	Release() []int
	Pre() *Stage
	Post() *Stage
	Dev() *Stage
}

type Version struct {
	epoch   int
	release []int
	dev     *Stage
	pre     *Stage
	post    *Stage
	local   []string
}

// ParseVersion implements a standard version parser with reference to packaging, an official
// pypi packaging library https://github.com/pypa/packaging/blob/21.3/packaging/version.py#L257.
func ParseVersion(version string) (*Version, error) {
	match := irregularVersionRe.FindStringSubmatch(version)
	if match == nil {
		return nil, fmt.Errorf("'%s' is not a valid version", version)
	}

	stdVersion := &Version{}

	if epoch := match[irregularVersionRe.SubexpIndex("epoch")]; epoch != "" {
		stdVersion.epoch, _ = strconv.Atoi(epoch)
	}

	if release := match[irregularVersionRe.SubexpIndex("release")]; release != "" {
		stdVersion.release = gslice.Map(strings.Split(release, "."), func(r string) int {
			digit, _ := strconv.Atoi(r)
			return digit
		})
	}

	var preL, preN string
	if l := match[irregularVersionRe.SubexpIndex("pre_l")]; l != "" {
		preL = l
	}
	if n := match[irregularVersionRe.SubexpIndex("pre_n")]; n != "" {
		preN = n
	}
	stdVersion.pre = parseLetterVersion(preL, preN)

	var postL, postN string
	if l := match[irregularVersionRe.SubexpIndex("post_l")]; l != "" {
		postL = l
	}
	if n1 := match[irregularVersionRe.SubexpIndex("post_n1")]; n1 != "" {
		postN = n1
	} else if n2 := match[irregularVersionRe.SubexpIndex("post_n2")]; n2 != "" {
		postN = n2
	}
	stdVersion.post = parseLetterVersion(postL, postN)

	var devL, devN string
	if l := match[irregularVersionRe.SubexpIndex("dev_l")]; l != "" {
		devL = l
	}
	if n := match[irregularVersionRe.SubexpIndex("dev_n")]; n != "" {
		devN = n
	}
	stdVersion.dev = parseLetterVersion(devL, devN)

	if local := match[irregularVersionRe.SubexpIndex("local")]; local != "" {
		stdVersion.local = parseLocalVersion(local)
	}

	return stdVersion, nil
}

func parseLetterVersion(letter string, number string) *Stage {
	if letter != "" {
		letter = strings.ToLower(letter)
		if letter == "alpha" {
			letter = "a"
		} else if letter == "beta" {
			letter = "b"
		} else if set.New("c", "pre", "preview").Contains(letter) {
			letter = "rc"
		} else if set.New("rev", "r").Contains(letter) {
			letter = "post"
		}

		var num int
		if number != "" {
			num, _ = strconv.Atoi(number)
		}

		return &Stage{
			Name:   letter,
			Number: num,
		}
	}

	if number != "" {
		// this is using the implicit post release syntax (e.g. 1.0-1)
		letter = "post"
		num, _ := strconv.Atoi(number)

		return &Stage{
			Name:   letter,
			Number: num,
		}
	}

	return nil
}

var localVersionSeparators = regexp.MustCompile(`[-_.]`)

func parseLocalVersion(local string) []string {
	if local != "" {
		local = strings.ToLower(local)
		return localVersionSeparators.Split(local, -1)
	}

	return nil
}

func (v *Version) String() string {
	return fmt.Sprintf("IVersion<%s>", v.Complete())
}

// Complete returns complete version.
func (v *Version) Complete() string {
	var parts []string

	// epoch
	if v.epoch != 0 {
		parts = append(parts, strconv.Itoa(v.epoch)+"!")
	}
	// release
	parts = append(parts, strings.Join(gslice.Map(v.release, strconv.Itoa), "."))
	// pre-release
	if v.pre != nil {
		parts = append(parts, v.pre.String())
	}
	// post-release
	if v.post != nil {
		parts = append(parts, "."+v.post.String())
	}
	// dev-release
	if v.dev != nil {
		parts = append(parts, "."+v.dev.String())
	}
	// local version segment
	if len(v.local) != 0 {
		parts = append(parts, "+"+v.Local())
	}

	return strings.Join(parts, "")
}

// Public returns version except local segment.
func (v *Version) Public() string {
	return strings.SplitN(v.Complete(), "+", 2)[0]
}

// Base returns base version only including epoch and release
func (v *Version) Base() string {
	var parts []string

	// epoch
	if v.epoch != 0 {
		parts = append(parts, strconv.Itoa(v.epoch)+"!")
	}
	// release
	parts = append(parts, strings.Join(gslice.Map(v.release, strconv.Itoa), "."))

	return strings.Join(parts, "")
}

func (v *Version) Local() string {
	return strings.Join(v.local, ".")
}

func (v *Version) Epoch() int {
	return v.epoch
}

func (v *Version) Release() []int {
	return v.release
}

func (v *Version) Pre() *Stage {
	return v.pre
}

func (v *Version) Post() *Stage {
	return v.post
}

func (v *Version) Dev() *Stage {
	return v.dev
}

// LegacyVersion is an irregular version for the compatibility of bdist_dumb format, an old-aged
// pypi package format, for details: https://peps.python.org/pep-0527/#bdist-dumb. this type of
// version has been deprecated in packaging, but now still available in pip.
type LegacyVersion struct {
	version string
}

// ParseLegacyVersion implements a legacy version parser with reference to packaging, an official
// pypi packaging library https://github.com/pypa/packaging/blob/21.3/packaging/version.py#L106.
func ParseLegacyVersion(version string) (*LegacyVersion, error) {
	return &LegacyVersion{
		version: version,
	}, nil
}

func (v *LegacyVersion) String() string {
	return fmt.Sprintf("LegacyVersion<%s>", v.version)
}

func (v *LegacyVersion) Complete() string {
	return v.version
}

func (v *LegacyVersion) Public() string {
	return v.version
}

func (v *LegacyVersion) Base() string {
	return v.version
}

func (v *LegacyVersion) Local() string {
	return ""
}

func (v *LegacyVersion) Epoch() int {
	return -1
}

func (v *LegacyVersion) Release() []int {
	return nil
}

func (v *LegacyVersion) Pre() *Stage {
	return nil
}

func (v *LegacyVersion) Post() *Stage {
	return nil
}

func (v *LegacyVersion) Dev() *Stage {
	return nil
}

// Parse canonicalizes a Version from original version string, fallback to LegacyVersion if version
// string is irregular, https://github.com/pypa/packaging/blob/21.3/packaging/version.py#L42.
func Parse(version string) (IVersion, error) {
	if v, err := ParseVersion(version); err == nil {
		return v, err
	}

	return ParseLegacyVersion(version)
}
