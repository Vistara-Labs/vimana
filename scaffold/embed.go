package scaffold

import (
	"embed"
	_ "embed"
)

//go:embed templates/grpc.gotmpl
var templateGrpcGoTemplate embed.FS

//go:embed templates/*.*tmpl
var createTemplateFS embed.FS

////go:embed templates/resource_scaffold.gotmpl
// var resourceScaffoldGoTemplate string

// //go:embed templates/data_source_scaffold.gotmpl
// var dataSourceScaffoldGoTemplate string

////go:embed templates/provider_scaffold.gotmpl
// var providerScaffoldGoTemplate string
