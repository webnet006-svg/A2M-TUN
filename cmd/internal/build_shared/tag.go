package build_shared

import (
	"errors"
	"strings"

	"github.com/sagernet/sing-box/common/badversion"
	"github.com/sagernet/sing/common"
	F "github.com/sagernet/sing/common/format"
	"github.com/sagernet/sing/common/shell"
)

func latestTag() (string, error) {
	output, err := shell.Exec("git", "tag", "--list", "--sort=-version:refname", "v*").ReadOutput()
	if err != nil {
		return "", err
	}
	tags := strings.Fields(output)
	if len(tags) == 0 {
		return "", errors.New("no version tags found")
	}
	return tags[0], nil
}

func ReadTag() (string, error) {
	currentTag, err := shell.Exec("git", "describe", "--tags").ReadOutput()
	if err != nil {
		currentTag, err = latestTag()
		if err != nil {
			return currentTag, err
		}
		shortCommit, _ := shell.Exec("git", "rev-parse", "--short", "HEAD").ReadOutput()
		version := badversion.Parse(strings.TrimPrefix(currentTag, "v"))
		return version.String() + "-" + shortCommit, nil
	}
	currentTagRev, _ := shell.Exec("git", "describe", "--tags", "--abbrev=0").ReadOutput()
	if currentTagRev == currentTag {
		return currentTag[1:], nil
	}
	shortCommit, _ := shell.Exec("git", "rev-parse", "--short", "HEAD").ReadOutput()
	version := badversion.Parse(currentTagRev[1:])
	return version.String() + "-" + shortCommit, nil
}

func ReadTagVersionRev() (badversion.Version, error) {
	currentTagRev := common.Must1(shell.Exec("git", "describe", "--tags", "--abbrev=0").ReadOutput())
	return badversion.Parse(currentTagRev[1:]), nil
}

func ReadTagVersion() (badversion.Version, error) {
	currentTag, err := shell.Exec("git", "describe", "--tags").ReadOutput()
	if err != nil {
		currentTag, err = latestTag()
		if err != nil {
			return badversion.Version{}, err
		}
		version := badversion.Parse(strings.TrimPrefix(currentTag, "v"))
		version.Patch++
		return version, nil
	}
	currentTagRev := common.Must1(shell.Exec("git", "describe", "--tags", "--abbrev=0").ReadOutput())
	version := badversion.Parse(currentTagRev[1:])
	if currentTagRev != currentTag {
		if version.PreReleaseIdentifier == "" {
			version.Patch++
		}
	}
	return version, nil
}

func TestFlightVersion(version badversion.Version) string {
	return F.ToString(version.Major, ".", version.Minor, ".1000")
}
