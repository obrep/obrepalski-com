+++
title = 'Cost Observability and Guardrails'
date = 2024-02-18
tags = ["series", "docs"]
series =  ["Optimising Cloud Services At Scale"]
series_order =  2
+++

I've worked on large-scale Recommender Systems (RecSys) for the past few years—first at ByteDance and now at ShareChat, where we serve over 325 million monthly active users, delivering tens of thousands of personalized feeds every second. At this scale, every inefficiency adds up quickly, and optimising services and workloads becomes essential—not just for reducing costs, but also for overall system maintainability.

My recent focus has been on running our workloads as efficiently as possible, spanning everything from low-level performance tweaks to leading division-wide initiatives. While some projects seemed minor in isolation, their cumulative impact has been massive. **In 2024 alone, our team's costs reduced by 90%, with the entire RecSys costs dropping by over 75%**. The best part? At the same time user retention from our RecSys improved, proving that cost optimisation doesn't have to compromise quality.

> :memo: **Want to know more?** Check out our [official blog](https://sharechat.com/blogs/artificial-intelligence/building-world-class-recommendations-at-lower-costs)!

Over this time I've learned a lot about cost optimisations - what works effectively and what leads to failure. Successful large-scale efforts depend on getting the fundamentals right, with observability and guardrails being among the most critical elements. Observability helps teams understand their spending and measure optimisation impact, while guardrails prevent unintended regressions.

---

# Cost observability

As with all engineering problems, you can't optimise what you don't see. This is why observability is the first step toward enabling teams to run their workloads efficiently. Without proper visibility, cost-cutting efforts often feel like guesswork. A well-structured approach ensures that teams have the right tools to track and understand their spending, can measure the impact of changes, and set up guardrails to prevent regressions.

Understanding what you're actually paying for in the cloud is often more difficult than it seems. Costs don't always correspond directly to the services you interact with daily. Many inefficiencies stem from over-provisioned resources, managed service markups, and unexpected network transfer fees. While major cloud providers offer cost management tools (like [GCP's Cost Management](https://cloud.google.com/cost-management)), these alone often aren't enough for a complete picture, especially if you're using multiple clouds or external services.

**Clear Resource Ownership** is fundamental to cost management. Mapping resources to owners can reveal important signals - difficulties in identifying owners often indicate that resources are no longer needed or lack proper maintenance. For resources still in use, clear ownership is crucial not just for cost management but for reliability—trying to track down resource owners during an outage is a nightmare scenario you want to avoid.

**Comprehensive Resource Labeling** adds crucial context to cost data and enhances the utility of cost visualization tools. [Labels](https://cloud.google.com/resource-manager/docs/labels-overview) should be baked into your infrastructure processes, ideally enforced through Infrastructure as Code (e.g. Terraform). Based on my experience, the most useful dimensions include Team, Business Unit, Environment (staging/production), and Technology/Component (e.g. BigTable/Scylla).

**Multi-Source Cost Aggregation** is essential as infrastructure bills typically extend beyond products from a single provider. Organizations often use external APIs, deal with licensing costs, or pay for third-party SaaS platforms. Having a unified view across all these sources enables better decision-making and optimisation opportunities.

Finally, **engineers must have ownership over their resource spending and the tools to understand it**. This visibility and accountability should be considered as fundamental as access to performance metrics or logging, enabling teams to make informed decisions about infrastructure choices and identify optimisation opportunities.

> :memo: A few useful tools:
> - [OpenCost](https://github.com/opencost/opencost)/[Komiser](https://github.com/tailwarden/komiser) - monitor Kubernetes/Cloud cost
> - [Infracost](https://github.com/infracost/infracost) - estimate cloud cost from Terraform
> - [Superset](https://github.com/apache/superset) - visualise and explore your cost (and much more)
---

# Guardrails: Ensuring No Unexpected Regressions

Once teams have cost observability in place, they need guardrails to ensure optimisations don't degrade performance or reliability. Beyond real-time monitoring, it's critical to compare cost and performance trends continuously to catch unintended side effects. In future posts, we'll dive deep into implementing comprehensive observability for your services.

Some key service-level metrics to track:
- Latencies (p50, p99)
- Error rates (e.g. ratio of 5XX HTTP responses)
- Pod restarts & Out of Memory (OOM) events
- Resource utilization (CPU, memory)
- Runtime metrics - e.g. scheduling latencies and Garbage Collector (GC)

> :bulb: If your service runs on a mix of instance types, break down performance metrics by underlying infrastructure to spot anomalies.

In most organisations you will also have several user metrics to track, which should be included in decision-making process - for example user satisfaction, time spent, and retention.

There is a lot of cost optimisations that can be done without impacting user experience but at some point we will need to decide if given functionality/feature is worth its price. The key is making these trade-offs explicit and intentional. A/B testing has proven to be an invaluable tool allowing us to make data-driven decisions about trade-offs between cost and user experience.

>:bulb: When running A/B tests make sure that relevant engineering/cost metrics are included in the experimentation and that you have a way of displaying your service metrics based on the experimentation variant(s). 

---

# Getting Started

With observability and guardrails in place, you're ready to begin optimisation efforts. While some projects will require significant investment, there are several quick wins you can pursue immediately.

Start by conducting a thorough review of your current costs and eliminating unused resources. In large environments, it's common to have thousands of resources with a long tail of low-cost items. Focus initially on resources that account for 80-90% of your costs, leaving smaller optimisations for later.

As a next step take a look at your compute resources, where there are multiple simple steps to avoid spending more than necessary:
- Ensure no idle VMs are running - this is quite common with development VMS and especially ones with attached accelerators (GPU/TPU) as those are quite expensive. In some cases this will be supported out of the box (e.g. for [Vertex instances](https://cloud.google.com/vertex-ai/docs/workbench/instances/idle-shutdown)).
- Right size your deployments - review things like [requested](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) cpu/memory and min replicas
- Apply autoscalers on your workloads - load for services will usually vary throughout the day/week/year. In most cases it is as simple as setting [Horizontal Pod Autoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) (HPA) to target a specific CPU utilisation - 60-70% is usually a good starting point.

Moving on to storage - the most important question is how long your data truly needs to be retained. While some data might need indefinite storage, most becomes obsolete quickly. Most cloud providers offer different [storage classes](https://cloud.google.com/storage/docs/storage-classes) with use cases depending on access patterns and required durability. The default class is quite sensible but also the most expensive - even up to 10x compared to archive storage. [Lifecycle policies](https://cloud.google.com/storage/docs/lifecycle) are allow you to automatically manage retention and storage class transitions - a common pattern is to transfer data to coldline/archive storage after initial period when it might be accessed more often. Think about down sampling if historical data is required long-term.


You should be able to avoid network charges by staying within the same [zone/region](https://cloud.google.com/compute/docs/regions-zones) as incurring inter-zone/region charges can quickly increase your cost. Consider single-zone deployment options for some of your services and databases. While it is not recommended for crucial production workloads, this approach can significantly reduce costs for non-critical resources. A popular inefficiency is processing data in a different region than where it is stored. While sometimes this is the only viable strategy (e.g. due to regional availability of some features/products) it should generally be avoided.

---

# Common Pitfalls

Cost optimisation efforts often appear straightforward until you start measuring their actual impact. Cloud pricing models are complex, with different discounting strategies, provisioning models, and fluctuating rates. This makes it difficult to directly correlate usage reductions to final cost savings. This challenge is particularly evident with compute resources, where pricing depends on factors like provisioning model, availability of spot instances, and committed use discounts.

**Access and Visibility** issues often plague optimisation efforts. Restricted access to cost data prevents teams from making informed decisions, while lack of clear cost ownership leads to diffused responsibility. Insufficient tooling for tracking and analyzing costs makes it impossible to measure improvement.

**Process and Culture** problems can derail even well-planned initiatives. Organizations frequently lack standardized processes for resource labeling and cost attribution. Many treat cost optimisation as a one-time project rather than an ongoing practice, or fail to prioritize cost efficiency in engineering decisions. 

**Technical Oversights** are common, including overlooking managed service markups, ignoring accumulating storage and logging costs, and not accounting for network transfer fees. Teams often underutilize available discounts, over-provision resources "just in case",keep development environments running 24/7, or fail to clean up unused resources promptly.

---

# FinOps

The challenges of cloud cost management have given rise to FinOps, a framework for cloud financial management. Its core principles align closely with what we've discussed:

1. Data should be accessible and timely
2. Decisions are driven by business value
3. Everyone takes ownership
4. Teams need to collaborate
5. A centralized team drives FinOps
6. Take advantage of variable cost model

> :memo: You can learn more about FinOps [here](https://www.finops.org/framework/).

---

# Let's Go!

Cost optimisation is a vast and complex topic. Even establishing proper observability can be challenging when dealing with multiple cloud providers and services. While technical solutions are important (and often the most exciting), cost efficiency is as much about process and mindset as it is about specific optimisations. It's crucial to understand that this is a continuous journey, not a one-time effort that can be completed and forgotten.

In upcoming posts in this series, we'll explore fundamental cost optimisation techniques and share practical tips based on real-world experience.  We'll start with service-level optimisations before expanding into infrastructure, networking, and beyond. While some content will be specific to technologies I use daily (Go, GCP), most practices will be applicable to the majority of large scale workloads. 

Stay tuned! 🚀