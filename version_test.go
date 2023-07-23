package version

import (
	"fmt"
	"reflect"
	"testing"
)

// the test cases come from https://github.com/pypa/packaging/blob/21.3/tests/test_version.py

var validVersions = []string{
	//  Implicit epoch of 0
	"1.0.dev456",
	"1.0a1",
	"1.0a2.dev456",
	"1.0a12.dev456",
	"1.0a12",
	"1.0b1.dev456",
	"1.0b2",
	"1.0b2.post345.dev456",
	"1.0b2.post345",
	"1.0b2-346",
	"1.0c1.dev456",
	"1.0c1",
	"1.0rc2",
	"1.0c3",
	"1.0",
	"1.0.post456.dev34",
	"1.0.post456",
	"1.1.dev1",
	"1.2+123abc",
	"1.2+123abc456",
	"1.2+abc",
	"1.2+abc123",
	"1.2+abc123def",
	"1.2+1234.abc",
	"1.2+123456",
	"1.2.r32+123456",
	"1.2.rev33+123456",
	// Explicit epoch of 1
	"1!1.0.dev456",
	"1!1.0a1",
	"1!1.0a2.dev456",
	"1!1.0a12.dev456",
	"1!1.0a12",
	"1!1.0b1.dev456",
	"1!1.0b2",
	"1!1.0b2.post345.dev456",
	"1!1.0b2.post345",
	"1!1.0b2-346",
	"1!1.0c1.dev456",
	"1!1.0c1",
	"1!1.0rc2",
	"1!1.0c3",
	"1!1.0",
	"1!1.0.post456.dev34",
	"1!1.0.post456",
	"1!1.1.dev1",
	"1!1.2+123abc",
	"1!1.2+123abc456",
	"1!1.2+abc",
	"1!1.2+abc123",
	"1!1.2+abc123def",
	"1!1.2+1234.abc",
	"1!1.2+123456",
	"1!1.2.r32+123456",
	"1!1.2.rev33+123456",
}

func TestValidVersion(t *testing.T) {
	for _, version := range validVersions {
		_, err := ParseVersion(version)
		if err != nil {
			t.Error(err)
		}
	}
}

var invalidVersions = []string{
	// nonsensical versions should be invalid
	"french toast",
	// Versions with invalid local versions
	"1.0+a+",
	"1.0++",
	"1.0+_foobar",
	"1.0+foo&asd",
	"1.0+1+1",
}

func TestInvalidVersion(t *testing.T) {
	for _, version := range invalidVersions {
		if _, err := ParseVersion(version); err != nil {
			if err == nil {
				t.Error("should be invalid version")
			}
		}
	}
}

