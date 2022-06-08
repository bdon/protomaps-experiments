NAME = "my very cool tileset"
ATTRIBUTION = "my attribution"

layer(
  name = "buildings",
  tags = [
    "name"
  ],
  filter = lambda z,f:f["properties"]["building"]
)

# layer definition for roads
# it might have some associated functions (must have unique names)
def roads_filter(z,f):
  if z < 2:
    return f["properties"]["highway"] == "motorway"
  else:
    return f["properties"]["highway"]

layer(
  name = "roads",
  tags = [
    "name",
    tag(
      key = "height",
      value = 1
    )
  ],
  filter = roads_filter
)
