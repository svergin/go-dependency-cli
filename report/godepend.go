package report

import (
	"context"
	"errors"
	"go-dependency-cli/git"
	goproxy "go-dependency-cli/proxyclient"
	"io"

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
	currentVersion string
	newestVersion  string
	actual         bool
}

func (dr *DependencyResolver) ErstelleReport(ctx context.Context, repoURL string, optBranch string) (*DependencyReport, error) {
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
		var rpe ReportEntry
		rpe.currentVersion = req.Mod.Version
		info, err := dr.gpc.GetLatest(ctx, req.Mod.Path)
		if err != nil {
			return nil, err
		}
		rpe.newestVersion = info.Version
		if !semver.IsValid(rpe.currentVersion) || !semver.IsValid(rpe.newestVersion) {
			return nil, errors.New("version not valid")
		}
		c := semver.Compare(rpe.currentVersion, rpe.newestVersion)
		rpe.actual = c == 0
		report.entries = append(report.entries, rpe)
	}

	return &report, nil
}
