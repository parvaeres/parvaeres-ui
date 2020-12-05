package models

type MAOReportResult struct {
	CheckId        string
	Category       string
	Severity       string
	Action         string
	Message        string
	Recommendation string
	Url            string
	Platform       string
	Resource       MAOResource
}

type MAOResource struct {
	Kind    string
	Version string
	Group   string
	Name    string
}

type MAOReportComponent struct {
	ResourceName      string
	ResourceNamespace string
	ResourceUID       string
	ResourceKind      string
	Results           []MAOReportResult
}

type MAOReports struct {
	OpsConformance      *MAOReportComponent `json:"Ops Conformance"`
	PodSecurity         *MAOReportComponent `json:"Pod Security"`
	WorkloadSupplyChain *MAOReportComponent `json:"Workload Software Supply Chain"`
}

type MAOResponse struct {
	Reports *MAOReports
}
