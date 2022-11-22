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
	layout  reportlayout
}
type reportlayout struct {
	moduleWidth         int
	currentVersionWidth int
	newestVersionWidth  int
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

	replacements := make((map[string]string), len(modfile.Replace))

	for _, r := range modfile.Replace {
		_, e1 := replacements[r.Old.Path]
		if !e1 {
			replacements[r.Old.Path] = r.Old.Path
		}
	}

	for _, req := range modfile.Require {
		_, e2 := replacements[req.Mod.Path]
		if e2 {
			continue
		}

		rpe, err := dr.createEntry(ctx, req)
		if err != nil {
			return nil, err
		}
		report.refreshLayout(rpe)

		report.entries = append(report.entries, *rpe)
	}

	return &report, nil
}

func (report *DependencyReport) refreshLayout(rpe *ReportEntry) {
	if report.layout.moduleWidth < len(rpe.module) {
		report.layout.moduleWidth = len(rpe.module)
	}
	if report.layout.currentVersionWidth < len(rpe.currentVersion) {
		report.layout.currentVersionWidth = len(rpe.currentVersion)
	}
	if report.layout.newestVersionWidth < len(rpe.newestVersion) {
		report.layout.newestVersionWidth = len(rpe.newestVersion)
	}
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
	mw := report.layout.moduleWidth
	cvw := report.layout.currentVersionWidth
	nvw := report.layout.newestVersionWidth
	maxw := mw + cvw + nvw + 9 + 3
	sb.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", maxw)))
	sb.WriteString(fmt.Sprintf("|%*s|%*s|%*s|%9s|\n", mw, "module", cvw, "current", nvw, "newest", "actual"))
	sb.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", maxw)))
	for _, e := range report.entries {
		sb.WriteString(fmt.Sprintf("|%*s|%*s|%*s|%9v|\n", mw, e.module, cvw, e.currentVersion, nvw, e.newestVersion, e.actual))
	}
	sb.WriteString(fmt.Sprintf("+%s+\n", strings.Repeat("-", maxw)))

	fmt.Print(sb.String())
}