var normalizedVersions = []struct {
	original   string
	normalized string
}{
	// Various development release incarnations
	{"1.0dev", "1.0.dev0"},
	{"1.0.dev", "1.0.dev0"},
	{"1.0dev1", "1.0.dev1"},
	{"1.0dev", "1.0.dev0"},
	{"1.0-dev", "1.0.dev0"},
	{"1.0-dev1", "1.0.dev1"},
	{"1.0DEV", "1.0.dev0"},
	{"1.0.DEV", "1.0.dev0"},
	{"1.0DEV1", "1.0.dev1"},
	{"1.0DEV", "1.0.dev0"},
	{"1.0.DEV1", "1.0.dev1"},
	{"1.0-DEV", "1.0.dev0"},
	{"1.0-DEV1", "1.0.dev1"},
	// Various alpha incarnations
	{"1.0a", "1.0a0"},
	{"1.0.a", "1.0a0"},
	{"1.0.a1", "1.0a1"},
	{"1.0-a", "1.0a0"},
	{"1.0-a1", "1.0a1"},
	{"1.0alpha", "1.0a0"},
	{"1.0.alpha", "1.0a0"},
	{"1.0.alpha1", "1.0a1"},
	{"1.0-alpha", "1.0a0"},
	{"1.0-alpha1", "1.0a1"},
	{"1.0A", "1.0a0"},
	{"1.0.A", "1.0a0"},
	{"1.0.A1", "1.0a1"},
	{"1.0-A", "1.0a0"},
	{"1.0-A1", "1.0a1"},
	{"1.0ALPHA", "1.0a0"},
	{"1.0.ALPHA", "1.0a0"},
	{"1.0.ALPHA1", "1.0a1"},
	{"1.0-ALPHA", "1.0a0"},
	{"1.0-ALPHA1", "1.0a1"},
	// Various beta incarnations
	{"1.0b", "1.0b0"},
	{"1.0.b", "1.0b0"},
	{"1.0.b1", "1.0b1"},
	{"1.0-b", "1.0b0"},
	{"1.0-b1", "1.0b1"},
	{"1.0beta", "1.0b0"},
	{"1.0.beta", "1.0b0"},
	{"1.0.beta1", "1.0b1"},
	{"1.0-beta", "1.0b0"},
	{"1.0-beta1", "1.0b1"},
	{"1.0B", "1.0b0"},
	{"1.0.B", "1.0b0"},
	{"1.0.B1", "1.0b1"},
	{"1.0-B", "1.0b0"},
	{"1.0-B1", "1.0b1"},
	{"1.0BETA", "1.0b0"},
	{"1.0.BETA", "1.0b0"},
	{"1.0.BETA1", "1.0b1"},
	{"1.0-BETA", "1.0b0"},
	{"1.0-BETA1", "1.0b1"},
	// Various release candidate incarnations
	{"1.0c", "1.0rc0"},
	{"1.0.c", "1.0rc0"},
	{"1.0.c1", "1.0rc1"},
	{"1.0-c", "1.0rc0"},
	{"1.0-c1", "1.0rc1"},
	{"1.0rc", "1.0rc0"},
	{"1.0.rc", "1.0rc0"},
	{"1.0.rc1", "1.0rc1"},
	{"1.0-rc", "1.0rc0"},
	{"1.0-rc1", "1.0rc1"},
	{"1.0C", "1.0rc0"},
	{"1.0.C", "1.0rc0"},
	{"1.0.C1", "1.0rc1"},
	{"1.0-C", "1.0rc0"},
	{"1.0-C1", "1.0rc1"},
	{"1.0RC", "1.0rc0"},
	{"1.0.RC", "1.0rc0"},
	{"1.0.RC1", "1.0rc1"},
	{"1.0-RC", "1.0rc0"},
	{"1.0-RC1", "1.0rc1"},
	// Various post release incarnations
	{"1.0post", "1.0.post0"},
	{"1.0.post", "1.0.post0"},
	{"1.0post1", "1.0.post1"},
	{"1.0post", "1.0.post0"},
	{"1.0-post", "1.0.post0"},
	{"1.0-post1", "1.0.post1"},
	{"1.0POST", "1.0.post0"},
	{"1.0.POST", "1.0.post0"},
	{"1.0POST1", "1.0.post1"},
	{"1.0POST", "1.0.post0"},
	{"1.0r", "1.0.post0"},
	{"1.0rev", "1.0.post0"},
	{"1.0.POST1", "1.0.post1"},
	{"1.0.r1", "1.0.post1"},
	{"1.0.rev1", "1.0.post1"},
	{"1.0-POST", "1.0.post0"},
	{"1.0-POST1", "1.0.post1"},
	{"1.0-5", "1.0.post5"},
	{"1.0-r5", "1.0.post5"},
	{"1.0-rev5", "1.0.post5"},
	// Local version case insensitivity
	{"1.0+AbC", "1.0+abc"},
	// Integer Normalization
	{"1.01", "1.1"},
	{"1.0a05", "1.0a5"},
	{"1.0b07", "1.0b7"},
	{"1.0c056", "1.0rc56"},
	{"1.0rc09", "1.0rc9"},
	{"1.0.post000", "1.0.post0"},
	{"1.1.dev09000", "1.1.dev9000"},
	{"00!1.2", "1.2"},
	{"0100!0.0", "100!0.0"},
	// Various other normalizations
	{"v1.0", "1.0"},
	{"   v1.0\t\n", "1.0"},
}

func TestNormalizedVersion(t *testing.T) {
	for _, version := range normalizedVersions {
		v, err := ParseVersion(version.original)
		if err != nil {
			t.Error(err)
			continue
		}
		if v.Complete() != version.normalized {
			t.Errorf("original: %s, %s != %s", version.original, v.Complete(), version.normalized)
		}
	}
}

