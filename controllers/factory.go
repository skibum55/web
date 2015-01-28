package controllers

import (
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	bhimperrsrepo "github.com/cppforlife/bosh-hub/release/importerrsrepo"
	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhreltarsrepo "github.com/cppforlife/bosh-hub/release/releasetarsrepo"
	bhwatchersrepo "github.com/cppforlife/bosh-hub/release/watchersrepo"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type FactoryRepos interface {
	ReleasesRepo() bhrelsrepo.ReleasesRepository
	ReleaseTarsRepo() bhreltarsrepo.ReleaseTarballsRepository
	ReleaseVersionsRepo() bhrelsrepo.ReleaseVersionsRepository
	JobsRepo() bhjobsrepo.JobsRepository

	S3StemcellsRepo() bhstemsrepo.S3StemcellsRepository
	StemcellsRepo() bhstemsrepo.StemcellsRepository

	ImportsRepo() bhimpsrepo.ImportsRepository
	ImportErrsRepo() bhimperrsrepo.ImportErrsRepository
	WatchersRepo() bhwatchersrepo.WatchersRepository
}

type Factory struct {
	HomeController HomeController
	DocsController DocsController

	ReleasesController        ReleasesController
	ReleaseTarballsController ReleaseTarballsController

	StemcellsController StemcellsController

	JobsController     JobsController
	PackagesController PackagesController

	ReleaseWatchersController   ReleaseWatchersController
	ReleaseImportsController    ReleaseImportsController
	ReleaseImportErrsController ReleaseImportErrsController
}

func NewFactory(apiKey string, r FactoryRepos, runner boshsys.CmdRunner, logger boshlog.Logger) (Factory, error) {
	factory := Factory{
		HomeController: NewHomeController(r.ReleasesRepo(), r.StemcellsRepo(), logger),
		DocsController: NewDocsController(logger),

		ReleasesController: NewReleasesController(
			r.ReleasesRepo(),
			r.ReleaseVersionsRepo(),
			r.JobsRepo(),
			r.StemcellsRepo(),
			runner,
			logger,
		),

		ReleaseTarballsController: NewReleaseTarballsController(
			r.ReleasesRepo(),
			r.ReleaseTarsRepo(),
			logger,
		),

		StemcellsController: NewStemcellsController(r.StemcellsRepo(), logger),

		JobsController:     NewJobsController(r.ReleaseVersionsRepo(), r.JobsRepo(), logger),
		PackagesController: NewPackagesController(r.ReleaseVersionsRepo(), runner, logger),

		ReleaseWatchersController:   NewReleaseWatchersController(apiKey, r.WatchersRepo(), logger),
		ReleaseImportsController:    NewReleaseImportsController(r.ImportsRepo(), logger),
		ReleaseImportErrsController: NewReleaseImportErrsController(r.ImportErrsRepo(), logger),
	}

	return factory, nil
}