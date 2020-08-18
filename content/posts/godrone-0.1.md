---
title: "GoDrone - A Parrot AR Drone 2.0 Firmware written in Go"
date: 2013-12-25T20:24:00+01:00
---

Merry Christmas (or [Newtonmas][1] if you prefer) everybody.

Today I'm very happy to release the first version of my favorite side project, GoDrone.

GoDrone is a free software alternative firmware for the [Parrot AR Drone
2.0][3]. And yes, this hopefully makes it the first robotic visualizer for Go's
[garbage collector][2] : ).

At this point the firmware is good enough to fly and provide basic attitude stabilization (using a simple complementary filter + pid controllers), so I'd really love to get feedback from any adventurous AR Drone owners. I'm providing binary installers for OSX/Linux/Windows:

[http://www.godrone.io/en/latest/index.html](http://www.godrone.io/en/latest/index.html)

But you may also choose to install from [source](http://www.godrone.io/en/latest/contributor/install_from_source.html).

Depending on initial feedback, I'd love to turn GoDrone into a viable alternative to the official firmware, and breathe some fresh air into the development of robotics software. In particular I'd like to show that web technologies can rival native mobile/desktop apps in costs and UX for providing user interfaces to robots, and I'd also like to promote the idea of using high level languages for firmware development in linux powered robots.

If you're interested, please make sure to join the mailing list / come and say hello in IRC:

[http://www.godrone.io/en/latest/user/community_support.html](http://www.godrone.io/en/latest/user/community_support.html)

[1]: http://www.youtube.com/watch?v=EqiiCOFR0Y8
[2]: http://www.godrone.io/en/latest/user/faq.html#isn-t-go-unsuitable-for-real-time-applications-like-this
[3]: http://en.wikipedia.org/wiki/Parrot_AR.Drone#Version_2.0