var versionStrings = []struct {
	version  string
	expected string
}{
	{"1.0.dev456", "1.0.dev456"},
	{"1.0a1", "1.0a1"},
	{"1.0a2.dev456", "1.0a2.dev456"},
	{"1.0a12.dev456", "1.0a12.dev456"},
	{"1.0a12", "1.0a12"},
	{"1.0b1.dev456", "1.0b1.dev456"},
	{"1.0b2", "1.0b2"},
	{"1.0b2.post345.dev456", "1.0b2.post345.dev456"},
	{"1.0b2.post345", "1.0b2.post345"},
	{"1.0rc1.dev456", "1.0rc1.dev456"},
	{"1.0rc1", "1.0rc1"},
	{"1.0", "1.0"},
	{"1.0.post456.dev34", "1.0.post456.dev34"},
	{"1.0.post456", "1.0.post456"},
	{"1.0.1", "1.0.1"},
	{"0!1.0.2", "1.0.2"},
	{"1.0.3+7", "1.0.3+7"},
	{"0!1.0.4+8.0", "1.0.4+8.0"},
	{"1.0.5+9.5", "1.0.5+9.5"},
	{"1.2+1234.abc", "1.2+1234.abc"},
	{"1.2+123456", "1.2+123456"},
	{"1.2+123abc", "1.2+123abc"},
	{"1.2+123abc456", "1.2+123abc456"},
	{"1.2+abc", "1.2+abc"},
	{"1.2+abc123", "1.2+abc123"},
	{"1.2+abc123def", "1.2+abc123def"},
	{"1.1.dev1", "1.1.dev1"},
	{"7!1.0.dev456", "7!1.0.dev456"},
	{"7!1.0a1", "7!1.0a1"},
	{"7!1.0a2.dev456", "7!1.0a2.dev456"},
	{"7!1.0a12.dev456", "7!1.0a12.dev456"},
	{"7!1.0a12", "7!1.0a12"},
	{"7!1.0b1.dev456", "7!1.0b1.dev456"},
	{"7!1.0b2", "7!1.0b2"},
	{"7!1.0b2.post345.dev456", "7!1.0b2.post345.dev456"},
	{"7!1.0b2.post345", "7!1.0b2.post345"},
	{"7!1.0rc1.dev456", "7!1.0rc1.dev456"},
	{"7!1.0rc1", "7!1.0rc1"},
	{"7!1.0", "7!1.0"},
	{"7!1.0.post456.dev34", "7!1.0.post456.dev34"},
	{"7!1.0.post456", "7!1.0.post456"},
	{"7!1.0.1", "7!1.0.1"},
	{"7!1.0.2", "7!1.0.2"},
	{"7!1.0.3+7", "7!1.0.3+7"},
	{"7!1.0.4+8.0", "7!1.0.4+8.0"},
	{"7!1.0.5+9.5", "7!1.0.5+9.5"},
	{"7!1.1.dev1", "7!1.1.dev1"},
}

func TestVersionString(t *testing.T) {
	for _, str := range versionStrings {
		var v IVersion
		v, err := ParseVersion(str.version)
		if err != nil {
			t.Error(err)
			continue
		}
		if v.Complete() != str.expected {
			t.Errorf("version: %s, %s != %s", str.version, v.Complete(), str.expected)
		}
		if repr, expected := fmt.Sprint(v), fmt.Sprintf("IVersion<%s>", str.expected); repr != expected {
			t.Errorf("version: %s, %s != %s", str.version, repr, expected)
		}
	}
}

func TestVersionNormalizedRc(t *testing.T) {
	rc, err := ParseVersion("1.0rc1")
	if err != nil {
		t.Error(err)
		return
	}
	c, err := ParseVersion("1.0c1")
	if err != nil {
		t.Error(err)
		return
	}
	if rc.Complete() != c.Complete() {
		t.Errorf("%s != %s", rc.Complete(), rc.Complete())
	}
}

var versionPublics = []struct {
	version string
	public  string
}{
	{"1.0", "1.0"},
	{"1.0.dev0", "1.0.dev0"},
	{"1.0.dev6", "1.0.dev6"},
	{"1.0a1", "1.0a1"},
	{"1.0a1.post5", "1.0a1.post5"},
	{"1.0a1.post5.dev6", "1.0a1.post5.dev6"},
	{"1.0rc4", "1.0rc4"},
	{"1.0.post5", "1.0.post5"},
	{"1!1.0", "1!1.0"},
	{"1!1.0.dev6", "1!1.0.dev6"},
	{"1!1.0a1", "1!1.0a1"},
	{"1!1.0a1.post5", "1!1.0a1.post5"},
	{"1!1.0a1.post5.dev6", "1!1.0a1.post5.dev6"},
	{"1!1.0rc4", "1!1.0rc4"},
	{"1!1.0.post5", "1!1.0.post5"},
	{"1.0+deadbeef", "1.0"},
	{"1.0.dev6+deadbeef", "1.0.dev6"},
	{"1.0a1+deadbeef", "1.0a1"},
	{"1.0a1.post5+deadbeef", "1.0a1.post5"},
	{"1.0a1.post5.dev6+deadbeef", "1.0a1.post5.dev6"},
	{"1.0rc4+deadbeef", "1.0rc4"},
	{"1.0.post5+deadbeef", "1.0.post5"},
	{"1!1.0+deadbeef", "1!1.0"},
	{"1!1.0.dev6+deadbeef", "1!1.0.dev6"},
	{"1!1.0a1+deadbeef", "1!1.0a1"},
	{"1!1.0a1.post5+deadbeef", "1!1.0a1.post5"},
	{"1!1.0a1.post5.dev6+deadbeef", "1!1.0a1.post5.dev6"},
	{"1!1.0rc4+deadbeef", "1!1.0rc4"},
	{"1!1.0.post5+deadbeef", "1!1.0.post5"},
}

