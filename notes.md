
Topics covered:
(Observability - understanding how your application performs and setting up guardrails)
- Metrics
- Profiling
- Measuring impact

(Service-level optimisations)
- Runtime configuration
- Caching

(Compute Optimisations)
- Choosing right infrastructure
- Right-sizing
- Autoscaling
- Cheaper compute: commited use and spot

(Networking)
- What else am I paying for?
- Networking optimisations
- Weighted Load balancing
- Zone aware routing

(Other)
- Product tradeoffs + A/B experiments
- Closing thought and my experiences
    - Cost observability is hard on its own
    - Cross-team collaboration on costs - setups that worked and which did not
    - Aligning incentives


# Load testing
A topic we have not touched on in this series is load testing
Load testing is a topic I have connected to (cost) optimisations most of my career but won't cover in this series. While on paper this is a great way to ensure your application runs as efficiently as possible it is quite troublesome to do properly in huge % of real world scenarios. Some of the issues I've found:
- Downstream dependencies
- Simulating real-world traffic
- Not impacting external systems
-

Despite those limitations, I've had success with optimising using load testing - back when I was working on a vector database. In that case I was able not only to find optimal infrastructure for our deployments but also tune parameters to find optimal parameters for performance by treating it as optimisation problem (grid-scan etc)

Load testing is perfect for certain scenarios - like testing a stateless API service, a data processing pipeline, or a standalone web application. However, when dealing with services that have multiple downstream dependencies running in production, load testing becomes significantly more complex. Setting up a realistic load test would require either a sophisticated staging environment mirroring production or careful coordination to avoid impacting production dependencies. This is where profiling comes in as a powerful alternative.


If your use case is similar I'd highly recommend [k6s](https://k6.io/) which is very easy to get started with and powerfull at the same time. I've been using it for past for stress testing 

- When operating at large enough scale, self-hosting might be a better option, which also helps prevent vendor lock in
- Applying TTL etc. on logs/storage is a great starting point. This can be done through lifecycle policies 