ogr2ogr -f GeoJSON ne_10m_rivers_lake_centerlines.geojson ne_10m_rivers_lake_centerlines/ne_10m_rivers_lake_centerlines.shp
ogr2ogr -f GeoJSON ne_10m_rivers_north_america.geojson ne_10m_rivers_north_america/ne_10m_rivers_north_america.shp
tippecanoe -zg -o ne_10m_rivers.mbtiles -l rivers --coalesce-densest-as-needed --extend-zooms-if-still-dropping ne_10m_rivers_lake_centerlines.geojson ne_10m_rivers_north_america.geojson
pmtiles-convert ne_10m_rivers.mbtiles ne_10m_rivers.pmtiles