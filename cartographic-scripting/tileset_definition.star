NAME = "my very cool tileset"
ATTRIBUTION = "my attribution"

def height_from_tag(z,f):
  return f.props["height"]

layer(
  name = "buildings",
  tags = [
    "name", # this is passed straight through as a string value
    tag(
      key = "height",
      value = height_from_tag
    )
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
      key = "roadclass",
      value = "major"
    )
  ],
  filter = roads_filter
)
