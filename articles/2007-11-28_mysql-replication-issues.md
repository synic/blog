<!-- :metadata:

title: MySQL Replication Issues
tags: Programming, Linux
publishedAt: 2007-11-28T16:25:16-07:00
summary:

At <a href='http://www.sendoutcards.com'>SendOutCards</a>, we use <a
href='http://dev.mysql.com/doc/refman/5.0/en/replication.html'>MySQL
replication</a> to ensure that if our main database server goes down because of
hardware failure, we'll still have an server that is up to date with our data.
In a nutshell, our main database server sends another server every update that
is performed on itself...

-->

At <a href='http://www.sendoutcards.com'>SendOutCards</a>, we use <a
href='http://dev.mysql.com/doc/refman/5.0/en/replication.html'>MySQL
replication</a> to ensure that if our main database server goes down because of
hardware failure, we'll still have an server that is up to date with our data.
In a nutshell, our main database server sends another server every update that
is performed on itself.  <br><br>
 For non-hardware issues (IE, a bad SQL
statement that wipes out everyone's account balance), we've written a script on
our slave server that creates a complete backup of the database every night.
<br><br>
 This system is great, <i>in theory</i>.  Unfortunately, MySQL
replication seems to be highly problematic, at least with our database schema
and hardware setup.  Before long, the slave server always becomes corrupt after
an arbitrary length of time (if you have any idea of why this might be
happening, please contact me).  Due to the size of our database, we've had to
change the <code>expire_log_days</code> (the number of days before the binary
log is deleted) in <code>/etc/mysql/my.cnf</code> to 2 days; the reason is that
if the slave becomes corrupt and is not repaired immediately the log
information can become out of date and a repair without stopping the master
database is impossible.<br><br>To make sure we know about the situation as
soon as possible, I've written a simple python script to monitor the slave
status via cron that will send an email once an error is detected.

```python
#!/usr/bin/env python

import MySQLdb, MySQLdb.cursors, os

## REPLACE THE FOLLOWING WITH YOUR INFORMATION
db_user = 'someuser'
db_pass = 'somepassword'
email = 'administrators@email.tld'

db = MySQLdb.connect(user=db_user, passwd=db_pass, db='mysql',
        cursorclass=MySQLdb.cursors.DictCursor)
cur = db.cursor()
cur.execute('SHOW SLAVE STATUS')
row = cur.fetchone()

errno = row['Last_Errno']
error = row['Last_Error']

if errno > 0:
        os.system('echo \"%s\" | mail -s \"Slave error on `hostname`\" %s' %
                (error, email))
```

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
