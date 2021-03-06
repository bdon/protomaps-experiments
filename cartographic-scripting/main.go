package main

import "go.starlark.net/starlark"
import "fmt"
import "os"

type Feature struct {
	id int
	props *starlark.Dict
}

func NewFeature(id int, tags map[string]string) *Feature {
	props_dict := starlark.NewDict(1)
	for key, elem := range tags {
		props_dict.SetKey(starlark.String(key),starlark.String(elem))
	}
	return &Feature{id:id,props:props_dict}
}

func (f *Feature) Freeze() {

}

func (f *Feature) Hash() (uint32,error) {
	return 0, nil
}

func (f *Feature) String() string {
	return "abcd"
}

func (f *Feature) Truth() starlark.Bool {
	return true
}

func (f *Feature) Type() string {
	return "feature"
}

func (f *Feature) Attr(name string) (starlark.Value, error) {
	if (name == "props") {
		return f.props,nil
	}
	return nil,nil
}

func (f *Feature) AttrNames() []string {
	return []string{"props"}
}

type Layer struct {
	name string
	filter starlark.Callable
	tagDefinitions *starlark.List
}

type ComplexTag struct {
	key string
	value starlark.Value
}

func (t ComplexTag) Freeze() {

}

func (t ComplexTag) Hash() (uint32,error) {
	return 0,nil
}

func (t ComplexTag) String() string {
	return "complextag(key=" + t.key + ")"
}

func (t ComplexTag) Truth() starlark.Bool {
	return true
}

func (t ComplexTag) Type() string {
	return "complextag"
}

func main() {
	thread := &starlark.Thread{Name: "my thread"}

	var layers []Layer

	layerBuiltin := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var layer Layer
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "name", &layer.name, "tags?", &layer.tagDefinitions, "filter?", &layer.filter); err != nil {
			return nil, err
		}
		layers = append(layers,layer)
		return starlark.MakeInt(0), nil
	}

	tagBuiltin := func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var tag ComplexTag
		if err := starlark.UnpackArgs(b.Name(), args, kwargs, "key", &tag.key, "value", &tag.value); err != nil {
			return nil, err
		}
		return tag, nil
	}

	predeclared := starlark.StringDict{
		"layer":   starlark.NewBuiltin("layer", layerBuiltin),
		"tag":   starlark.NewBuiltin("tag", tagBuiltin),
	}

	globals, err := starlark.ExecFile(thread, "tileset_definition.star", nil, predeclared)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	attribution := globals["ATTRIBUTION"]
	fmt.Println(attribution)
	fmt.Println(attribution.Type())
	fmt.Println(layers)

	var features []*Feature
	features = append(features,NewFeature(1,map[string]string{
					"building":"yes",
					"height":"20",
			}))
	features = append(features,NewFeature(2,map[string]string{
					"highway":"motorway",
			}))
	features = append(features,NewFeature(3,map[string]string{
					"highway":"secondary",
			}))


	// for every feature in the input...
	for _, feature := range features {
		// feature_dict := starlark.NewDict(1)
		// feature_dict.SetKey(starlark.String("properties"),props_dict)

		// for every zoom level...
		for zoom := 0; zoom <= 4; zoom++ {

			for _, layer := range layers {
				v, err := starlark.Call(thread, layer.filter, starlark.Tuple{starlark.MakeInt(zoom),feature}, nil)
				if err != nil {
				} else if v.Truth() {
					fmt.Println("feature", feature.id, "appears in",layer.name,"at zoom", zoom)

					// execute tag definitions

					iter := layer.tagDefinitions.Iterate()
					defer iter.Done()
					var tagDefinition starlark.Value
					for iter.Next(&tagDefinition) {
						fmt.Println(tagDefinition)
						// check if it is a string
						basic_tag, is_string := tagDefinition.(starlark.String)
						if is_string {
							fmt.Println("Basic tag", basic_tag)
						}
						complex_tag, is_complextag := tagDefinition.(ComplexTag)
						if is_complextag {
							fmt.Println("Complex tag", complex_tag)	
							complex_tag_value := complex_tag.value
							literal_value, is_literal := complex_tag_value.(starlark.String)
							if is_literal {
								fmt.Println(literal_value)
							}
							fn_value, is_fn := complex_tag_value.(starlark.Callable)
							if is_fn {
								v, _ := starlark.Call(thread, fn_value, starlark.Tuple{starlark.MakeInt(zoom),feature}, nil)
								if err != nil {
									fmt.Println(v)
								}
							}
						}

					}

				}
			}


			// 	layer_dict, _ := layer.(*starlark.Dict)
			// 	name_val, _, _ := layer_dict.Get(starlark.String("name"))
			// 	filter_val, _, _ := layer_dict.Get(starlark.String("filter"))
			// 	filter_func, _ := filter_val.(*starlark.Function)
		}
	}
}
