<!-- :metadata:

title: Decorator fun in Python
tags: Programming, Python
published: 2007-03-26T18:25:00-0700
summary:

I came across the following code a while ago.  I can't take credit for it, and
I can't remember where I got it.  Oh well.  It's pretty cool nonetheless.

-->

I came across the following code a while ago.  I can't take credit for it, and
I can't remember where I got it.  Oh well.  It's pretty cool nonetheless.

```python
def threaded(f):
    """
        A decorator that will make any function run in a new thread
    """
    def wrapper(*args):
        t = threading.Thread(target=f, args=args)
        t.setDaemon(True)
        t.start()

    wrapper.name = func.name
    wrapper.dict = func.dict
    wrapper.doc = func.doc

    return wrapper

def synchronized(func):
    """
        A decorator to make a function synchronized – which means only one
        thread is allowed to access it at a time
    """
    def wrapper(self,__args,*__kw):
        try:
            rlock = self._sync_lock
        except AttributeError:
            from threading import RLock
            rlock = self.dict.setdefault(‘_sync_lock‘,RLock())
        rlock.acquire()
        try:
            return func(self,__args,*__kw)
        finally:
            rlock.release()
    wrapper.name = func.name
    wrapper.dict = func.dict
    wrapper.doc = func.doc
    return wrapper

# example:

@threaded
def this_is_a_long_running_function():
    connect_to_network_and_do_a_lot_of_stuff
# any time the above function is called, it will run in a new thread
```

<div class='vimtip'>
<h3>
<b>vim tip:</b> <i>A,  D, and C</i>
</h3>

<p>
<b>A</b> will take you to the end of the line and put you in insert mode.
<b>D</b> will delete until the end of the line, and <b>C</b> will delete to the
end of the line and put you in insert mode.
</p>
</div>

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
