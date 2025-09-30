package construct

import (
	_ "embed"
)

// Vue frontend templates
//go:embed templates/frontend/index.vue
var vueIndexTemplate string

//go:embed templates/frontend/composable.ts
var vueComposableTemplate string

//go:embed templates/frontend/types.ts
var vueTypesTemplate string