func TestVersionPublic(t *testing.T) {
	for _, public := range versionPublics {
		v, err := ParseVersion(public.version)
		if err != nil {
			t.Error(err)
			continue
		}
		if v.Public() != public.public {
			t.Errorf("version: %s, %s != %s", public.version, v.Public(), public.public)
		}
	}
}

var versionBases = []struct {
	version string
	base    string
}{
	{"1.0", "1.0"},
	{"1.0.dev0", "1.0"},
	{"1.0.dev6", "1.0"},
	{"1.0a1", "1.0"},
	{"1.0a1.post5", "1.0"},
	{"1.0a1.post5.dev6", "1.0"},
	{"1.0rc4", "1.0"},
	{"1.0.post5", "1.0"},
	{"1!1.0", "1!1.0"},
	{"1!1.0.dev6", "1!1.0"},
	{"1!1.0a1", "1!1.0"},
	{"1!1.0a1.post5", "1!1.0"},
	{"1!1.0a1.post5.dev6", "1!1.0"},
	{"1!1.0rc4", "1!1.0"},
	{"1!1.0.post5", "1!1.0"},
	{"1.0+deadbeef", "1.0"},
	{"1.0.dev6+deadbeef", "1.0"},
	{"1.0a1+deadbeef", "1.0"},
	{"1.0a1.post5+deadbeef", "1.0"},
	{"1.0a1.post5.dev6+deadbeef", "1.0"},
	{"1.0rc4+deadbeef", "1.0"},
	{"1.0.post5+deadbeef", "1.0"},
	{"1!1.0+deadbeef", "1!1.0"},
	{"1!1.0.dev6+deadbeef", "1!1.0"},
	{"1!1.0a1+deadbeef", "1!1.0"},
	{"1!1.0a1.post5+deadbeef", "1!1.0"},
	{"1!1.0a1.post5.dev6+deadbeef", "1!1.0"},
	{"1!1.0rc4+deadbeef", "1!1.0"},
	{"1!1.0.post5+deadbeef", "1!1.0"},
}

func TestVersionBase(t *testing.T) {
	for _, base := range versionBases {
		v, err := ParseVersion(base.version)
		if err != nil {
			t.Error(err)
			continue
		}
		if v.Base() != base.base {
			t.Errorf("version: %s, %s != %s", base.version, v.Base(), base.base)
		}
	}
}

var versionEpochs = []struct {
	version string
	epoch   int64
}{
	{"1.0", 0},
	{"1.0.dev0", 0},
	{"1.0.dev6", 0},
	{"1.0a1", 0},
	{"1.0a1.post5", 0},
	{"1.0a1.post5.dev6", 0},
	{"1.0rc4", 0},
	{"1.0.post5", 0},
	{"1!1.0", 1},
	{"1!1.0.dev6", 1},
	{"1!1.0a1", 1},
	{"1!1.0a1.post5", 1},
	{"1!1.0a1.post5.dev6", 1},
	{"1!1.0rc4", 1},
	{"1!1.0.post5", 1},
	{"1.0+deadbeef", 0},
	{"1.0.dev6+deadbeef", 0},
	{"1.0a1+deadbeef", 0},
	{"1.0a1.post5+deadbeef", 0},
	{"1.0a1.post5.dev6+deadbeef", 0},
	{"1.0rc4+deadbeef", 0},
	{"1.0.post5+deadbeef", 0},
	{"1!1.0+deadbeef", 1},
	{"1!1.0.dev6+deadbeef", 1},
	{"1!1.0a1+deadbeef", 1},
	{"1!1.0a1.post5+deadbeef", 1},
	{"1!1.0a1.post5.dev6+deadbeef", 1},
	{"1!1.0rc4+deadbeef", 1},
	{"1!1.0.post5+deadbeef", 1},
}

