import os, sqlite3, datetime

backup_dir_location = '/backup/' + str(datetime.datetime.now())[:19]
if not os.path.exists(backup_dir_location): os.mkdir(backup_dir_location) 

database_dir_location = '/database'

database_path = [ 
    filename for filename in os.listdir(database_dir_location) 
    if os.path.isfile(os.join(database_dir_location, filename)) and filename.endswith('.db')
]

for database_location in database_path:
    db = sqlite3.connect(os.path.join(database_dir_location, database_location))
    f = open(os.path.join(backup_dir_location, database_location), 'w', encoding='utf8')
    f.write('\n'.join(db.iterdump()))
    f.close()
    db.close()