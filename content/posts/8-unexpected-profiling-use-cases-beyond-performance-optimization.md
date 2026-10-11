---
title: "9 Unexpected Profiling Use Cases Beyond Performance Optimization"
date: 2023-04-03T22:17:31.000Z
slug: "8-unexpected-profiling-use-cases-beyond-performance-optimization"
aliases:
  - /8-unexpected-profiling-use-cases-beyond-performance-optimization/
---

Profiling is often associated with performance optimization, but it has many other valuable use cases. Here are 9 of my favorites:

1.  **Detect Workload Shifts**: Handler frames get wider when they get more work.
2.  **Detect Data Nesting**: Frames often get deeper when processing more nested data.
3.  **Detect Dead Code**: If a profile shows a function, you can reject the hypothesis that the code is unused. This can also help with prioritizing security fixes for new vulnerabilities.
4.  **Debug Unresponsiveness**: I once used profiling to diagnose a PostgreSQL vacuum process stuck in an infinite loop due to index corruption. Loops with O(N^2) scaling issues also usually stand out in profiles.
5.  **Identify Configuration Issues**: Profiling can be used to identify bottlenecks caused by misconfiguration. For example spending a lot of CPU time in read/write syscalls is often caused by a lack of buffering. Similarly, a high amount of CPU time in TLS can indicate [connection churn problems](https://dev.to/gkampitakis/http-connection-churn-in-go-34pl).
6.  **Debug Errors**: If your error rates get high, your profile will show you the error path in the code.
7.  **Verify Deployments**: New frames appear, old frames disappear, regressions become visible.
8.  **Improve Maintainability**: The call hierarchies in your profiles can reveal deficiencies in the structure of your application.
9.  **Faster Onboarding**: Studying the hot code paths of an application is a good way to ramp up quickly.

Profiling may not be the ideal tool for all of these use cases, but no observability signal can give you such a broad visibility. Logging, tracing and metrics all depend on having the right instrumentation at the right time.

Have you used profiling for any other unexpected use cases? Please share them in the comments, [twitter](https://twitter.com/felixge/status/1643015560588845057) or [mastodon](https://home.social/@felixge/110137226772072415).

*Disclaimer: I work on [*profiling for Datadog*](https://www.datadoghq.com/product/code-profiling/), but I'm not barking on behalf of my employer here. I've been meaning to write this down for a while now since it often comes up in conversations.*
