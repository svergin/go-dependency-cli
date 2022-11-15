package report

import (
	"context"
	"errors"
	"fmt"
	"go-dependency-cli/git"
	goproxy "go-dependency-cli/proxyclient"
	"io"
	"strings"

	"github.com/rogpeppe/go-internal/modfile"
	"github.com/rogpeppe/go-internal/semver"
)

//TODO: https://dev.azure.com/nortal-de-dev-training/_git/learning-go?path=/doc/1.3.md&version=GBmain&_a=preview

type DependencyResolver struct {
	gpc goproxy.Client
}
type DependencyReport struct {
	entries []ReportEntry
}

type ReportEntry struct {
	module         string
	currentVersion string
	newestVersion  string
	actual         bool
}

func (dr *DependencyResolver) CreateReport(ctx context.Context, repoURL string, optBranch string) (*DependencyReport, error) {
	report := DependencyReport{}

	filesystem, err := git.Clone(ctx, repoURL, optBranch)
	if err != nil {
		return nil, err
	}
	gm, err := filesystem.Open("go.mod")
	if err != nil {
		return nil, err
	}
	fb, err := io.ReadAll(gm)
	if err != nil {
		return nil, err
	}
	modfile, err := modfile.Parse("go.mod", fb, nil)
	if err != nil {
		return nil, err
	}

	for _, req := range modfile.Require {
		rpe, err := dr.createEntry(ctx, req)
		if err != nil {
			return nil, err
		}
		report.entries = append(report.entries, *rpe)
	}

	return &report, nil
}

func (dr *DependencyResolver) createEntry(ctx context.Context, req *modfile.Require) (*ReportEntry, error) {
	info, err := dr.gpc.GetLatest(ctx, req.Mod.Path)
	if err != nil {
		return nil, err
	}
	if !semver.IsValid(req.Mod.Version) || !semver.IsValid(info.Version) {
		return nil, errors.New("version not valid")
	}
	c := semver.Compare(req.Mod.Version, info.Version)
	rpe := ReportEntry{
		module:         req.Mod.Path,
		currentVersion: req.Mod.Version,
		newestVersion:  info.Version,
		actual:         c == 0,
	}
	return &rpe, nil
}

func (report *DependencyReport) Print() {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", 68)))
	sb.WriteString(fmt.Sprintf("|%40s|%8s|%8s|%9s|\n", "module", "current", "newest", "actual"))
	sb.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", 68)))
	for _, e := range report.entries {
		sb.WriteString(fmt.Sprintf("|%40s|%8s|%8s|%9v|\n", e.module, e.currentVersion, e.newestVersion, e.actual))
	}
	sb.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", 68)))

	fmt.Print(sb.String())
}