func TestVersionEpoch(t *testing.T) {
	for _, epoch := range versionEpochs {
		v, err := ParseVersion(epoch.version)
		if err != nil {
			t.Error(err)
			continue
		}
		if v.Epoch() != epoch.epoch {
			t.Errorf("version: %s, %d != %d", epoch.version, v.Epoch(), epoch.epoch)
		}
	}
}

var versionReleases = []struct {
	version string
	release []int64
}{
	{"1.0", []int64{1, 0}},
	{"1.0.dev0", []int64{1, 0}},
	{"1.0.dev6", []int64{1, 0}},
	{"1.0a1", []int64{1, 0}},
	{"1.0a1.post5", []int64{1, 0}},
	{"1.0a1.post5.dev6", []int64{1, 0}},
	{"1.0rc4", []int64{1, 0}},
	{"1.0.post5", []int64{1, 0}},
	{"1!1.0", []int64{1, 0}},
	{"1!1.0.dev6", []int64{1, 0}},
	{"1!1.0a1", []int64{1, 0}},
	{"1!1.0a1.post5", []int64{1, 0}},
	{"1!1.0a1.post5.dev6", []int64{1, 0}},
	{"1!1.0rc4", []int64{1, 0}},
	{"1!1.0.post5", []int64{1, 0}},
	{"1.0+deadbeef", []int64{1, 0}},
	{"1.0.dev6+deadbeef", []int64{1, 0}},
	{"1.0a1+deadbeef", []int64{1, 0}},
	{"1.0a1.post5+deadbeef", []int64{1, 0}},
	{"1.0a1.post5.dev6+deadbeef", []int64{1, 0}},
	{"1.0rc4+deadbeef", []int64{1, 0}},
	{"1.0.post5+deadbeef", []int64{1, 0}},
	{"1!1.0+deadbeef", []int64{1, 0}},
	{"1!1.0.dev6+deadbeef", []int64{1, 0}},
	{"1!1.0a1+deadbeef", []int64{1, 0}},
	{"1!1.0a1.post5+deadbeef", []int64{1, 0}},
	{"1!1.0a1.post5.dev6+deadbeef", []int64{1, 0}},
	{"1!1.0rc4+deadbeef", []int64{1, 0}},
	{"1!1.0.post5+deadbeef", []int64{1, 0}},
}

func TestVersionRelease(t *testing.T) {
	for _, release := range versionReleases {
		v, err := ParseVersion(release.version)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(v.Release(), release.release) {
			t.Errorf("version: %s, %v != %v", release.version, v.Release(), release.release)
		}
	}
}

var versionPres = []struct {
	version string
	pre     *Stage
}{
	{"1.0", nil},
	{"1.0.dev0", nil},
	{"1.0.dev6", nil},
	{"1.0a1", &Stage{"a", 1}},
	{"1.0a1.post5", &Stage{"a", 1}},
	{"1.0a1.post5.dev6", &Stage{"a", 1}},
	{"1.0rc4", &Stage{"rc", 4}},
	{"1.0.post5", nil},
	{"1!1.0", nil},
	{"1!1.0.dev6", nil},
	{"1!1.0a1", &Stage{"a", 1}},
	{"1!1.0a1.post5", &Stage{"a", 1}},
	{"1!1.0a1.post5.dev6", &Stage{"a", 1}},
	{"1!1.0rc4", &Stage{"rc", 4}},
	{"1!1.0.post5", nil},
	{"1.0+deadbeef", nil},
	{"1.0.dev6+deadbeef", nil},
	{"1.0a1+deadbeef", &Stage{"a", 1}},
	{"1.0a1.post5+deadbeef", &Stage{"a", 1}},
	{"1.0a1.post5.dev6+deadbeef", &Stage{"a", 1}},
	{"1.0rc4+deadbeef", &Stage{"rc", 4}},
	{"1.0.post5+deadbeef", nil},
	{"1!1.0+deadbeef", nil},
	{"1!1.0.dev6+deadbeef", nil},
	{"1!1.0a1+deadbeef", &Stage{"a", 1}},
	{"1!1.0a1.post5+deadbeef", &Stage{"a", 1}},
	{"1!1.0a1.post5.dev6+deadbeef", &Stage{"a", 1}},
	{"1!1.0rc4+deadbeef", &Stage{"rc", 4}},
	{"1!1.0.post5+deadbeef", nil},
}

