tippecanoe \
--force \
-zg \
--projection=EPSG:4326 \
--no-tile-compression \
--no-feature-limit --no-tile-size-limit \
-o cb_2018_us_zcta510_500k_nolimit.mbtiles -l zcta cb_2018_us_zcta510_500k.json

tippecanoe \
--force \
-zg \
--projection=EPSG:4326 \
--no-tile-compression \
--coalesce-densest-as-needed \
-o cb_2018_us_zcta510_500k_coalesced.mbtiles -l zcta cb_2018_us_zcta510_500k.json