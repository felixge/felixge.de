---
title: Let's Fix File Uploading
date: 2013-03-27T15:52:00+01:00
aliases:
  - "/2013/03/27/lets-fix-file-uploading.html"
---

**tl;dr:** We are creating an open source project to provide a turnkey file
uploading solution called [tus](http://www.tus.io/). Today we are happy to show
you our first [resumable file uploading demo](http://www.tus.io/demo.html).

It's 2013, and adding reliable file uploading to your app is still too damn
hard. And if content is king, this poses a serious problem to those interested
in acquiring it. In this post I'll describe the current uploading landscape and how we
are planning to fix it.

## Resumable Uploads

Have you ever tried uploading a big file over an unreliable internet
connection? Depending on the site you were using, the experience can be
extremely frustrating. Except for a few bigger sites, network errors will
force you to restart your upload from the beginning. Even worse, some sites
will be unable to detect this problem, leaving it up to you to figure out why
the progress bar is stuck.

The solution of course is to offer resumable file uploads, however, this is
easier said than done. While there is no shortage of bits and pieces that can
help you with building your own solution, there is no simple solution that you
can get up and running in a few minutes.

The first thing you'll need to figure out is a server side API. Google has
published several http protocols for resumable uploading (e.g.
[YouTube](https://developers.google.com/youtube/v3/guides/using_resumable_upload_protocol),
[Google Drive](https://developers.google.com/drive/manage-uploads)), but they
all rely on a non-standard http code `308 Resume Incomplete` which is on
collision course with [upcoming
standards](http://tools.ietf.org/html/draft-reschke-http-status-308-07). Amazon
S3 offers [multipart
uploads](http://docs.aws.amazon.com/AmazonS3/latest/dev/UsingRESTAPImpUpload.html)
(not to be confused with `multipart/form-data`), but they are hard to use and
the minimum chunk size is 5MB which is far too large to be useful under bad
networking conditions.

And even if you figure out the server side, you'll still need to take care of
the client side. For HTML5, your best bet is to start out with something
like [jQuery File Upload](http://blueimp.github.com/jQuery-File-Upload/) which
has some basic support for resumable uploads, but you'll still have a lot of
tweaking ahead of you in order to make it work with your server side API. If
you are also running native iOS / Android apps, you will be entirely on your
own to come up with a solution.

Fuck this. It's 2013, and you shouldn't have to waste several days for what
will only be a rudimentary solution, hard to scale, and a pain to maintain.
This is where we come in with [tus](http://tus.io/). Tus is an open source
project that aims to provide you with a turnkey solution to all your file
uploading needs. We are designing a
[protocol](http://www.tus.io/docs/http-protocol.html), building a
[server](https://github.com/tus/tusd), and are working on clients for HTML5,
iOS and Android. Once these are solid, we will also provide a hosted service,
making it even easier for people to get going, as well as allowing us to recoup
our investment in this project.

Today is a great day to have a first look, as we have just completed our first
[resumable file uploading demo](http://www.tus.io/demo.html) for HTML5. The
demo is pointed at a [tusd](http://github.com/tus/tusd) instance running on
EC2, and uses [jQuery File
Upload](http://blueimp.github.com/jQuery-File-Upload/) under the hood. What's
exciting about it, is that it not only handles network errors, but also browser
crashes and user errors like closing the page. No matter what happens, a
previous upload will always continue where it left off.

So check it out, and give us your
[feedback](https://github.com/tus/tus.io/issues/new?labels=feedback). We are
also very interested to hear what other issues you'd like to see solved. If you
are not sure about the possibilities, just read on.

## A world of possibilities

Resumable file uploads are just the tip of the iceberg, and there are many more
problems ripe for a solution. However, we need your help in prioritizing them,
so please let us know which of the things below you would find valuable.

* **File upload acceleration:** There are many situations where file uploads
  fail to achieve optimal performance due to latency and [TCP
  limitations](http://en.wikipedia.org/wiki/Bandwidth-delay_product). We hope
  to tackle this from multiple angles: Hosting a reverse CDN to improve latency
  regardless of user location, utilizing multiple TCP connections to increase
  the effective tcp window size, and maybe in the far future implement a
  specialized UDP protocol.
* **Huge files:** With HD cameras on the rise, uploads become bigger and bigger. A
  good solution should be able to handle files > 2GB with ease.
* **Streaming:** A lot of files can be processed as a stream, so tus should
  make it easy to let you download and process an active upload.
* **Service Integrations:** It should be easy to let your users select files
  from their Dropbox, Facebook, Google Drive or other cloud services.
* **Processing pipelines:** Most files are media files and require additional
  processing before being consumed by other users. Tus should easily integrate
  with existing processing solutions such as
  [Transloadit](http://transloadit.com/), [Zencoder](http://zencoder.com/) and
  others.
* **Scalability:** It should be easy to run a cluster of tus instances, allowing
  you to add or remove capacity as needed.
* **Security:** HTTPS is a no-brainer, but we are hoping to raise the bar by
  supporting file encryption on the fly for applications that can benefit from it.
* **Protocols:** HTTP is not the only protocol out there. We are interested in
  exploring uploads via FTP, E-Mail and other commonly used protocols.

If any of the above sounds exciting to you, please [let us
know](https://github.com/tus/tus.io/issues/new?labels=feedback). For now our
priority is to release a solid server along with great client libraries for
HTML5, iOS and Android, but at the end of the day it will be your feedback
that will drive our roadmap.
