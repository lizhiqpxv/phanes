package template

func init() {
	register(MappingTemplate, mapping)
}

var mapping = `
{{$true := true}}
{{$time := "time.Time"}}

package mapping

import (
	"{{.ProjectName}}/model"
	"{{.ProjectName}}/model/entity"
)


// {{.StructName}}sEntityToDto entity data transfer
func {{.StructName}}sEntityToDto({{.CamelName}}s []*entity.{{.StructName}}) []*{{.StructName}}Info {
	out := make([]*{{.StructName}}Info, 0, len({{.CamelName}}s))
	for _, c := range {{.CamelName}}s  {
		out = append(out, {{.StructName}}EntityToDto(c))
	}
	return out
}

// {{.StructName}}EntityToDto entity data transfer
func {{.StructName}}EntityToDto(e *entity.{{.StructName}}) *{{.StructName}}Info {
	return &{{.StructName}}Info{
		{{range $v :=.Fields}}
			{{.Name}}: {{if eq .Type $time}}e.{{.Name}}.Unix(),{{else}}e.{{.Name}},{{end}}
		{{end}}
	}
}
`
