---
layout: post
title: "PostgreSQL operations that you can't EXPLAIN"
redirect_from:
  - /2017/05/05/postgresql-operations-that-you-cant-explain.html
date: 2017-05-23T16:25:00+02:00
updated: 2017-05-23T16:25:00+02:00
---

I am currently dabbling in PostgreSQL internals in my spare time, including
trying to understand the query planner and executor. It's a deeply humbling
experience, but occasionally I'm delighted to be able to answer simple
questions I'm researching.

One of these questions is how PostgreSQL is executing `ORDER BY` clauses inside
of aggregate expressions.

Consider a simple table that holds 3 rows of an int column ranging from 1 to 3:

```
# CREATE TABLE foo AS
SELECT * FROM generate_series(1, 3) i;
SELECT 3
# SELECT * FROM foo;
 i
---
 1
 2
 3
(3 rows)
```

Now let's assume we want to convert those rows into a json array. This can
be accomplished using the [json\_agg](https://www.postgresql.org/docs/9.6/static/functions-aggregate.html#FUNCTIONS-AGGREGATE-TABLE) 
aggregate function:

```
# SELECT json_agg(i) FROM foo;
 json_agg
-----------
 [1, 2, 3]
(1 row)
```

But what if we want to order the array elements in descending order? No
problem, we can simply use the order-by-clause support available to all
[aggregate expressions](https://www.postgresql.org/docs/9.6/static/sql-expressions.html#SYNTAX-AGGREGATES):

```
# SELECT json_agg(i ORDER BY i DESC) FROM foo;
 json_agg
-----------
 [3, 2, 1]
(1 row)
```

If you're like me, you'd probably imagine the `EXPLAIN` output for this query
to contain three nodes: a `Seq Scan`, a `Sort`, and an `Aggregate`. However,
when we run `EXPLAIN`, it turns out there is no `Sort` node. 

```
# EXPLAIN SELECT json_agg(i ORDER BY i DESC) FROM foo;
                         QUERY PLAN
-------------------------------------------------------------
 Aggregate  (cost=41.88..41.89 rows=1 width=32)
   ->  Seq Scan on foo  (cost=0.00..35.50 rows=2550 width=4)
(2 rows)
```

In fact, the plan even remains unchanged even if we change the sort order or
remove the clause entirely:

```
# EXPLAIN SELECT json_agg(i ORDER BY i ASC) FROM foo;
                         QUERY PLAN
-------------------------------------------------------------
 Aggregate  (cost=41.88..41.89 rows=1 width=32)
   ->  Seq Scan on foo  (cost=0.00..35.50 rows=2550 width=4)
(2 rows)

# EXPLAIN SELECT json_agg(i) FROM foo;
                         QUERY PLAN
-------------------------------------------------------------
 Aggregate  (cost=41.88..41.89 rows=1 width=32)
   ->  Seq Scan on foo  (cost=0.00..35.50 rows=2550 width=4)
(2 rows)
```

But how is this possible? If there is no `Sort` node, how does the aggregate
perform its sorting?

Well, it turns out that aggregate plan nodes can [perform their own
sorting](https://github.com/postgres/postgres/blob/aa3bcba08d466bc6fd2558f8f0bf0e6d6c89b58b/src/backend/executor/nodeAgg.c#L520-L541).
These sorts do not explicitly show up in `EXPLAIN`, but we can trace them using
the
[trace\_sort](https://www.postgresql.org/docs/9.6/static/runtime-config-developer.html#GUC-TRACE-SORT)
developer option:

```
# SET trace_sort=true;
SET
# SET client_min_messages='log';
SET
# SELECT json_agg(i ORDER BY i DESC) FROM foo;
LOG:  begin datum sort: workMem = 4096, randomAccess = f
LOG:  performsort starting: CPU 0.00s/0.00u sec elapsed 0.00 sec
LOG:  performsort done: CPU 0.00s/0.00u sec elapsed 0.00 sec
LOG:  internal sort ended, 25 KB used: CPU 0.00s/0.00u sec elapsed 0.00 sec
 json_agg
-----------
 [3, 2, 1]
(1 row)
```

Another way to see that the sort is part of the plan is to use
[debug\_print\_plan](https://www.postgresql.org/docs/9.6/static/runtime-config-logging.html#RUNTIME-CONFIG-LOGGING-WHAT)
option, but the output is not for the faint of heart. Specifically you need to
look for the `:aggorder` field below:

```
# SET debug_print_plan=true;
SET
# SELECT json_agg(i ORDER BY i DESC) FROM foo;
LOG:  plan:
DETAIL:     {PLANNEDSTMT
   :commandType 1
   :queryId 0
   :hasReturning false
   :hasModifyingCTE false
   :canSetTag true
   :transientPlan false
   :dependsOnRole false
   :parallelModeNeeded false
   :planTree
      {AGG
      :startup_cost 41.88
      :total_cost 41.89
      :plan_rows 1
      :plan_width 32
      :parallel_aware false
      :plan_node_id 0
      :targetlist (
         {TARGETENTRY
         :expr
            {AGGREF
            :aggfnoid 3175
            :aggtype 114
            :aggcollid 0
            :inputcollid 0
            :aggtranstype 2281
            :aggargtypes (o 23)
            :aggdirectargs <>
            :args (
               {TARGETENTRY
               :expr
                  {VAR
                  :varno 65001
                  :varattno 1
                  :vartype 23
                  :vartypmod -1
                  :varcollid 0
                  :varlevelsup 0
                  :varnoold 1
                  :varoattno 1
                  :location 16
                  }
               :resno 1
               :resname <>
               :ressortgroupref 1
               :resorigtbl 0
               :resorigcol 0
               :resjunk false
               }
            )
            :aggorder (
               {SORTGROUPCLAUSE
               :tleSortGroupRef 1
               :eqop 96
               :sortop 521
               :nulls_first true
               :hashable true
               }
            )
            :aggdistinct <>
            :aggfilter <>
            :aggstar false
            :aggvariadic false
            :aggkind n
            :agglevelsup 0
            :aggsplit 0
            :location 7
            }
         :resno 1
         :resname json_agg
         :ressortgroupref 0
         :resorigtbl 0
         :resorigcol 0
         :resjunk false
         }
      )
      :qual <>
      :lefttree
         {SEQSCAN
         :startup_cost 0.00
         :total_cost 35.50
         :plan_rows 2550
         :plan_width 4
         :parallel_aware false
         :plan_node_id 1
         :targetlist (
            {TARGETENTRY
            :expr
               {VAR
               :varno 1
               :varattno 1
               :vartype 23
               :vartypmod -1
               :varcollid 0
               :varlevelsup 0
               :varnoold 1
               :varoattno 1
               :location -1
               }
            :resno 1
            :resname <>
            :ressortgroupref 0
            :resorigtbl 0
            :resorigcol 0
            :resjunk false
            }
         )
         :qual <>
         :lefttree <>
         :righttree <>
         :initPlan <>
         :extParam (b)
         :allParam (b)
         :scanrelid 1
         }
      :righttree <>
      :initPlan <>
      :extParam (b)
      :allParam (b)
      :aggstrategy 0
      :aggsplit 0
      :numCols 0
      :grpColIdx
      :grpOperators
      :numGroups 1
      :aggParams (b)
      :groupingSets <>
      :chain <>
      }
   :rtable (
      {RTE
      :alias <>
      :eref
         {ALIAS
         :aliasname foo
         :colnames ("i")
         }
      :rtekind 0
      :relid 404407
      :relkind r
      :tablesample <>
      :lateral false
      :inh false
      :inFromCl true
      :requiredPerms 2
      :checkAsUser 0
      :selectedCols (b 9)
      :insertedCols (b)
      :updatedCols (b)
      :securityQuals <>
      }
   )
   :resultRelations <>
   :utilityStmt <>
   :subplans <>
   :rewindPlanIDs (b)
   :rowMarks <>
   :relationOids (o 404407)
   :invalItems <>
   :nParamExec 0
   }

LOG:  begin datum sort: workMem = 4096, randomAccess = f
LOG:  performsort starting: CPU 0.00s/0.00u sec elapsed 0.00 sec
LOG:  performsort done: CPU 0.00s/0.00u sec elapsed 0.00 sec
LOG:  internal sort ended, 25 KB used: CPU 0.00s/0.00u sec elapsed 0.00 sec
 json_agg
-----------
 [3, 2, 1]
(1 row)
```

In conlusion: While `EXPLAIN` is a powerful weapon for those seeking to
understand query performance, you should be aware of the fact that there are
things that you can't `EXPLAIN` :).

The research for this article was done using PostgreSQL 9.6.3 and LLDB, but
should also apply to recent versions as well as the upcoming PostgreSQL 10.

```
# SELECT version();
                                                   version
--------------------------------------------------------------------------------------------------------------
 PostgreSQL 9.6.3 on x86_64-apple-darwin14.5.0, compiled by Apple LLVM version 7.0.0 (clang-700.1.76), 64-bit
(1 row)
```
