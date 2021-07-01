import pymongo

myclient = pymongo.MongoClient("mongodb://localhost:27017/")
mydb = myclient["go-crud"]
mycol = mydb["fiber-crud"]
for i in range(500002, 2000000):
    mydict = { "id": i, "title": f"title {i}", "content": f"content {i}"}
    mycol.insert_one(mydict)