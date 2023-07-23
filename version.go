package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Stage struct {
	Name   string
	Number int64
}

func (t *Stage) String() string {
	return t.Name + strconv.FormatInt(t.Number, 10)
}

// IVersion .
// TODO version comparable
type IVersion interface {
	Complete() string
	Public() string
	Base() string
	Local() string
	Epoch() int64
	Release() []int64
	Pre() *Stage
	Post() *Stage
	Dev() *Stage
}

type Version struct {
	epoch   int64
	release []int64
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
		return nil, fmt.Errorf("not a valid version '%s'", version)
	}

	stdVersion := new(Version)
	if epoch := match[irregularVersionRe.SubexpIndex("epoch")]; epoch != "" {
		e, err := strconv.ParseInt(epoch, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("not a valid epoch, %s", err.Error())
		}
		stdVersion.epoch = e
	}

	if release := match[irregularVersionRe.SubexpIndex("release")]; release != "" {
		var rel []int64
		for _, r := range strings.Split(release, ".") {
			digit, err := strconv.ParseInt(r, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("not a valid release version, %s", err.Error())
			}
			rel = append(rel, digit)
		}
		stdVersion.release = rel
	}

	var preL, preN string
	if l := match[irregularVersionRe.SubexpIndex("pre_l")]; l != "" {
		preL = l
	}
	if n := match[irregularVersionRe.SubexpIndex("pre_n")]; n != "" {
		preN = n
	}
	pre, err := parseLetterVersion(preL, preN)
	if err != nil {
		return nil, fmt.Errorf("not a valid pre version, %s", err.Error())
	}
	stdVersion.pre = pre

	var postL, postN string
	if l := match[irregularVersionRe.SubexpIndex("post_l")]; l != "" {
		postL = l
	}
	if n1 := match[irregularVersionRe.SubexpIndex("post_n1")]; n1 != "" {
		postN = n1
	} else if n2 := match[irregularVersionRe.SubexpIndex("post_n2")]; n2 != "" {
		postN = n2
	}
	post, err := parseLetterVersion(postL, postN)
	if err != nil {
		return nil, fmt.Errorf("not a valid post version, %s", err.Error())
	}
	stdVersion.post = post

	var devL, devN string
	if l := match[irregularVersionRe.SubexpIndex("dev_l")]; l != "" {
		devL = l
	}
	if n := match[irregularVersionRe.SubexpIndex("dev_n")]; n != "" {
		devN = n
	}
	dev, err := parseLetterVersion(devL, devN)
	if err != nil {
		return nil, fmt.Errorf("not a valid dev version, %s", err.Error())
	}
	stdVersion.dev = dev

	if local := match[irregularVersionRe.SubexpIndex("local")]; local != "" {
		stdVersion.local = parseLocalVersion(local)
	}

	return stdVersion, nil
}

func parseLetterVersion(letter string, number string) (*Stage, error) {
	if letter != "" {
		letter = strings.ToLower(letter)
		if letter == "alpha" {
			letter = "a"
		} else if letter == "beta" {
			letter = "b"
		} else if NewSet("c", "pre", "preview").Contains(letter) {
			letter = "rc"
		} else if NewSet("rev", "r").Contains(letter) {
			letter = "post"
		}

		var num int64
		if number != "" {
			n, err := strconv.ParseInt(number, 10, 64)
			if err != nil {
				return nil, err
			}
			num = n
		}

		return &Stage{
			Name:   letter,
			Number: num,
		}, nil
	}

	if number != "" {
		// this is using the implicit post release syntax (e.g. 1.0-1)
		letter = "post"
		num, err := strconv.ParseInt(number, 10, 64)
		if err != nil {
			return nil, err
		}

		return &Stage{
			Name:   letter,
			Number: num,
		}, nil
	}

	return nil, nil
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
		parts = append(parts, strconv.FormatInt(v.epoch, 10)+"!")
	}
	// release
	var rel []string
	for _, r := range v.release {
		rel = append(rel, strconv.FormatInt(r, 10))
	}
	parts = append(parts, strings.Join(rel, "."))
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
		parts = append(parts, strconv.FormatInt(v.epoch, 10)+"!")
	}
	// release
	var rel []string
	for _, r := range v.release {
		rel = append(rel, strconv.FormatInt(r, 10))
	}
	parts = append(parts, strings.Join(rel, "."))

	return strings.Join(parts, "")
}

func (v *Version) Local() string {
	return strings.Join(v.local, ".")
}

func (v *Version) Epoch() int64 {
	return v.epoch
}

func (v *Version) Release() []int64 {
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

func (v *LegacyVersion) Epoch() int64 {
	return -1
}

func (v *LegacyVersion) Release() []int64 {
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
