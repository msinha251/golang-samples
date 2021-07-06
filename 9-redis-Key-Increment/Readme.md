This repo contains Fiber api in go with mongo as db and redis cache layer.

<br>
It contains below functions:
<br>

* `Get Article By Id` : This tries to get data from cache for given ID, if data is not there in cache then it brings it from mongo db and persist in redis cache for future use.

