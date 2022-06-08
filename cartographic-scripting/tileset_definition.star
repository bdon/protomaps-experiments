NAME = "my very cool tileset"
ATTRIBUTION = "my attribution"

LAYERS = []

# layer definition for buildings
LAYERS.append({
  "name":"buildings",
  "tags":[
    {
      "key":"height",
      "value":1
    },
    {
      "key":"name",
      "value":lambda z,f:f["name"]
    }
  ],
  "filter":lambda z,f:f["properties"]["building"]
})

# layer definition for roads
# it might have some associated functions (must have unique names)
def roads_filter(z,f):
  if z < 2:
    return f["properties"]["highway"] == "motorway"
  else:
    return f["properties"]["highway"]

LAYERS.append({
  "name":"roads",
  "tags":[
    {
      "key":"height",
      "value":1
    },
    {
      "key":"name",
      "value":lambda z,f:f["name"]
    }
  ],
  "filter":roads_filter
})

