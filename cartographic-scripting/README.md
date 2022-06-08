## Cartographic Scripting

There are two kinds of scripting:

* OpenStreetMap Data > tagged vector features in tiles - "tileset scripting". Requires detailed cleanup on messy OSM data, and filtering to cut down on tile size.
* tagged vector feature in tiles > visual shapes and labels - "client scripting". Requires expressiveness for flexible display.

## Prototype

This folder has an [example of a tileset](./tileset_definition.star) scripting design written in the [Starlark](https://github.com/bazelbuild/starlark) language, a subset of Python. 

The program [main.go](./main.go) will execute `tileset_definition.star` against dummy data. Control which features are assigned to which layers by defining functions in the file:

```py
def roads_filter(z,f):
  if z < 2:
    return f["properties"]["highway"] == "motorway"
  else:
    return f["properties"]["highway"]
```

## Other Designs

### Mapbox GL Expressions

[osm-bright-gl-style](https://github.com/openmaptiles/osm-bright-gl-style/blob/master/style.json)

Client scripting. Functional Lisp-inside-JSON. Must run on a unique environment (GL), not turing-complete; limited to built-in functions, can require preprocessor.

### Tangram Styles

[refill-style.yaml](https://github.com/tangrams/refill-style/blob/gh-pages/refill-style.yaml)

Client scripting. Runs in WebGL. Variables, inline shader code, embeddable JavaScript expressions.

### Protomaps JS Rules

(disclaimer: i designed this)

[default style.ts](https://github.com/protomaps/protomaps.js/blob/master/src/default_style/style.ts)

Client scripting. Everything is TypeScript, turing complete, variables and parameterization, runs 100% on CPU.

### osm2pgsql configuration files

[default.style for osm2pgsql](https://github.com/openstreetmap/osm2pgsql/blob/master/default.style)

Declarative, custom textual format for creating PostgreSQL tables. geometrytype-based selection, key wildcards.

### imposm3 mappings

[example-mapping.yml](https://github.com/omniscale/imposm3/blob/master/example-mapping.yml)

Declarative YAML for creating PostgreSQL tables. 

### osmium-tool configuration

[osmium-export docs](https://docs.osmcode.org/osmium/latest/osmium-export.html)

Declarative JSON format for Simple Features like GeoJSON from OSM.

### HOTOSM Export Tool YAML

(disclaimer: i designed this)

[default.yml](https://github.com/hotosm/osm-export-tool-python/blob/master/osm_export_tool/mappings/default.yml)

Declarative YAML specification for transforming OSM to Simple Features. Can embed subset of SQL inside `where` expressions for layer definitions:

    where:
      - amenity IN ('university','school','library','fuel','hospital','fire_station','police','townhall')

### Tilemaker LUA

[process-openmaptiles.lua](https://github.com/systemed/tilemaker/blob/master/resources/process-openmaptiles.lua)

True tileset scripting. Write Lua functions that are executed against OSM features to define tileset features. Depends on Lua C runtime. 

### Planetiler CustomMap YAML

[power.yml](https://github.com/onthegomap/planetiler/blob/main/planetiler-custommap/src/main/resources/samples/power.yml)

Declarative YAML for transforming OSM into tileset features.