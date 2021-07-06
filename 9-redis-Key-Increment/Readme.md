This repo contains Go module for monitoring Redis status and backup of redis keys in mongo. 
<br>
This will scan redis service as soon as it find redis is down go pipleine will increase the redis keys in mongo siteid collection and then it'll will be restored in redis once Redis will be UP again.

<br>
Go pipeline contains below functions:
<br>

* Function 1: It scans for redis keys and make a copy of that in mongo siteid collection. 
	
* Function 2: As soon as it founds redis == DOWN and siteidIncreased == False (default), it'll increase the Mongo siteids with 500 and update siteidIncreased bool value to True, which is also stored and updated in another mongodb collection siteIdIncreased.
	
* Function 3: As soon as it founds redis == UP and siteidIncreased == True, then it'll update the Redis key's from mongo siteid collection and update the siteidIncreased to False.


![Scheme](/9-redis-Key-Increment/images/Redis_SiteID_Update_with_Mongo_1.png){width=60%}