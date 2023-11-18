import pandas as pd
import sqlite3

conn = sqlite3.connect('Komica_Reply.db')
cursor = conn.cursor()

reply_df = pd.read_sql_query('SELECT * FROM "Reply_Data"', con=conn)
ignore_df = pd.read_sql_query('SELECT * FROM "Ignore_Data"', con=conn)
conn.close()

conn = sqlite3.connect('ReplyDatabase.db')
cursor = conn.cursor()
# create table
sql = 'CREATE TABLE "ignore_word" ( \
	"word"	TEXT NOT NULL UNIQUE \
);'
cursor.execute(sql)
sql = 'CREATE TABLE "reply_message" ( \
	"id"	INTEGER NOT NULL UNIQUE, \
	"message"	TEXT NOT NULL, \
	PRIMARY KEY("id" AUTOINCREMENT) \
);'
cursor.execute(sql)
sql = 'CREATE TABLE "word" ( \
	"word"	TEXT NOT NULL UNIQUE, \
	"target_id"	TEXT, \
	PRIMARY KEY("word") \
);'
cursor.execute(sql)
#ignore
for i in range(len(ignore_df)):
    word = ignore_df.loc[i]["Keyword"]
    sql = f'INSERT INTO "ignore_word" ("word") VALUES ("{word}");'
    cursor.execute(sql)
# reply
# keyword
for i in range(len(reply_df)):
    message = reply_df.loc[i]['Reply']
    sql = f'INSERT INTO "main"."reply_message" ("message") VALUES ("{message}");'
    cursor.execute(sql)
    df_keyword = reply_df.loc[i]['Keyword'].split(",")
    index = i + 1
    for keyword in df_keyword:
        result = pd.read_sql_query(f'SELECT * FROM "word" WHERE word = "{keyword}" LIMIT 1;', con=conn)
        isHasWord = len(result) > 0
        if isHasWord:
            oriTargetId = result.loc[0]["target_id"]
            oriTargetId += f',{index}'
            sql = f'UPDATE word SET target_id="{oriTargetId}" WHERE word="{keyword}";'
        else:
            sql = f'INSERT INTO word ("word","target_id") VALUES ("{keyword}","{index}");'
        print(sql)
        cursor.execute(sql)
conn.commit()
conn.close()