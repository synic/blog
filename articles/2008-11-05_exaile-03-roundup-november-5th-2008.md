<!-- :metadata:

title: Exaile 0.3 Roundup for November 5th, 2008
tags: Exaile
published: 2008-11-05T23:18:35-0700
summary:

This week we've worked on the following:

-->

<p>This week we've worked on the following:</p>
 <p><ul>
 <li>Aren converted
the database save from Pickle to Shelve, in an attempt to make better use of
memory and to speed up loading and saving.  This changed the format of the
music.db, so you'll need to remove this and rescan your library if you're
following the weekly tarballs.  Just rm ~/.local/share/exaile/music.db</li>
<li>The addition of the "tagcovers" plugin.  This will search for album art
embedded in id3 tags of mp3 files.</li>
 <li>Album compilation support:  If
there is more than one artist in the same directory with the same album name,
it's marked as a compilation.  They will show up under "Various Artists" (as
well as just under the artist) in the collection panel, and will be treated as
one album where the album art manager is concerned.</li>
 </ul>
 </p>
<p>Have fun!</p>

<div class="restored-from-archive">
  <h3>Restored from VimTips archive</h3>
  <p>
  This article was restored from the VimTips archive. There's probably
  missing images and broken links (and even some flash references), but it
  was still important to me to bring them back.
  </p>
</div>
