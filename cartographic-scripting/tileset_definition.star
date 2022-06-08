NAME = "my very cool tileset"
ATTRIBUTION = "my attribution"

layer(
  name = "buildings",
  tags = [
    "name"
  ],
  filter = lambda z,f:f.props["building"]
)

# layer definition for roads
# it might have some associated functions (must have unique names)
def roads_filter(z,f):
  if z < 2:
    return f.props["highway"] == "motorway"
  else:
    return f.props["highway"]

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
