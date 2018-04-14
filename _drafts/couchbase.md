---
layout: post
title: "Why we migrated from Couchbase to PostgreSQL"
---

**tl;dr:** After finding ourselves debugging our kernel's TCP/IP stack because
Couchbase didn't support joins, we decided to migrate to PostgreSQL and are
still happy with it 2+ years later. If you like epic DevOps stories, or even
better, preventing them, our team at Apple is also hiring: email-here

Picking a database can be a tricky problem. Ideally you'd like to pick the
"best tool for the job", but often times the reward for good work is more work,
and with more work comes new requirements that tend to invalidate earlier
assumptions.

Our team's decision to pick Couchbase goes back to before we even had a team
and there was just a single developer. The initial scope was to ETL some
manufacturing data from another system and to build an analytical dashboard for
it. Writing incremental map reduce code in JavaScript was a fun and fast way to
get the job done, and the project was a big success.

But there were deep problems with the system of record that the analytical
dashboard couldn't solve, so more people were hired and tasked to build a
replacement system. We gathered some initial requirements and determined we'd
only require modest functionality for the database and we are mostly worried
about data volumes and scaling, so we decided to double down on our investment
in Couchbase.

Fast forward a year, and the new system has been up and running for a while.
People are excited, and reward our efforts by giving us new use cases. One of
those use cases involved packing parts into boxes and turned out to require
transactions and JOINing tons of data. Unfortunately Couchbase 3 supported
neither of those features at the time.

We were already forced to implement application level JOINs for previous use
cases, and determined we could leverage Couchbase's CAS (Check and Set) feature
to fake the transactional guarantees we needed. We knew this was a bad idea,
but the idea of switching databases of what had become a critical system seemed
even less appealing. So the work was done, and the feature deployed, and
everybody was as happy as can be expected when forced to operate on themselves
using blunt instruments.

A few months later, in late 2015, our worst nightmare strikes. Our system goes
down in the middle of the night and a factory has come to a grinding halt as a
result. After some debugging we determine that the database is very slow, which
is causing requests to time out all over the place. However, we have no idea
why. There are no signs of exhausting CPU, disks or network, yet even trivial
key listing operations usually measured in milliseconds take over 6 seconds.
The logs are full of messages, but none seem to be new or relevant to what we
are seeing.

We decide to fail over to our secondary database cluster, and to our big
relief, things start to work again. Unfortunately the relief is only temporary.
A few hours later the problem reappears, and a new symptom appears: Some
database queries return incomplete results. At this point we're truly scared of
having lost data, but for a lack of better ideas attempt another failover, this
time back to the primary cluster. The system recovers yet again, including the
missing data, which is a relief. But the only thing we learned so far is that
the problem doesn't seem to be isolated to a certain machine or cluster. This
situation continues until we determine that a certain use case of the packing
feature seems to be responsible for triggering the issue, and we ask the
factory to stop using the feature while we investigate.

The post-mortem investigation turned out to be very hard. There are tons of log
files and messages to look at, but they all lead to dead ends. We know that the
packing feature can trigger a bit of an I/O storm by sending thousands of
concurrent requests to the database, but Couchbase shoud be able to handle
those, and we can't come up with any good theories that fully explain the
symptoms we are seeing.

So we decide we need more data, and the only way to get it is to be able to
reproduce the issue in an isolated environment. It takes a few failed attempts,
but eventually we manage to build a script that can set couchbase on fire and
trigger the issue on our secondary cluster by simulating some http requests
related to the packing feature.

<div style="text-align: center;">
<img width="150" src="/assets/posts/couchbase-to-postgres/couch-on-fire.gif">
</div>

This is fantastic as it allows us to capture much more data and to perform a
lot of experiments, but there are a lot of roads to go down, and a lot of them
turn into dead ends. But thanks to our sustained efforts, we finally get lucky
while analyzing tcpdump traffic between our cluster nodes in Wireshark. We
notice that 



# Todos

- This article switches from past to present tense in the middle ...
  is that a bad idea?
- Mention that we actually considered PostgreSQL from the beginning but decided
  against it because we were under time pressure.
