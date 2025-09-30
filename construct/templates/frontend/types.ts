export interface {{.ResourceName}} {
  id: number
  {{range .Fields}}{{.Name}}: {{.TypeScriptType}}
  {{end}}created_at: string
  updated_at: string
}

export interface {{.ResourceName}}CreateRequest {
  {{range .Fields}}{{.Name}}: {{.TypeScriptType}}
  {{end}}
}

export interface {{.ResourceName}}UpdateRequest {
  {{range .Fields}}{{.Name}}?: {{.TypeScriptType}}
  {{end}}
}