func TestVersionPre(t *testing.T) {
	for _, pre := range versionPres {
		v, err := ParseVersion(pre.version)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(v.Pre(), pre.pre) {
			t.Errorf("version: %s, %v != %v", pre.version, v.Pre(), pre.pre)
		}
	}
}

var versionPosts = []struct {
	version string
	post    *Stage
}{
	{"1.0", nil},
	{"1.0.dev0", nil},
	{"1.0.dev6", nil},
	{"1.0a1", nil},
	{"1.0a1.post5", &Stage{"post", 5}},
	{"1.0a1.post5.dev6", &Stage{"post", 5}},
	{"1.0rc4", nil},
	{"1.0.post5", &Stage{"post", 5}},
	{"1!1.0", nil},
	{"1!1.0.dev6", nil},
	{"1!1.0a1", nil},
	{"1!1.0a1.post5", &Stage{"post", 5}},
	{"1!1.0a1.post5.dev6", &Stage{"post", 5}},
	{"1!1.0rc4", nil},
	{"1!1.0.post5", &Stage{"post", 5}},
	{"1.0+deadbeef", nil},
	{"1.0.dev6+deadbeef", nil},
	{"1.0a1+deadbeef", nil},
	{"1.0a1.post5+deadbeef", &Stage{"post", 5}},
	{"1.0a1.post5.dev6+deadbeef", &Stage{"post", 5}},
	{"1.0rc4+deadbeef", nil},
	{"1.0.post5+deadbeef", &Stage{"post", 5}},
	{"1!1.0+deadbeef", nil},
	{"1!1.0.dev6+deadbeef", nil},
	{"1!1.0a1+deadbeef", nil},
	{"1!1.0a1.post5+deadbeef", &Stage{"post", 5}},
	{"1!1.0a1.post5.dev6+deadbeef", &Stage{"post", 5}},
	{"1!1.0rc4+deadbeef", nil},
	{"1!1.0.post5+deadbeef", &Stage{"post", 5}},
}

func TestVersionPost(t *testing.T) {
	for _, post := range versionPosts {
		v, err := ParseVersion(post.version)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(v.Post(), post.post) {
			t.Errorf("version: %s, %v != %v", post.version, v.Post(), post.post)
		}
	}
}

var versionDevs = []struct {
	version string
	dev     *Stage
}{
	{"1.0", nil},
	{"1.0.dev0", &Stage{"dev", 0}},
	{"1.0.dev6", &Stage{"dev", 6}},
	{"1.0a1", nil},
	{"1.0a1.post5", nil},
	{"1.0a1.post5.dev6", &Stage{"dev", 6}},
	{"1.0rc4", nil},
	{"1.0.post5", nil},
	{"1!1.0", nil},
	{"1!1.0.dev6", &Stage{"dev", 6}},
	{"1!1.0a1", nil},
	{"1!1.0a1.post5", nil},
	{"1!1.0a1.post5.dev6", &Stage{"dev", 6}},
	{"1!1.0rc4", nil},
	{"1!1.0.post5", nil},
	{"1.0+deadbeef", nil},
	{"1.0.dev6+deadbeef", &Stage{"dev", 6}},
	{"1.0a1+deadbeef", nil},
	{"1.0a1.post5+deadbeef", nil},
	{"1.0a1.post5.dev6+deadbeef", &Stage{"dev", 6}},
	{"1.0rc4+deadbeef", nil},
	{"1.0.post5+deadbeef", nil},
	{"1!1.0+deadbeef", nil},
	{"1!1.0.dev6+deadbeef", &Stage{"dev", 6}},
	{"1!1.0a1+deadbeef", nil},
	{"1!1.0a1.post5+deadbeef", nil},
	{"1!1.0a1.post5.dev6+deadbeef", &Stage{"dev", 6}},
	{"1!1.0rc4+deadbeef", nil},
	{"1!1.0.post5+deadbeef", nil},
}

func TestVersionDev(t *testing.T) {
	for _, dev := range versionDevs {
		v, err := ParseVersion(dev.version)
		if err != nil {
			t.Error(err)
			continue
		}
		if !reflect.DeepEqual(v.Dev(), dev.dev) {
			t.Errorf("version: %s, %v != %v", dev.version, v.Dev(), dev.dev)
		}
	}
}
